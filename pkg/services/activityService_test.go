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
		expectedReturn *models.ActivityWithSellDetail
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
			expectedReturn: &models.ActivityWithSellDetail{
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
			expectedReturn: &models.ActivityWithSellDetail{
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
			expectedReturn: &models.ActivityWithSellDetail{
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
			expectedError:  errors.New("the activity already exist on the given date"),
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

func TestGetActivity(t *testing.T) {
	var (
		mockActivityRepo   *mocks.IActivityRepository
		mockSellDetailRepo *mocks.ISellDetailRepository
		activityService    services.IActivityService
	)

	beforeEach := func() {
		mockActivityRepo = new(mocks.IActivityRepository)
		mockSellDetailRepo = new(mocks.ISellDetailRepository)
		activityService = services.NewActivityService(mockActivityRepo, mockSellDetailRepo)
	}

	t.Run("Get activity success", func(t *testing.T) {
		beforeEach()

		mockActivityRepo.On("TakeById", 1).Return(&models.Activity{
			Id:           1,
			ActivePondId: 1,
			Mode:         "SELL",
			ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			MerchantId:   indirect.Int(1),
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}, nil)

		mockSellDetailRepo.On("ListByQuery", "\"SellId\" = ? AND \"DelFlag\" = ?", 1, false).Return([]models.SellDetail{
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
		result, err := activityService.Get(1)

		assert.Nil(t, err)
		assert.Equal(t, &models.ActivityWithSellDetail{
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
		}, result)
	})

	t.Run("Should return error when ListByQuery error", func(t *testing.T) {
		beforeEach()

		mockActivityRepo.On("TakeById", 1).Return(&models.Activity{
			Id:           1,
			ActivePondId: 1,
			Mode:         "SELL",
			ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			MerchantId:   indirect.Int(1),
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}, nil)

		mockSellDetailRepo.On("ListByQuery", "\"SellId\" = ? AND \"DelFlag\" = ?", 1, false).Return([]models.SellDetail{}, nil)
		result, err := activityService.Get(1)

		assert.Error(t, err)
		assert.Equal(t, "sell detail not found", err.Error())
		assert.Nil(t, result)
	})

	t.Run("Should return error when ListByQuery eror", func(t *testing.T) {
		beforeEach()

		mockActivityRepo.On("TakeById", 1).Return(&models.Activity{
			Id:           1,
			ActivePondId: 1,
			Mode:         "SELL",
			ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			MerchantId:   indirect.Int(1),
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}, nil)

		mockSellDetailRepo.On("ListByQuery", "\"SellId\" = ? AND \"DelFlag\" = ?", 1, false).Return(nil, assert.AnError)
		result, err := activityService.Get(1)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("Should return error when TakeById error", func(t *testing.T) {
		beforeEach()

		mockActivityRepo.On("TakeById", 1).Return(nil, assert.AnError)
		result, err := activityService.Get(1)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}

func TestUpdateActivity(t *testing.T) {
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

	t.Run("Update activity success", func(t *testing.T) {
		beforeEach()

		mockActivityRepo.On("Update", &models.Activity{
			Id:           1,
			ActivePondId: 1,
			Mode:         "SELL",
			ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			MerchantId:   indirect.Int(1),
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}).Return(nil)
		mockActivityRepo.On("WithTrx", mock.Anything).Return(mockActivityRepo)
		mockSellDetailRepo.On("WithTrx", mock.Anything).Return(mockSellDetailRepo)

		mockSellDetailRepo.On("Update", &models.SellDetail{
			Id:           1,
			SellId:       1,
			FishType:     "Kraphong",
			Size:         "M",
			Amount:       100,
			FishUnit:     "Kilogram",
			PricePerUnit: 100,
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}).Return(nil)

		mockSellDetailRepo.On("Create", &models.SellDetail{
			Id:           0,
			SellId:       1,
			FishType:     "Nil",
			Size:         "M",
			Amount:       100,
			FishUnit:     "Kilogram",
			PricePerUnit: 100,
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}).Return(&models.SellDetail{
			Id:           2,
			SellId:       1,
			FishType:     "Nil",
			Size:         "M",
			Amount:       100,
			FishUnit:     "Kilogram",
			PricePerUnit: 100,
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}, nil)

		result, err := activityService.Update(&models.ActivityWithSellDetail{
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
					Id:           0,
					SellId:       1,
					FishType:     "Nil",
					Size:         "M",
					Amount:       100,
					FishUnit:     "Kilogram",
					PricePerUnit: 100,
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				},
			},
		}, "testUser")

		assert.Nil(t, err)
		assert.Equal(t, []*models.SellDetail{
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
		}, result)
	})

	t.Run("Should return error with create", func(t *testing.T) {
		beforeEach()

		mockActivityRepo.On("Update", &models.Activity{
			Id:           1,
			ActivePondId: 1,
			Mode:         "SELL",
			ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			MerchantId:   indirect.Int(1),
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}).Return(nil)
		mockActivityRepo.On("WithTrx", mock.Anything).Return(mockActivityRepo)
		mockSellDetailRepo.On("WithTrx", mock.Anything).Return(mockSellDetailRepo)

		mockSellDetailRepo.On("Create", &models.SellDetail{
			Id:           0,
			SellId:       1,
			FishType:     "Nil",
			Size:         "M",
			Amount:       100,
			FishUnit:     "Kilogram",
			PricePerUnit: 100,
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}).Return(nil, assert.AnError)

		result, err := activityService.Update(&models.ActivityWithSellDetail{
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
					Id:           0,
					SellId:       1,
					FishType:     "Nil",
					Size:         "M",
					Amount:       100,
					FishUnit:     "Kilogram",
					PricePerUnit: 100,
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				},
			},
		}, "testUser")

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.Nil(t, result)
	})

	t.Run("Should return error with update", func(t *testing.T) {
		beforeEach()

		mockActivityRepo.On("Update", &models.Activity{
			Id:           1,
			ActivePondId: 1,
			Mode:         "SELL",
			ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			MerchantId:   indirect.Int(1),
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}).Return(nil)
		mockActivityRepo.On("WithTrx", mock.Anything).Return(mockActivityRepo)
		mockSellDetailRepo.On("WithTrx", mock.Anything).Return(mockSellDetailRepo)

		mockSellDetailRepo.On("Update", &models.SellDetail{
			Id:           1,
			SellId:       1,
			FishType:     "Kraphong",
			Size:         "M",
			Amount:       100,
			FishUnit:     "Kilogram",
			PricePerUnit: 100,
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}).Return(assert.AnError)

		result, err := activityService.Update(&models.ActivityWithSellDetail{
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
					Id:           0,
					SellId:       1,
					FishType:     "Nil",
					Size:         "M",
					Amount:       100,
					FishUnit:     "Kilogram",
					PricePerUnit: 100,
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				},
			},
		}, "testUser")

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.Nil(t, result)
	})

	t.Run("Should return error with update activity", func(t *testing.T) {
		beforeEach()

		mockActivityRepo.On("Update", &models.Activity{
			Id:           1,
			ActivePondId: 1,
			Mode:         "SELL",
			ActivityDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			MerchantId:   indirect.Int(1),
			Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
		}).Return(assert.AnError)
		mockActivityRepo.On("WithTrx", mock.Anything).Return(mockActivityRepo)

		result, err := activityService.Update(&models.ActivityWithSellDetail{
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
					Id:           0,
					SellId:       1,
					FishType:     "Nil",
					Size:         "M",
					Amount:       100,
					FishUnit:     "Kilogram",
					PricePerUnit: 100,
					Base:         models.Base{CreatedBy: "testUser", UpdatedBy: "testUser"},
				},
			},
		}, "testUser")

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.Nil(t, result)
	})
}
