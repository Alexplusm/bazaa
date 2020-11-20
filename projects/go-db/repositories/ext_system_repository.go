package repositories

import (
	"context"
	"fmt"

	"github.com/Alexplusm/bazaa/projects/go-db/interfaces"
	"github.com/Alexplusm/bazaa/projects/go-db/objects/dao"
)

type ExtSystemRepository struct {
	DBConn interfaces.IDBHandler
}

const (
	insertExtSystemWithIDStatement = `
INSERT INTO ext_systems ("ext_system_id", "description", "post_results_url")
VALUES ($1, $2, $3)
RETURNING "ext_system_id";
`
	insertExtSystemWithoutIDStatement = `
INSERT INTO ext_systems ("description", "post_results_url")
VALUES ($1, $2)
RETURNING "ext_system_id";
`
)

func (repo *ExtSystemRepository) InsertExtSystem(
	extSystemDAO dao.ExtSystemDAO,
) (string, error) {
	p := repo.DBConn.GetPool()
	ctx := context.Background()
	conn, err := p.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("insert extSystem: acquire connection: %v", err)
	}
	defer conn.Release()

	var args []interface{}
	var statement string

	if extSystemDAO.HasID() {
		statement = insertExtSystemWithIDStatement
		args = []interface{}{extSystemDAO.ID, extSystemDAO.Description, extSystemDAO.PostResultsURL}
	} else {
		statement = insertExtSystemWithoutIDStatement
		args = []interface{}{extSystemDAO.Description, extSystemDAO.PostResultsURL}
	}

	row := conn.QueryRow(ctx, statement, args...)

	var extSystemID string
	if err = row.Scan(&extSystemID); err != nil {
		return "", fmt.Errorf("insert extSystem: %v", err)
	}

	return extSystemID, nil
}

func (repo *ExtSystemRepository) SelectExtSystems() ([]dao.ExtSystemDAO, error) {
	// TODO: for web client need list of extSystems [{id, description}, ...]
	return nil, nil
}
