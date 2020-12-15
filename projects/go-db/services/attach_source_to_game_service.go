package services

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Alexplusm/bazaa/projects/go-db/consts"
	"github.com/Alexplusm/bazaa/projects/go-db/interfaces"
	"github.com/Alexplusm/bazaa/projects/go-db/objects/bo"
	"github.com/Alexplusm/bazaa/projects/go-db/objects/dao"
	"github.com/Alexplusm/bazaa/projects/go-db/utils/fileutils"
	"github.com/Alexplusm/bazaa/projects/go-db/utils/logutils"
)

type AttachSourceToGameService struct {
	GameRepo       interfaces.IGameRepo
	SourceRepo     interfaces.ISourceRepo
	ScreenshotRepo interfaces.IScreenshotRepo
	FileService    interfaces.IFileService
}

func (service *AttachSourceToGameService) AttachArchives(
	gameID string, archives []*multipart.FileHeader,
) error {
	filenames, err := service.FileService.CopyFiles(archives, consts.MediaTempDir)
	if err != nil {
		return fmt.Errorf("attach zip archive: %v", err)
	}

	images, err := service.FileService.UnzipImages(filenames)
	if err != nil {
		return fmt.Errorf("attach zip archive: %v", err)
	}

	// TODO: Source Service
	// TODO: another func
	archivesFilename := make([]string, 0, len(archives))
	for _, archive := range archives {
		archivesFilename = append(archivesFilename, archive.Filename)
	}
	value := strings.Join(archivesFilename, ",")

	source := dao.SourceInsertDAO{
		SourceBaseDAO: dao.SourceBaseDAO{
			Type: consts.ArchiveSourceType, CreatedAt: time.Now().Unix(), GameID: gameID, Value: value,
		},
	}

	sourceID, err := service.SourceRepo.InsertOne(source)
	if err != nil {
		return fmt.Errorf("attach zip archive: %v", err)
	}

	a, b := split(images, gameID, sourceID)
	err = service.ScreenshotRepo.InsertList(a)
	if err != nil {
		return fmt.Errorf("attach zip archive: %v", err)
	}

	err = service.ScreenshotRepo.InsertListWithExpertAnswer(b)
	if err != nil {
		return fmt.Errorf("attach zip archive: %v", err)
	}

	removeArchives(filenames)

	return nil
}

func (service *AttachSourceToGameService) AttachSchedules(gameID string) error {
	fmt.Println("Schedules attaching coming soon ... : gameID =", gameID)
	return nil
}

func (service *AttachSourceToGameService) AttachGameResults(gameID string, params bo.AttachGameParams) error {
	// TODO: sourceService.method
	sources, err := service.SourceRepo.SelectListByGame(gameID)
	if err != nil {
		return fmt.Errorf("%v AttachGameResults: %v", logutils.GetStructName(service), err)
	}

	exist := false
	for _, source := range sources {
		if source.Value == params.SourceGameID {
			exist = true
		}
	}

	if exist {
		// TODO: kek
		return fmt.Errorf("source exist")
	}
	// TODO: sourceService.method

	source := dao.SourceInsertDAO{
		SourceBaseDAO: dao.SourceBaseDAO{
			Type: consts.ArchiveSourceType, CreatedAt: time.Now().Unix(), GameID: gameID, Value: params.SourceGameID,
		},
	}

	sourceID, err := service.SourceRepo.InsertOne(source)
	if err != nil {
		return fmt.Errorf("TOODOOO: %v", err)
	}

	screenshots, err := service.ScreenshotRepo.SelectListByGameID(params.SourceGameID)
	if err != nil {
		return fmt.Errorf("%v AttachGameResults: %v", logutils.GetStructName(service), err)
	}

	newScreenshots := make([]dao.ScreenshotDAO, 0, len(screenshots))

	for _, screenshot := range screenshots {
		if string(screenshot.UsersAnswer) == params.Answer {
			ddao := dao.ScreenshotDAO{
				Filename: screenshot.Filename,
				GameID:   screenshot.GameID,
				SourceID: sourceID,
			}
			newScreenshots = append(newScreenshots, ddao)
		}
	}

	err = service.ScreenshotRepo.InsertList(newScreenshots)
	if err != nil {
		return fmt.Errorf("%v AttachGameResults: %v", logutils.GetStructName(service), err)
	}

	return nil
}

func removeArchives(filenames []string) {
	for _, fn := range filenames {
		err := fileutils.RemoveFile(consts.MediaTempDir, fn)
		if err != nil {
			log.Error("remove archive: ", err)
		}
	}
}

func split(
	images []bo.ImageParsingResult, gameID, sourceID string,
) ([]dao.ScreenshotDAO, []dao.ScreenshotWithExpertAnswerDAO) {
	mmap := make(map[string]bool)
	imagesWithoutExpertAnswer := make([]dao.ScreenshotDAO, 0, len(images))
	imagesWithExpertAnswer := make([]dao.ScreenshotWithExpertAnswerDAO, 0, len(images))

	// INFO: Когда загружаем несколько архивов могут быть попасться одинаковые файлы
	// -> обрабатываем эту ситуацию
	for _, image := range images {
		if !mmap[image.Filename] {
			mmap[image.Filename] = true
			if image.Category == UndefinedCategory {
				screen := dao.ScreenshotDAO{image.Filename, gameID, sourceID}

				imagesWithoutExpertAnswer = append(imagesWithoutExpertAnswer, screen)
			} else {
				screen := dao.ScreenshotWithExpertAnswerDAO{
					ScreenshotDAO: dao.ScreenshotDAO{
						Filename: image.Filename, GameID: gameID, SourceID: sourceID,
					},
					ExpertAnswer: image.Category,
				}
				imagesWithExpertAnswer = append(imagesWithExpertAnswer, screen)
			}
		}
	}

	return imagesWithoutExpertAnswer, imagesWithExpertAnswer
}
