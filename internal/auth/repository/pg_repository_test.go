package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/lib/tracer"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthRepo_CreateUser(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewRepository(sqlxDB, tracer.InitTracer("main", "test"))

	t.Run("CreateUser", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(
			"1")

		user := domain.SignUpRequest{
			Username: "username",
			Email:    "email@gmail.com",
			Password: "password",
		}

		mock.ExpectQuery("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id").WithArgs(&user.Username, &user.Email,
			&user.Password).WillReturnRows(rows)

		createdUser, err := authRepo.CreateUser(context.Background(), user)

		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, 1)
	})
}

func TestAuthRepo_ValidateUser(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewRepository(sqlxDB, tracer.InitTracer("main", "test"))

	t.Run("ValidateUser", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "avatar"}).AddRow(
			"1", "username", "email@gmail.com", "password", "")

		user := domain.SignInRequest{
			Email:    "email@gmail.com",
			Password: "password",
		}

		mock.ExpectQuery("SELECT * FROM users WHERE email = $1 AND password = $2").WithArgs(&user.Email,
			&user.Password).WillReturnRows(rows)

		validatedUser, err := authRepo.ValidateUser(context.Background(), user)

		require.NoError(t, err)
		require.NotNil(t, validatedUser)
		require.Equal(t, validatedUser.Email, user.Email)
	})
}
