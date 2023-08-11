package service

import (
	"github.com/modaniru/moex-telegram-bot/internal/service/services"
	"github.com/modaniru/moex-telegram-bot/internal/storage"
	"github.com/modaniru/moex-telegram-bot/internal/storage/gen"
)

type UserService interface {
	Register(id int) error
	GetUserById(id int) (gen.User, error)
	DeleteUserById(id int) error
	FollowUser(id int) error
	UnfollowUser(id int) error
}

type Service struct {
	UserService
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		UserService: services.NewUserService(storage),
	}
}
