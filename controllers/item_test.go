package controllers_test

import (
	"testing"
	"time"

	"github.com/nikolausreza131192/pos/controllers"
	"github.com/nikolausreza131192/pos/entity"
	"github.com/nikolausreza131192/pos/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	now := time.Now()
	tcs := []struct {
		name            string
		itemRP          repository.ItemRepo
		expectedResults []entity.Item
	}{
		{
			name: "Empty result",
			itemRP: &fakeItemRepo{
				GetAllResults: []entity.Item{},
			},
			expectedResults: []entity.Item{},
		},
		{
			name: "Result is not empty",
			itemRP: &fakeItemRepo{
				GetAllResults: []entity.Item{
					entity.Item{
						CategoryID:   1,
						BrandID:      1,
						Code:         "X",
						Name:         "Item 1",
						Length:       10,
						Width:        10,
						InitialStock: 1,
						InitialArea:  100,
						Remark:       "Lorem ipsum",
						CreatedBy:    "Wewe",
						UpdatedBy:    "Wewe",
						CreatedAt:    now,
						UpdatedAt:    now,
					},
				},
			},
			expectedResults: []entity.Item{
				entity.Item{
					CategoryID:   1,
					BrandID:      1,
					Code:         "X",
					Name:         "Item 1",
					Length:       10,
					Width:        10,
					InitialStock: 1,
					InitialArea:  100,
					Remark:       "Lorem ipsum",
					CreatedBy:    "Wewe",
					UpdatedBy:    "Wewe",
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			param := controllers.ItemControllerParam{
				ItemRP: tc.itemRP,
			}
			controller := controllers.NewItem(param)
			results := controller.GetAll()

			assert.Equal(t, tc.expectedResults, results)
		})
	}
}

func TestGetByID(t *testing.T) {
	now := time.Now()
	tcs := []struct {
		name            string
		itemRP          repository.ItemRepo
		expectedResults entity.Item
	}{
		{
			name: "Empty result",
			itemRP: &fakeItemRepo{
				GetByIDResult: entity.Item{},
			},
			expectedResults: entity.Item{},
		},
		{
			name: "Result is not empty",
			itemRP: &fakeItemRepo{
				GetByIDResult: entity.Item{
					CategoryID:   1,
					BrandID:      1,
					Code:         "X",
					Name:         "Item 1",
					Length:       10,
					Width:        10,
					InitialStock: 1,
					InitialArea:  100,
					Remark:       "Lorem ipsum",
					CreatedBy:    "Wewe",
					UpdatedBy:    "Wewe",
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
			expectedResults: entity.Item{
				CategoryID:   1,
				BrandID:      1,
				Code:         "X",
				Name:         "Item 1",
				Length:       10,
				Width:        10,
				InitialStock: 1,
				InitialArea:  100,
				Remark:       "Lorem ipsum",
				CreatedBy:    "Wewe",
				UpdatedBy:    "Wewe",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			param := controllers.ItemControllerParam{
				ItemRP: tc.itemRP,
			}
			controller := controllers.NewItem(param)
			results := controller.GetByID(1)

			assert.Equal(t, tc.expectedResults, results)
		})
	}
}
