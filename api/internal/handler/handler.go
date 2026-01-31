package handler

import "github.com/ElshadHu/vulnly/api/internal/repository"

type API struct {
	repo *repository.DynamoDB
}

func New(repo *repository.DynamoDB) *API {
	return &API{repo: repo}
}
