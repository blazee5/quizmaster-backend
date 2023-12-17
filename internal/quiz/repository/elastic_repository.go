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
	"math"
	"strconv"
)

type ElasticRepository struct {
	client *elasticsearch.Client
	tracer trace.Tracer
}

func NewElasticRepository(client *elasticsearch.Client, tracer trace.Tracer) *ElasticRepository {
	return &ElasticRepository{client: client, tracer: tracer}
}

func (repo *ElasticRepository) CreateIndex(ctx context.Context, input models.Quiz) error {
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

func (repo *ElasticRepository) SearchIndex(ctx context.Context, input, sortBy, sortDir string, page, size int) (models.QuizList, error) {
	ctx, span := repo.tracer.Start(ctx, "quizElasticRepo.SearchIndex")
	defer span.End()

	query := map[string]any{
		"from": (page - 1) * size,
		"size": size,
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
		"sort": []map[string]any{
			{
				sortBy: sortDir,
			},
		},
	}

	dataBytes, err := json.Marshal(&query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.QuizList{}, err
	}

	response, err := repo.client.Search(
		repo.client.Search.WithContext(ctx),
		repo.client.Search.WithIndex("quizzes"),
		repo.client.Search.WithBody(bytes.NewReader(dataBytes)),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.QuizList{}, err
	}
	defer response.Body.Close()

	if response.IsError() {
		span.RecordError(err)
		span.SetStatus(codes.Error, response.String())

		return models.QuizList{}, err
	}
	type EsHits struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source models.Quiz `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	var hits EsHits
	if err = json.NewDecoder(response.Body).Decode(&hits); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.QuizList{}, err
	}

	quizzes := make([]models.Quiz, len(hits.Hits.Hits))

	for i, source := range hits.Hits.Hits {
		quizzes[i] = source.Source
	}

	total := hits.Hits.Total.Value

	return models.QuizList{
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(size))),
		Page:       page,
		Size:       size,
		Quizzes:    quizzes,
	}, nil
}

func (repo *ElasticRepository) UpdateIndex(ctx context.Context, id int, input models.Quiz) error {
	ctx, span := repo.tracer.Start(ctx, "quizElasticRepo.DeleteIndex")
	defer span.End()

	body, err := json.Marshal(map[string]any{
		"script": map[string]any{
			"source": fmt.Sprintf("ctx._source.title = '%s'; ctx._source.description = '%s'; ctx._source.image = '%s'", input.Title, input.Description, input.Image),
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
