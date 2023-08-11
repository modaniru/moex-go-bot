package services

import (
	"context"

	"github.com/modaniru/moex-telegram-bot/internal/storage/gen"
)

type UserStorage interface {
	CreateUser(ctx context.Context, id int32) error
	GetUser(ctx context.Context, id int32) (gen.User, error)
	DeleteUser(ctx context.Context, id int32) error
	Follow(ctx context.Context, id int32) error
	Unfollow(ctx context.Context, id int32) error
}

type UserService struct {
	storage UserStorage
}

func NewUserService(userStorage UserStorage) *UserService {
	return &UserService{storage: userStorage}
}

func (u *UserService) Register(id int) error {
	return u.storage.CreateUser(context.Background(), int32(id))
}

func (u *UserService) GetUserById(id int) (gen.User, error) {
	return u.storage.GetUser(context.Background(), int32(id))
}

func (u *UserService) DeleteUserById(id int) error {
	return u.storage.DeleteUser(context.Background(), int32(id))
}

func (u *UserService) FollowUser(id int) error{
	return u.storage.Follow(context.Background(), int32(id))
}

func (u *UserService) UnfollowUser(id int) error{
	return u.storage.Unfollow(context.Background(), int32(id))
}