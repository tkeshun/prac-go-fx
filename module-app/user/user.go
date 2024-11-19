package user

import "context"

type Service interface {
	SaveUser(ctx context.Context, username string) error
	User(ctx context.Context, id string) (string, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) SaveUser(ctx context.Context, username string) error {
	return s.repository.Save(ctx, username)
}

func (s *service) User(ctx context.Context, id string) (string, error) {
	return s.repository.Find(ctx, id)
}
