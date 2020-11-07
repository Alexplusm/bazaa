package interfaces

import (
	"github.com/Alexplusm/bazaa/projects/go-db/objects/dao"
)

type IGameRepository interface {
	InsertGame(game dao.GameDAO) (string, error)
	HasNotStartedGameWithSameID(gameID string) (bool, error)
}

type ISourceRepository interface {
	InsertSource(source dao.SourceDAO) (string, error)
}

type IScreenshotRepository interface {
	InsertScreenshots(screenshots []dao.ScreenshotDAO) error
	InsertScreenshotsWithExpertAnswer(screenshots []dao.ScreenshotWithExpertAnswerDAO) error
}
