package services

import (
	"fmt"

	"github.com/Alexplusm/bazaa/projects/go-db/interfaces"
	"github.com/Alexplusm/bazaa/projects/go-db/objects/bo"
	"github.com/Alexplusm/bazaa/projects/go-db/objects/dao"
)

// todo: remove! unused
type ICreateGameService interface {
	CreateGame()
}

type CreateGameService struct {
	Repository interfaces.IGameRepository
}

func (service *CreateGameService) CreateGame(game bo.GameBO) (string, error) {
	fmt.Printf("CreateGame service: %+v\n", game)

	gameDAO := dao.GameDAO{}
	gameDAO.FromBO(game)

	return service.Repository.InsertGame(gameDAO)
}