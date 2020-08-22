package controllers_test

import (
	"github.com/golang/mock/gomock"
	"testing"
	"time"

	"github.com/nikolausreza131192/pos/controllers"
	"github.com/nikolausreza131192/pos/entity"
	"github.com/nikolausreza131192/pos/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	now := time.Now()
	tcs := []struct {
		name            string
		itemRP          repository.ItemRepo
		mockItemRP      func() repository.ItemRepo
		expectedResults []entity.Item
	}{
		{
			name: "Empty result",
			mockItemRP: func() repository.ItemRepo {
				itemRepo := repository.NewMockItemRepo(mockController)

				itemRepo.EXPECT().GetAll().Return([]entity.Item{})

				return itemRepo
			},
			expectedResults: []entity.Item{},
		},
		{
			name: "Result is not empty",
			mockItemRP: func() repository.ItemRepo {
				itemRepo := repository.NewMockItemRepo(mockController)

				response := []entity.Item{
					{
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
				}
				itemRepo.EXPECT().GetAll().Return(response)

				return itemRepo
			},
			expectedResults: []entity.Item{
				{
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
				ItemRP: tc.mockItemRP(),
			}
			controller := controllers.NewItem(param)
			results := controller.GetAll()

			assert.Equal(t, tc.expectedResults, results)
		})
	}
}

func TestGetByID(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	now := time.Now()
	tcs := []struct {
		name            string
		mockItemRP      func() repository.ItemRepo
		expectedResults entity.Item
	}{
		{
			name: "Empty result",
			mockItemRP: func() repository.ItemRepo {
				itemRepo := repository.NewMockItemRepo(mockController)

				itemRepo.EXPECT().GetByID(1).Return(entity.Item{})

				return itemRepo
			},
			expectedResults: entity.Item{},
		},
		{
			name: "Result is not empty",
			mockItemRP: func() repository.ItemRepo {
				itemRepo := repository.NewMockItemRepo(mockController)

				response := entity.Item{
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
				}
				itemRepo.EXPECT().GetByID(1).Return(response)

				return itemRepo
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
				ItemRP: tc.mockItemRP(),
			}
			controller := controllers.NewItem(param)
			results := controller.GetByID(1)

			assert.Equal(t, tc.expectedResults, results)
		})
	}
}
