package repository

import (
	"context"
	"encoding/json"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticRepository struct {
	esclient *elasticsearch.Client
}

func NewElasticRepository(esclient *elasticsearch.Client) *ElasticRepository {
	return &ElasticRepository{esclient: esclient}
}

func (repo *ElasticRepository) CreateIndex(ctx context.Context, input domain.Quiz) error {
	_, err := json.Marshal(input)

	if err != nil {
		return err
	}

	//res, err := repo.esclient.Create("quizzes", input., quizBytes)
	//fmt.Println(res)
	//
	if err != nil {
		return err
	}

	return nil
}

func (repo *ElasticRepository) UpdateIndex(ctx context.Context, input domain.Quiz) error {
	//TODO implement me
	panic("implement me")
}

func (repo *ElasticRepository) DeleteIndex(ctx context.Context, input domain.Quiz) error {
	//TODO implement me
	panic("implement me")
}

func (repo *ElasticRepository) SearchIndex(ctx context.Context, input string) ([]domain.Quiz, error) {
	//TODO implement me
	panic("implement me")
}
