package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blazee5/testhub-backend/internal/domain"
	"github.com/blazee5/testhub-backend/internal/quiz"
	"github.com/blazee5/testhub-backend/lib/http_errors"
	"github.com/blazee5/testhub-backend/lib/http_utils"
	"github.com/blazee5/testhub-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
)

type Handler struct {
	log     *zap.SugaredLogger
	service quiz.Service
}

func NewHandler(log *zap.SugaredLogger, service quiz.Service) *Handler {
	return &Handler{log: log, service: service}
}

func (h *Handler) GetAllQuizzes(c echo.Context) error {
	quizzes, err := h.service.GetAll(c.Request().Context())

	if err != nil {
		h.log.Infof("error while get all quizzes: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, quizzes)
}

func (h *Handler) GetQuiz(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	quiz, err := h.service.GetById(c.Request().Context(), id)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if err != nil {
		h.log.Infof("error while get quiz: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, quiz)
}

func (h *Handler) CreateQuiz(c echo.Context) error {
	var input domain.Quiz

	userId := c.Get("userId").(int)

	//FIXME: переделать парсинг вопросов
	questions := c.FormValue("questions")

	err := json.Unmarshal([]byte(questions), &input.Questions)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Questions is a required field",
		})
	}

	if err := c.Bind(&input); err != nil {
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
		err = os.Mkdir("public", os.ModePerm)
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

	for idx := range input.Questions {
		file, err := c.FormFile(fmt.Sprintf("question_img%d", idx+1))

		if err == nil {
			if err := http_utils.UploadFile(file, "public/"+file.Filename); err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": "server error",
				})
			}

			input.Questions[idx].Image = file.Filename
		}
	}

	input.UserId = userId

	id, err := h.service.Create(c.Request().Context(), input)

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

func (h *Handler) SaveResult(c echo.Context) error {
	var input domain.Result

	userId := c.Get("userId").(int)
	quizId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

	if err := c.Bind(&input); err != nil {
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

	result, err := h.service.SaveResult(c.Request().Context(), userId, quizId, input)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if err != nil {
		h.log.Infof("error while save results: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"correct_answers": result,
	})
}

func (h *Handler) DeleteQuiz(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	userId := c.Get("userId").(int)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	err = h.service.Delete(c.Request().Context(), userId, id)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if errors.Is(err, http_errors.PermissionDenied) {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "permission denied",
		})
	}

	if err != nil {
		h.log.Infof("error while delete quiz: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "sucess")
}
