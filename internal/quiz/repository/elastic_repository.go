package repository

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticRepository struct {
	client *elasticsearch.Client
}

func NewElasticRepository(client *elasticsearch.Client) *ElasticRepository {
	return &ElasticRepository{client: client}
}

func (repo *ElasticRepository) CreateIndex(ctx context.Context, input models.Quiz) error {
	data, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal quiz: %w", err)
	}

	resp, err := repo.client.Index(
		"quizzes",
		bytes.NewReader(data),
		repo.client.Index.WithContext(ctx),
	)

	if err != nil {
		return fmt.Errorf("failed to index quiz: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("error indexing quiz: %v", resp.String())
	}

	return nil
}

var (
	searchFields = []string{"title", "description"}
)

func (repo *ElasticRepository) SearchIndex(ctx context.Context, input string) ([]models.QuizInfo, error) {
	query := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":     input,
				"fuzziness": 2,
				"fields":    searchFields,
			},
		},
	}

	dataBytes, err := json.Marshal(&query)
	if err != nil {
		return nil, err
	}

	response, err := repo.client.Search(
		repo.client.Search.WithContext(ctx),
		repo.client.Search.WithIndex("quizzes"),
		repo.client.Search.WithBody(bufio.NewReader(bytes.NewReader(dataBytes))),
	)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.IsError() {
		return nil, err
	}

	type EsHits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source models.QuizInfo `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	var hits EsHits
	if err := json.NewDecoder(response.Body).Decode(&hits); err != nil {
		return nil, err
	}

	quizzes := make([]models.QuizInfo, len(hits.Hits.Hits))
	for i, source := range hits.Hits.Hits {
		quizzes[i] = source.Source
	}
	return quizzes, nil
}
