package services

import (
	"fmt"

	"github.com/Alexplusm/bazaa/projects/go-db/dao"
	"github.com/Alexplusm/bazaa/projects/go-db/domain"
	"github.com/Alexplusm/bazaa/projects/go-db/interfaces"
)

type ICreateGameService interface {
	CreateGame()
}

type CreateGameService struct {
	Repository interfaces.IGameRepository
}

func (service *CreateGameService) CreateGame(game domain.GameBO) (string, error) {
	fmt.Printf("CreateGame service: %+v\n", game)

	dao := dao.GameDAO{}
	dao.FromBO(game)

	return service.Repository.CreateGame(dao)
}
