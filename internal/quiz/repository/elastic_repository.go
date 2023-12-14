package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/elastic/go-elasticsearch/v8"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"strconv"
)

type ElasticRepository struct {
	client *elasticsearch.Client
	tracer trace.Tracer
}

func NewElasticRepository(client *elasticsearch.Client, tracer trace.Tracer) *ElasticRepository {
	return &ElasticRepository{client: client, tracer: tracer}
}

func (repo *ElasticRepository) CreateIndex(ctx context.Context, input models.QuizInfo) error {
	ctx, span := repo.tracer.Start(ctx, "quizElasticRepo.CreateIndex")
	defer span.End()

	data, err := json.Marshal(input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	resp, err := repo.client.Index(
		"quizzes",
		bytes.NewReader(data),
		repo.client.Index.WithContext(ctx),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		span.RecordError(err)
		span.SetStatus(codes.Error, resp.String())

		return err
	}

	return nil
}

func (repo *ElasticRepository) SearchIndex(ctx context.Context, input string) ([]models.QuizInfo, error) {
	ctx, span := repo.tracer.Start(ctx, "quizElasticRepo.SearchIndex")
	defer span.End()

	query := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"should": []map[string]any{
					{
						"prefix": map[string]any{
							"title": input,
						},
					},
					{
						"fuzzy": map[string]any{
							"title": map[string]any{
								"value":     input,
								"fuzziness": 2,
							},
						},
					},
				},
			},
		},
	}

	dataBytes, err := json.Marshal(&query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	response, err := repo.client.Search(
		repo.client.Search.WithContext(ctx),
		repo.client.Search.WithIndex("quizzes"),
		repo.client.Search.WithBody(bytes.NewReader(dataBytes)),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}
	defer response.Body.Close()

	if response.IsError() {
		span.RecordError(err)
		span.SetStatus(codes.Error, response.String())

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
	if err = json.NewDecoder(response.Body).Decode(&hits); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	quizzes := make([]models.QuizInfo, len(hits.Hits.Hits))
	for i, source := range hits.Hits.Hits {
		quizzes[i] = source.Source
	}
	return quizzes, nil
}

func (repo *ElasticRepository) UpdateIndex(ctx context.Context, id int, input models.QuizInfo) error {
	ctx, span := repo.tracer.Start(ctx, "quizElasticRepo.DeleteIndex")
	defer span.End()

	body, err := json.Marshal(map[string]any{
		"script": map[string]any{
			"source": fmt.Sprintf("ctx._source.title = '%s'; ctx._source.description = '%s'", input.Title, input.Description),
			"lang":   "painless",
		},
		"query": map[string]any{
			"match": map[string]any{
				"id": id,
			},
		},
	})

	resp, err := repo.client.UpdateByQuery(
		[]string{"quizzes"},
		repo.client.UpdateByQuery.WithBody(bytes.NewReader(body)),
		repo.client.UpdateByQuery.WithContext(ctx),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		span.RecordError(err)
		span.SetStatus(codes.Error, resp.String())

		return err
	}

	return nil
}

func (repo *ElasticRepository) DeleteIndex(ctx context.Context, id int) error {
	ctx, span := repo.tracer.Start(ctx, "quizElasticRepo.DeleteIndex")
	defer span.End()

	documentID := strconv.Itoa(id)

	body, err := json.Marshal(map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				"id": documentID,
			},
		},
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	resp, err := repo.client.DeleteByQuery(
		[]string{"quizzes"},
		bytes.NewReader(body),
		repo.client.DeleteByQuery.WithContext(ctx),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		span.RecordError(err)
		span.SetStatus(codes.Error, resp.String())

		return err
	}

	return nil
}
