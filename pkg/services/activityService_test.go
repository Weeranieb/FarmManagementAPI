package services_test

import (
	dbContext "boonmafarm/api/pkg/dbcontext"
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"boonmafarm/api/utils/indirect"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres" // Change the driver to postgres
	"gorm.io/gorm"
)

func TestCreateActivity(t *testing.T) {

	var (
		mockActivityRepo   *mocks.IActivityRepository
		mockSellDetailRepo *mocks.ISellDetailRepository
		activityService    services.IActivityService
		db                 *gorm.DB
	)

	beforeEach := func() {
		mockActivityRepo = new(mocks.IActivityRepository)
		mockSellDetailRepo = new(mocks.ISellDetailRepository)
		activityService = services.NewActivityService(mockActivityRepo, mockSellDetailRepo)

		mockDB, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer mockDB.Close()

		// Create a gorm.DB instance from the mockDB using the postgres driver
		dialector := postgres.New(postgres.Config{
			Conn:       mockDB,
			DriverName: "postgres",
		})

		db, err = gorm.Open(dialector, &gorm.Config{})
		if err != nil {
			t.Fatalf("An error occurred while creating gorm.DB: %s", err)
		}

		// Inject the db instance into the dbcontext
		dbContext.Context.Postgresql = db
	}

	tests := []struct {
		name           string
		request        models.CreateActivityRequest
		mockService    func()
		expectedReturn *models.CreateActivityResponse
		expectedError  error
	}{
		{
			name: "Create activity fill success",
			request: models.CreateActivityRequest{
				ActivePondId: 1,
				Mode:         "FILL",
				ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Amount:       indirect.Int(100),
				FishUnit:     indirect.String("Kilogram"),
				FishType:     indirect.String("Kraphong"),
				FishWeight:   indirect.Float64(100),
				PricePerUnit: indirect.Float64(100),
			},
			mockService: func() {
				mockActivityRepo.On("FirstByQuery", "\"Mode\" = ? AND \"ActivityDate\" = ? AND \"DelFlag\" = ?", "FILL", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), false).Return(nil, nil)
				mockActivityRepo.On("WithTrx", mock.Anything).Return(mockActivityRepo)
				mockActivityRepo.On("Create", &models.Activity{
					ActivePondId: 1,
					Mode:         "FILL",
					ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Amount:       indirect.Int(100),
					FishUnit:     indirect.String("Kilogram"),
					FishType:     indirect.String("Kraphong"),
					FishWeight:   indirect.Float64(100),
					PricePerUnit: indirect.Float64(100),
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				}).Return(&models.Activity{
					Id:           1,
					ActivePondId: 1,
					Mode:         "FILL",
					ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Amount:       indirect.Int(100),
					FishUnit:     indirect.String("Kilogram"),
					FishType:     indirect.String("Kraphong"),
					FishWeight:   indirect.Float64(100),
					PricePerUnit: indirect.Float64(100),
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				}, nil)
			},
			expectedReturn: &models.CreateActivityResponse{
				Activity: models.Activity{
					Id:           1,
					ActivePondId: 1,
					Mode:         "FILL",
					ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Amount:       indirect.Int(100),
					FishUnit:     indirect.String("Kilogram"),
					FishType:     indirect.String("Kraphong"),
					FishWeight:   indirect.Float64(100),
					PricePerUnit: indirect.Float64(100),
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				},
			},
			expectedError: nil,
		},
		{
			name: "Create activity move success",
			request: models.CreateActivityRequest{
				ActivePondId:   1,
				ToActivePondId: indirect.Int(2),
				Mode:           "MOVE",
				ActivityDate:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Amount:         indirect.Int(100),
				FishUnit:       indirect.String("Kilogram"),
				FishType:       indirect.String("Kraphong"),
				FishWeight:     indirect.Float64(100),
				PricePerUnit:   indirect.Float64(100),
			},
			mockService: func() {
				mockActivityRepo.On("FirstByQuery", "\"Mode\" = ? AND \"ActivityDate\" = ? AND \"DelFlag\" = ?", "MOVE", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), false).Return(nil, nil)
				mockActivityRepo.On("WithTrx", mock.Anything).Return(mockActivityRepo)
				mockActivityRepo.On("Create", &models.Activity{
					ActivePondId:   1,
					ToActivePondId: indirect.Int(2),
					Mode:           "MOVE",
					ActivityDate:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Amount:         indirect.Int(100),
					FishUnit:       indirect.String("Kilogram"),
					FishType:       indirect.String("Kraphong"),
					FishWeight:     indirect.Float64(100),
					PricePerUnit:   indirect.Float64(100),
					Base:           models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				}).Return(&models.Activity{
					Id:             1,
					ActivePondId:   1,
					ToActivePondId: indirect.Int(2),
					Mode:           "MOVE",
					ActivityDate:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Amount:         indirect.Int(100),
					FishUnit:       indirect.String("Kilogram"),
					FishType:       indirect.String("Kraphong"),
					FishWeight:     indirect.Float64(100),
					PricePerUnit:   indirect.Float64(100),
					Base:           models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				}, nil)
			},
			expectedReturn: &models.CreateActivityResponse{
				Activity: models.Activity{
					Id:             1,
					ActivePondId:   1,
					ToActivePondId: indirect.Int(2),
					Mode:           "MOVE",
					ActivityDate:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Amount:         indirect.Int(100),
					FishUnit:       indirect.String("Kilogram"),
					FishType:       indirect.String("Kraphong"),
					FishWeight:     indirect.Float64(100),
					PricePerUnit:   indirect.Float64(100),
					Base:           models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				},
			},
			expectedError: nil,
		},
		{
			name: "Create activity sell success",
			request: models.CreateActivityRequest{
				ActivePondId: 1,
				Mode:         "SELL",
				ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				MerchantId:   indirect.Int(1),
				SellDetail: []models.AddSellDetail{
					{
						FishType:     "Kraphong",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
					},
					{
						FishType:     "Nil",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
					},
				},
			},
			mockService: func() {
				mockActivityRepo.On("FirstByQuery", "\"Mode\" = ? AND \"ActivityDate\" = ? AND \"DelFlag\" = ?", "SELL", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), false).Return(nil, nil)
				mockActivityRepo.On("WithTrx", mock.Anything).Return(mockActivityRepo)
				mockActivityRepo.On("Create", &models.Activity{
					ActivePondId: 1,
					Mode:         "SELL",
					ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					MerchantId:   indirect.Int(1),
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				}).Return(&models.Activity{
					Id:           1,
					ActivePondId: 1,
					Mode:         "SELL",
					ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					MerchantId:   indirect.Int(1),
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				}, nil)
				mockSellDetailRepo.On("WithTrx", mock.Anything).Return(mockSellDetailRepo)
				mockSellDetailRepo.On("BulkCreate", []models.SellDetail{
					{
						SellId:       1,
						FishType:     "Kraphong",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
						Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
					},
					{
						SellId:       1,
						FishType:     "Nil",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
						Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
					},
				}).Return([]models.SellDetail{
					{
						Id:           1,
						SellId:       1,
						FishType:     "Kraphong",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
						Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
					},
					{
						Id:           2,
						SellId:       1,
						FishType:     "Nil",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
						Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
					},
				}, nil)
			},
			expectedReturn: &models.CreateActivityResponse{
				Activity: models.Activity{
					Id:           1,
					ActivePondId: 1,
					Mode:         "SELL",
					ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					MerchantId:   indirect.Int(1),
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				},
				SellDetail: []models.SellDetail{
					{
						Id:           1,
						SellId:       1,
						FishType:     "Kraphong",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
						Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
					},
					{
						Id:           2,
						SellId:       1,
						FishType:     "Nil",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
						Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Create activity sell fail at BulkCreate",
			request: models.CreateActivityRequest{
				ActivePondId: 1,
				Mode:         "SELL",
				ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				MerchantId:   indirect.Int(1),
				SellDetail: []models.AddSellDetail{
					{
						FishType:     "Kraphong",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
					},
					{
						FishType:     "Nil",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
					},
				},
			},
			mockService: func() {
				mockActivityRepo.On("FirstByQuery", "\"Mode\" = ? AND \"ActivityDate\" = ? AND \"DelFlag\" = ?", "SELL", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), false).Return(nil, nil)
				mockActivityRepo.On("WithTrx", mock.Anything).Return(mockActivityRepo)
				mockActivityRepo.On("Create", &models.Activity{
					ActivePondId: 1,
					Mode:         "SELL",
					ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					MerchantId:   indirect.Int(1),
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				}).Return(&models.Activity{
					Id:           1,
					ActivePondId: 1,
					Mode:         "SELL",
					ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					MerchantId:   indirect.Int(1),
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				}, nil)
				mockSellDetailRepo.On("WithTrx", mock.Anything).Return(mockSellDetailRepo)
				mockSellDetailRepo.On("BulkCreate", []models.SellDetail{
					{
						SellId:       1,
						FishType:     "Kraphong",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
						Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
					},
					{
						SellId:       1,
						FishType:     "Nil",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
						Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
					},
				}).Return(nil, errors.New("failed to bulk create sell details"))
			},
			expectedReturn: nil,
			expectedError:  errors.New("failed to bulk create sell details"),
		},
		{
			name: "Create activity fail at Create",
			request: models.CreateActivityRequest{
				ActivePondId: 1,
				Mode:         "SELL",
				ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				MerchantId:   indirect.Int(1),
				SellDetail: []models.AddSellDetail{
					{
						FishType:     "Kraphong",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
					},
					{
						FishType:     "Nil",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
					},
				},
			},
			mockService: func() {
				mockActivityRepo.On("FirstByQuery", "\"Mode\" = ? AND \"ActivityDate\" = ? AND \"DelFlag\" = ?", "SELL", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), false).Return(nil, nil)
				mockActivityRepo.On("WithTrx", mock.Anything).Return(mockActivityRepo)
				mockActivityRepo.On("Create", &models.Activity{
					ActivePondId: 1,
					Mode:         "SELL",
					ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					MerchantId:   indirect.Int(1),
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				}).Return(nil, assert.AnError)
			},
			expectedReturn: nil,
			expectedError:  assert.AnError,
		},
		{
			name: "Create activity fail at FirstByQuery",
			request: models.CreateActivityRequest{
				ActivePondId: 1,
				Mode:         "SELL",
				ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				MerchantId:   indirect.Int(1),
				SellDetail: []models.AddSellDetail{
					{
						FishType:     "Kraphong",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
					},
				},
			},
			mockService: func() {
				mockActivityRepo.On("FirstByQuery", "\"Mode\" = ? AND \"ActivityDate\" = ? AND \"DelFlag\" = ?", "SELL", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), false).Return(nil, assert.AnError)
			},
			expectedReturn: nil,
			expectedError:  assert.AnError,
		},
		{
			name: "Create activity found at FirstByQuery",
			request: models.CreateActivityRequest{
				ActivePondId: 1,
				Mode:         "SELL",
				ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				MerchantId:   indirect.Int(1),
				SellDetail: []models.AddSellDetail{
					{
						FishType:     "Kraphong",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
					},
				},
			},
			mockService: func() {
				mockActivityRepo.On("FirstByQuery", "\"Mode\" = ? AND \"ActivityDate\" = ? AND \"DelFlag\" = ?", "SELL", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), false).Return(&models.Activity{
					ActivePondId: 1,
					Mode:         "SELL",
					ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					MerchantId:   indirect.Int(1),
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				}, nil)
			},
			expectedReturn: nil,
			expectedError:  assert.AnError,
		},
		{
			name: "Create activity fail at validation",
			request: models.CreateActivityRequest{
				ActivePondId: 1,
				Mode:         "SELL",
				ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				SellDetail: []models.AddSellDetail{
					{
						FishType:     "Kraphong",
						Size:         "M",
						Amount:       100,
						FishUnit:     "Kilogram",
						PricePerUnit: 100,
					},
				},
			},
			expectedReturn: nil,
			expectedError:  errors.New("merchant id is empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeEach()

			if tt.mockService != nil {
				tt.mockService()
			}

			result, err := activityService.Create(tt.request, "testUser")

			if tt.expectedError == nil {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedReturn, result)
			} else {
				assert.Nil(t, result)
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err)
			}
		})
	}
}
