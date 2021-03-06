package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"github.com/Alexplusm/bazaa/projects/go-db/interfaces"
	"github.com/Alexplusm/bazaa/projects/go-db/objects/dao"
)

type AnswerRepo struct {
	DBConn interfaces.IDBHandler
}

const (
	insertAnswerStatement = `
INSERT INTO answers ("screenshot_id", "game_id", "user_id", "value", "answer_date")
VALUES ($1, $2, $3, $4, $5);
`
	selectAnsweredScreenshotsCountStatement = `
select COUNT(DISTINCT screenshot_id) FROM answers
WHERE answers.game_id = ($1)
`
	selectUniqueUsersInGameStatement = `
SELECT DISTINCT user_id FROM answers
WHERE answers.game_id = ($1)
`
	selectAnswersByGame = `
SELECT 
ans.game_id, ans.user_id, ans.screenshot_id, ans.answer_date, ans.value, s.users_answer, s.expert_answer
FROM answers ans
INNER JOIN screenshots s
ON s.screenshot_id = ans.screenshot_id
WHERE
ans.game_id = ($1) AND
(ans.answer_date BETWEEN ($2) AND ($3))
`
	selectAnswersByUserAndGameStatement = `
SELECT 
ans.game_id, ans.user_id, ans.screenshot_id, ans.answer_date, ans.value, s.users_answer, s.expert_answer
FROM answers ans
INNER JOIN screenshots s
ON s.screenshot_id = ans.screenshot_id
WHERE
ans.user_id = ($1) AND
ans.game_id = ($2) AND
(ans.answer_date BETWEEN ($3) AND ($4))
`
	selectScreenshotResultsStatement = `
SELECT 
ans.game_id, ans.user_id, ans.screenshot_id, ans.answer_date, ans.value, s.users_answer, s.expert_answer
FROM answers ans
INNER JOIN screenshots s
ON s.screenshot_id = ans.screenshot_id
WHERE
ans.game_id = ($1) AND ans.screenshot_id = ($2) 
`
)

func (repo *AnswerRepo) InsertList(answers []dao.AnswerInsertDAO) {
	for _, answer := range answers {
		err := repo.InsertOne(answer)
		if err != nil {
			log.Error("answer repo: insert list: ", err)
		}
	}
}

func (repo *AnswerRepo) InsertOne(answer dao.AnswerInsertDAO) error {
	p := repo.DBConn.GetPool()
	conn, err := p.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("insert game: acquire connection: %v", err)
	}
	defer conn.Release()

	nowTimestamp := time.Now().Unix()
	row, err := conn.Query(
		context.Background(),
		insertAnswerStatement,
		answer.ScreenshotID, answer.GameID, answer.UserID, answer.Value, nowTimestamp,
	)
	if err != nil {
		return fmt.Errorf("insert answer: %v", err)
	}
	row.Close()

	return nil
}

func (repo *AnswerRepo) SelectScreenshotResult(gameID, screenshotID string) ([]dao.AnswerScreenshotRetrieveDAO, error) {
	p := repo.DBConn.GetPool()
	conn, err := p.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("select screenshot result: acquire connection: %v", err)
	}
	defer conn.Release()

	rows, err := conn.Query(
		context.Background(), selectScreenshotResultsStatement,
		gameID, screenshotID,
	)
	if err != nil {
		return nil, fmt.Errorf("select screenshot result: %v", err)
	}
	defer rows.Close()

	return retrieveAnswerScreenshotList(rows)
}

// TODO: check
func (repo *AnswerRepo) SelectAnsweredScreenshotsByGame(
	gameID string,
) (dao.AnsweredScreenshotsDAO, error) {
	p := repo.DBConn.GetPool()
	conn, err := p.Acquire(context.Background())
	if err != nil {
		return dao.AnsweredScreenshotsDAO{},
			fmt.Errorf("select answered screenshots by game: acquire connection: %v", err)
	}
	defer conn.Release()

	var count int64

	row := conn.QueryRow(context.Background(), selectAnsweredScreenshotsCountStatement, gameID)
	if row.Scan(&count) != nil {
		return dao.AnsweredScreenshotsDAO{},
			fmt.Errorf("select answered screenshots by game: %v", err)
	}

	rows, err := conn.Query(context.Background(), selectUniqueUsersInGameStatement, gameID)
	if err != nil {
		return dao.AnsweredScreenshotsDAO{},
			fmt.Errorf("select answered screenshots by game: %v", err)
	}
	defer rows.Close()

	listUsers := make([]string, 0, 1024)

	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		if err != nil {
			log.Error("select answered screenshots by game: ", err)
			continue
		}
		listUsers = append(listUsers, user)
	}
	if rows.Err() != nil {
		log.Error("select answered screenshots by game: ", rows.Err())
	}

	res := dao.AnsweredScreenshotsDAO{
		Count:  int(count),
		UserID: listUsers,
	}

	fmt.Printf("reees: %+v\n", res)

	return res, nil
}

func (repo *AnswerRepo) SelectListTODO(gameID string, from, to time.Time) ([]dao.AnswerScreenshotRetrieveDAO, error) {
	p := repo.DBConn.GetPool()
	conn, err := p.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("select answers by user: acquire connection: %v", err)
	}
	defer conn.Release()

	rows, err := conn.Query(
		context.Background(), selectAnswersByGame,
		gameID, from.Unix(), to.Unix(),
	)

	if err != nil {
		return nil, fmt.Errorf("todooo: %v", err)
	}
	defer rows.Close()

	return retrieveAnswerScreenshotList(rows)
}

func (repo *AnswerRepo) SelectListByUserAndGame(
	userID string, gameID string, from, to time.Time,
) ([]dao.AnswerScreenshotRetrieveDAO, error) {
	p := repo.DBConn.GetPool()
	conn, err := p.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("select answers by user and game: acquire connection: %v", err)
	}
	defer conn.Release()

	rows, err := conn.Query(
		context.Background(), selectAnswersByUserAndGameStatement,
		userID, gameID, from.Unix(), to.Unix(),
	)

	if err != nil {
		return nil, fmt.Errorf("select answers by user and game: %v", err)
	}
	defer rows.Close()

	return retrieveAnswerScreenshotList(rows)
}

func retrieveAnswerScreenshotList(rows pgx.Rows) ([]dao.AnswerScreenshotRetrieveDAO, error) {
	list := make([]dao.AnswerScreenshotRetrieveDAO, 0, 1024)

	for rows.Next() {
		v := dao.AnswerScreenshotRetrieveDAO{}

		err := rows.Scan(
			&v.GameID, &v.UserID, &v.ScreenshotID,
			&v.AnswerDate, &v.Value, &v.UsersAnswer, &v.ExpertAnswer,
		)
		if err != nil {
			return nil, fmt.Errorf("retrieve list: %v", err)
		}
		list = append(list, v)
	}

	return list, nil
}
