package services

import (
	"fmt"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Alexplusm/bazaa/projects/go-db/consts"
	"github.com/Alexplusm/bazaa/projects/go-db/interfaces"
	"github.com/Alexplusm/bazaa/projects/go-db/objects/bo"
	"github.com/Alexplusm/bazaa/projects/go-db/objects/dao"
	"github.com/Alexplusm/bazaa/projects/go-db/objects/dto"
	"github.com/Alexplusm/bazaa/projects/go-db/utils/timeutils"
)

type AnswerService struct {
	AnswerRepo interfaces.IAnswerRepo
}

func (service *AnswerService) GetScreenshotResults(
	gameID, screenshotID string,
) ([]dto.UserAnswerForScreenshotResultDTO, error) {
	res, err := service.AnswerRepo.SelectScreenshotResult(gameID, screenshotID)
	if err != nil {
		return nil, fmt.Errorf("get screenshot results: %v", err)
	}

	res_len := len(res)
	list := make([]dto.UserAnswerForScreenshotResultDTO, 0, res_len)

	for _, r := range res {
		dtoo := dto.UserAnswerForScreenshotResultDTO{}
		dtoo.UserID = r.UserID
		dtoo.Answer = r.Value
		if res_len < consts.RequiredAnswerCountToFinishScreenshot {
			dtoo.Result = "inProcess"
		} else {
			if string(r.UsersAnswer) == "-1" { // TODO: refactor
				dtoo.Result = "undefined"
			} else {
				if r.Value == string(r.UsersAnswer) { // TODO: refactor
					dtoo.Result = "right"
				} else {
					dtoo.Result = "wrong"
				}
			}
		}
		list = append(list, dtoo)
	}

	fmt.Printf("res: %+v\n", res)

	return list, nil
}

// TODO: refactor | screenshotResult (добавить поле usersResult)
func (service *AnswerService) GetUserStatistics(
	userID string, gameIDs []string, from, to time.Time,
) ([]bo.StatisticAnswersDateSlicedBO, error) {
	userAnswers := make([]dao.AnswerScreenshotRetrieveDAO, 0, 1024)

	for _, gameID := range gameIDs {
		oneRes, err := service.AnswerRepo.SelectListByUserAndGame(userID, gameID, from, to)
		if err != nil {
			log.Error("user statistics service: ", err)
			continue
		}
		userAnswers = append(userAnswers, oneRes...)
		fmt.Printf("oneRes: %+v\n", oneRes)
	}

	sort.SliceStable(userAnswers, func(i, j int) bool {
		return userAnswers[i].AnswerDate < userAnswers[j].AnswerDate
	})

	if len(userAnswers) == 0 {
		return make([]bo.StatisticAnswersDateSlicedBO, 0, 0), nil
	}

	start := timeutils.TrimTime(time.Unix(userAnswers[0].AnswerDate, 0))
	end := time.Unix(userAnswers[len(userAnswers)-1].AnswerDate, 0)
	end = end.AddDate(0, 0, 1)
	end = timeutils.TrimTime(end)

	results := countRes(userAnswers, start, end)

	return results, nil
}

func (service *AnswerService) GetUsersAndScreenshotCountByGame(
	gameID string,
) (dao.AnsweredScreenshotsDAO, error) {
	return service.AnswerRepo.SelectAnsweredScreenshotsByGame(gameID)
}

func (service *AnswerService) ABC(gameID string, from, to time.Time) ([]dao.AnswerScreenshotRetrieveDAO, error) {
	return service.AnswerRepo.SelectListTODO(gameID, from, to)
}

func countRes(userAnswers []dao.AnswerScreenshotRetrieveDAO, start, end time.Time) []bo.StatisticAnswersDateSlicedBO {
	results := make([]bo.StatisticAnswersDateSlicedBO, 0, len(userAnswers))

	fmt.Printf("val: %+v\n", userAnswers)

	currentDay := start
	for currentDay.Before(end) {
		for _, userAnswer := range userAnswers {
			date := time.Unix(userAnswer.AnswerDate, 0)
			nextDay := currentDay.AddDate(0, 0, 1)

			if currentDay.Before(date) && date.Before(nextDay) {
				curIdx := -1
				for i := range results {
					if results[i].Date.Equal(currentDay) {
						curIdx = i
						break
					}
				}
				if curIdx == -1 {
					stat := bo.StatisticAnswersBO{MatchWithExpert: -1}
					stat.Increase(userAnswer.Value, string(userAnswer.ExpertAnswer), string(userAnswer.UsersAnswer))
					results = append(results, bo.StatisticAnswersDateSlicedBO{
						Date:       currentDay,
						Statistics: stat,
					})
				} else {
					results[curIdx].Statistics.Increase(
						userAnswer.Value, string(userAnswer.ExpertAnswer), string(userAnswer.UsersAnswer),
					)
				}
			}
		}
		currentDay = currentDay.AddDate(0, 0, 1)
	}

	fmt.Printf("RESULTS: %+v\n", results)

	return results

}
