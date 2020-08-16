package repository_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/nikolausreza131192/pos/entity"
	"github.com/nikolausreza131192/pos/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetByUsername(t *testing.T) {
	now := time.Now()
	db, mock, _ := sqlmock.New()
	mockDB := sqlx.NewDb(db, "sqlmock")
	tcs := []struct {
		name           string
		expectedQuery  func()
		ID             int
		expectedResult entity.User
	}{
		{
			name: "Get from DB success",
			expectedQuery: func() {
				rows := sqlmock.NewRows([]string{"id", "nama", "username", "email", "role", "status", "created_by", "updated_by", "created_at", "updated_at"}).AddRow(1, "Wewe", "wewe", "nikolaus.reza@tokopedia.com", "admin", 1, "wewe", "wewe", now, now)
				mock.ExpectQuery("SELECT (.+) FROM m_user").WillReturnRows(rows)
			},
			ID: 2,
			expectedResult: entity.User{
				ID:        1,
				Name:      "Wewe",
				Username:  "wewe",
				Email:     "nikolaus.reza@tokopedia.com",
				Role:      "admin",
				Status:    1,
				CreatedBy: "wewe",
				UpdatedBy: "wewe",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "Get from DB; Error in struct scan",
			expectedQuery: func() {
				rows := sqlmock.NewRows([]string{"x"}).AddRow(1)
				mock.ExpectQuery("SELECT (.+) FROM m_user").WillReturnRows(rows)
			},
			ID:             2,
			expectedResult: entity.User{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectedQuery()
			repoParam := repository.UserRepoParam{
				DB: mockDB,
			}
			repo := repository.NewUser(repoParam)
			result := repo.GetByUsername("admin")

			assert.Nil(t, mock.ExpectationsWereMet())
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGetUserPassword(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mockDB := sqlx.NewDb(db, "sqlmock")
	tcs := []struct {
		name           string
		expectedQuery  func()
		ID             int
		expectedResult string
		expectedError  error
	}{
		{
			name: "Get from DB success",
			expectedQuery: func() {
				rows := sqlmock.NewRows([]string{"password"}).AddRow("secrethashedpassword")
				mock.ExpectQuery("SELECT (.+) FROM m_user").WillReturnRows(rows)
			},
			ID:             2,
			expectedResult: "secrethashedpassword",
		},
		{
			name: "User not found in database",
			expectedQuery: func() {
				mock.ExpectQuery("SELECT (.+) FROM m_user").WillReturnError(sql.ErrNoRows)
			},
			ID:             2,
			expectedResult: "",
			expectedError:  errors.New("User is not found"),
		},
		{
			name: "Error database",
			expectedQuery: func() {
				mock.ExpectQuery("SELECT (.+) FROM m_user").WillReturnError(errors.New("Some error"))
			},
			ID:             2,
			expectedResult: "",
			expectedError:  errors.New("Some error"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectedQuery()
			repoParam := repository.UserRepoParam{
				DB: mockDB,
			}
			repo := repository.NewUser(repoParam)
			result, err := repo.GetUserPassword("admin")

			assert.Nil(t, mock.ExpectationsWereMet())
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCreateUser(t *testing.T) {
	now := time.Now()
	db, mock, _ := sqlmock.New()
	mockDB := sqlx.NewDb(db, "sqlmock")
	tcs := []struct {
		name               string
		dbTransExpectation func() error
		expectedQuery      func()
		user               entity.User
		expectError        bool
		expectUserPassword bool
	}{
		{
			name:          "Error on begin transaction",
			expectedQuery: func() {},
			user: entity.User{
				Name:      "Test",
				Username:  "test",
				Email:     "test@mail.com",
				Role:      "admin",
				CreatedBy: "Foo Bar",
				UpdatedBy: "Foo Bar",
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectError: true,
		},
		{
			name: "Error query",
			expectedQuery: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO m_user`).WillReturnError(errors.New("Some error"))
			},
			user: entity.User{
				Name:      "Test",
				Username:  "test",
				Email:     "test@mail.com",
				Role:      "admin",
				CreatedBy: "Foo Bar",
				UpdatedBy: "Foo Bar",
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectError: true,
		},
		{
			name: "Error commit transaction",
			expectedQuery: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO m_user`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			user: entity.User{
				Name:      "Test",
				Username:  "test",
				Email:     "test@mail.com",
				Role:      "admin",
				CreatedBy: "Foo Bar",
				UpdatedBy: "Foo Bar",
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectError: true,
		},
		{
			name: "Successfully insert to DB",
			expectedQuery: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO m_user`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			user: entity.User{
				Name:      "Test",
				Username:  "test",
				Email:     "test@mail.com",
				Role:      "admin",
				CreatedBy: "Foo Bar",
				UpdatedBy: "Foo Bar",
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectError: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectedQuery()
			repoParam := repository.UserRepoParam{
				DB: mockDB,
			}
			repo := repository.NewUser(repoParam)
			pass, err := repo.CreateUser(tc.user)

			assert.Nil(t, mock.ExpectationsWereMet())
			if tc.expectError {
				assert.NotNil(t, err)
				assert.Equal(t, "", pass)
			} else {
				assert.Nil(t, err)
				assert.NotEqual(t, "", pass)
			}
		})
	}
}
