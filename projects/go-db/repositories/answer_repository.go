package repositories

import (
	"context"
	"fmt"

	"github.com/Alexplusm/bazaa/projects/go-db/interfaces"
	"github.com/Alexplusm/bazaa/projects/go-db/objects/dao"
)

type AnswerRepository struct {
	DBConn interfaces.IDBHandler
}

const (
	insertAnswerStatement = `
INSERT INTO answers ("screenshot_id", "game_id", "user_id", "value")
VALUES ($1, $2, $3, $4);
`
)

func (repo *AnswerRepository) InsertAnswers(answers []dao.AnswerDAO) {
	for _, answer := range answers {
		err := repo.InsertAnswer(answer)
		if err != nil {
			fmt.Println("err: insert answers: ", err) // TODO: log error | return error
		}
	}
}

func (repo *AnswerRepository) InsertAnswer(answer dao.AnswerDAO) error {
	p := repo.DBConn.GetPool()
	conn, err := p.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("insert game: acquire connection: %v", err)
	}
	defer conn.Release()

	row, err := conn.Query(
		context.Background(),
		insertAnswerStatement,
		answer.ScreenshotID, answer.GameID, answer.UserID, answer.Value,
	)
	if err != nil {
		return fmt.Errorf("insert answer: %v", err)
	}
	row.Close()

	return nil
}
