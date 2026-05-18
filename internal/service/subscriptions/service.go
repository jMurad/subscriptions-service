package subscriptions

import (
	"time"
)

const (
	shortTimeout  = 3 * time.Second
	mediumTimeout = 5 * time.Second
	longTimeout   = 10 * time.Second
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}
