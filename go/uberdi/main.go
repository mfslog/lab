package main

import (
	"fmt"
	"go.uber.org/dig"
)

type User struct {
	ID int64
}

type Repo interface {
	Get(int64) (*User, error)
}

type repo struct{}

func (r *repo) Get(id int64) (*User, error) {
	return &User{ID: id}, nil
}

func NewRepo() Repo {
	return &repo{}
}

type Service struct {
	r Repo
}

func NewService(r Repo) *Service {
	return &Service{r: r}
}

func main() {
	container := dig.New()

	container.Provide(NewRepo)
	container.Provide(NewService)

	container.Invoke(func(s *Service) {
		u, err := s.r.Get(1)
		fmt.Printf("user %+v err %+v", u, err)
	})

}
