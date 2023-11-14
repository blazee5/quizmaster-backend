package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/question"
	"github.com/blazee5/quizmaster-backend/lib/http_utils"
	"github.com/blazee5/quizmaster-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
)

type Handler struct {
	log     *zap.SugaredLogger
	service question.Service
}

func NewHandler(log *zap.SugaredLogger, service question.Service) *Handler {
	return &Handler{log: log, service: service}
}

func (h *Handler) CreateQuestion(c echo.Context) error {
	var input domain.Question

	userId := c.Get("userId").(int)
	quizId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	input.QuizId = quizId

	if err := c.Bind(&input); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	answers := c.FormValue("answers")

	if answers == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Answers is required field.",
		})
	}

	if err := json.Unmarshal([]byte(answers), &input.Answers); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	if err := c.Validate(&input); err != nil {
		validateErr := err.(validator.ValidationErrors)

		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": response.ValidationError(validateErr),
		})
	}

	if _, err := os.Stat("public"); os.IsNotExist(err) {
		_ = os.Mkdir("public", os.ModePerm)
	}

	file, err := c.FormFile("image")

	if err == nil {
		if err := http_utils.UploadFile(file, "public/"+file.Filename); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "server error",
			})
		}
		input.Image = file.Filename
	}

	id, err := h.service.Create(c.Request().Context(), userId, quizId, input)

	if err != nil {
		h.log.Infof("error while create quiz: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": id,
	})
}

func (h *Handler) GetQuizQuestions(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	questions, err := h.service.GetQuestionsById(c.Request().Context(), id)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if err != nil {
		h.log.Infof("error while get questions by quiz id: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, questions)
}
