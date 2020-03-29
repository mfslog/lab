//+build wireinject

package main

import (
	"fmt"
	"log"
)

import "github.com/google/wire"

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

var (
	Set = wire.NewSet(
		NewRepo,
	)
)

type Service struct {
	R Repo
}

func NewService(r Repo) (*Service, error) {
	return &Service{R: r}, nil
}

func InitApp() (*Service, error) {
	panic(wire.Build(Set, NewService))
}

func main() {
	s, err := InitApp()
	if err != nil {
		log.Fatal(err)
		return
	}
	u, err := s.R.Get(1)
	fmt.Printf("%+v,+%v", u, err)
}
