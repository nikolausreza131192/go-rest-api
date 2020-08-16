package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nikolausreza131192/pos/entity"
	"github.com/nikolausreza131192/pos/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	now := time.Now()
	db, mock, _ := sqlmock.New()
	mockDB := sqlx.NewDb(db, "sqlmock")
	tcs := []struct {
		name            string
		expectedQuery   func()
		expectedResults []entity.Item
	}{
		{
			name: "Query error",
			expectedQuery: func() {
				mock.ExpectQuery("SELECT (.+) FROM m_barang").WillReturnError(errors.New("Foo Bar"))
			},
			expectedResults: []entity.Item{},
		},
		{
			name: "Struct scan failed",
			expectedQuery: func() {
				rows := sqlmock.NewRows([]string{"unknown_column"}).AddRow("xxx")
				mock.ExpectQuery("SELECT (.+) FROM m_barang").WillReturnRows(rows)
			},
			expectedResults: []entity.Item{},
		},
		{
			name: "Success",
			expectedQuery: func() {
				rows := sqlmock.NewRows([]string{"id", "kategori_id", "merk_id", "kode", "nama", "panjang", "lebar", "initial_stok", "initial_m2", "keterangan", "created_by", "updated_by", "created_at", "updated_at"}).AddRow(1, 1, 1, "x", "xx", 1, 1, 1, 1, "", "wewe", "wewe", now, now)
				mock.ExpectQuery("SELECT (.+) FROM m_barang").WillReturnRows(rows)
			},
			expectedResults: []entity.Item{
				entity.Item{
					ID:           1,
					CategoryID:   1,
					BrandID:      1,
					Code:         "x",
					Name:         "xx",
					Length:       1,
					Width:        1,
					InitialStock: 1,
					InitialArea:  1,
					CreatedBy:    "wewe",
					UpdatedBy:    "wewe",
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectedQuery()
			repoParam := repository.ItemRepoParam{
				DB: mockDB,
			}
			repo := repository.NewItem(repoParam)
			results := repo.GetAll()

			assert.Nil(t, mock.ExpectationsWereMet())
			assert.Equal(t, tc.expectedResults, results)
		})
	}
}

func TestGetByID(t *testing.T) {
	now := time.Now()
	db, mock, _ := sqlmock.New()
	mockDB := sqlx.NewDb(db, "sqlmock")
	tcs := []struct {
		name           string
		expectedQuery  func()
		ID             int
		expectedResult entity.Item
	}{
		{
			name: "Get from cache",
			expectedQuery: func() {
				rows := sqlmock.NewRows([]string{"id", "kategori_id", "merk_id", "kode", "nama", "panjang", "lebar", "initial_stok", "initial_m2", "keterangan", "created_by", "updated_by", "created_at", "updated_at"}).AddRow(1, 1, 1, "x", "xx", 1, 1, 1, 1, "", "wewe", "wewe", now, now)
				mock.ExpectQuery("SELECT (.+) FROM m_barang").WillReturnRows(rows)
			},
			ID: 1,
			expectedResult: entity.Item{
				ID:           1,
				CategoryID:   1,
				BrandID:      1,
				Code:         "x",
				Name:         "xx",
				Length:       1,
				Width:        1,
				InitialStock: 1,
				InitialArea:  1,
				CreatedBy:    "wewe",
				UpdatedBy:    "wewe",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
		{
			name: "Not found in cache; Get from DB, but get error in struct scan",
			expectedQuery: func() {
				rows := sqlmock.NewRows([]string{"id", "kategori_id", "merk_id", "kode", "nama", "panjang", "lebar", "initial_stok", "initial_m2", "keterangan", "created_by", "updated_by", "created_at", "updated_at"}).AddRow(1, 1, 1, "x", "xx", 1, 1, 1, 1, "", "wewe", "wewe", now, now)
				mock.ExpectQuery("SELECT (.+) FROM m_barang").WillReturnRows(rows)
				mock.ExpectQuery(`SELECT (.+) FROM m_barang`).WithArgs(2).WillReturnError(errors.New("Foo Bar"))
			},
			ID:             2,
			expectedResult: entity.Item{},
		},
		{
			name: "Not found in cache; Get from DB",
			expectedQuery: func() {
				rows := sqlmock.NewRows([]string{"id", "kategori_id", "merk_id", "kode", "nama", "panjang", "lebar", "initial_stok", "initial_m2", "keterangan", "created_by", "updated_by", "created_at", "updated_at"}).AddRow(1, 1, 1, "x", "xx", 1, 1, 1, 1, "", "wewe", "wewe", now, now)
				mock.ExpectQuery("SELECT (.+) FROM m_barang").WillReturnRows(rows)
				row2 := sqlmock.NewRows([]string{"id", "kategori_id", "merk_id", "kode", "nama", "panjang", "lebar", "initial_stok", "initial_m2", "keterangan", "created_by", "updated_by", "created_at", "updated_at"}).AddRow(2, 1, 1, "y", "yy", 1, 1, 1, 1, "", "wewe", "wewe", now, now)
				mock.ExpectQuery(`SELECT (.+) FROM m_barang`).WithArgs(2).WillReturnRows(row2)
			},
			ID: 2,
			expectedResult: entity.Item{
				ID:           2,
				CategoryID:   1,
				BrandID:      1,
				Code:         "y",
				Name:         "yy",
				Length:       1,
				Width:        1,
				InitialStock: 1,
				InitialArea:  1,
				CreatedBy:    "wewe",
				UpdatedBy:    "wewe",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectedQuery()
			repoParam := repository.ItemRepoParam{
				DB: mockDB,
			}
			repo := repository.NewItem(repoParam)
			result := repo.GetByID(tc.ID)

			assert.Nil(t, mock.ExpectationsWereMet())
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
