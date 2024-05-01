package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBill(t *testing.T) {

	var (
		mockBillRepo *mocks.IBillRepository
		billService  services.IBillService
	)

	beforeEach := func() {
		mockBillRepo = new(mocks.IBillRepository)
		billService = services.NewBillService(mockBillRepo)
	}

	t.Run("Create bill", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddBill{
			Type:        "Electricity",
			FarmGroupId: 1,
			PaidAmount:  1000,
			PaymentDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockBillReturn := &models.Bill{
			Id:          1,
			Type:        "Electricity",
			FarmGroupId: 1,
			PaidAmount:  1000,
			PaymentDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockBillRepo.On("Create", mock.Anything).Return(mockBillReturn, nil)

		// Call the Create method
		bill, err := billService.Create(request, "testUser")

		// Assert the result
		assert.NoError(t, err)
		assert.NotNil(t, bill)
	})

	t.Run("Create bill with invalid data", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddBill{
			Type:        "Electricity",
			FarmGroupId: 1,
			PaidAmount:  1000,
		}

		// Call the Create method
		bill, err := billService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Equal(t, "payment date is empty", err.Error())
		assert.Nil(t, bill)
	})

	t.Run("Create bill with error on create", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := models.AddBill{
			Type:        "Electricity",
			FarmGroupId: 1,
			PaidAmount:  1000,
			PaymentDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		// Mock the necessary repository methods
		mockBillRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		// Call the Create method
		bill, err := billService.Create(request, "testUser")

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, bill)

	})
}

func TestGetBill(t *testing.T) {

	var (
		mockBillRepo *mocks.IBillRepository
		billService  services.IBillService
	)

	beforeEach := func() {
		mockBillRepo = new(mocks.IBillRepository)
		billService = services.NewBillService(mockBillRepo)
	}

	t.Run("Get bill", func(t *testing.T) {
		beforeEach()

		// Define test data
		mockBillReturn := &models.Bill{
			Id:          1,
			Type:        "Electricity",
			FarmGroupId: 1,
			PaidAmount:  1000,
			PaymentDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockBillRepo.On("TakeById", mock.Anything).Return(mockBillReturn, nil)

		// Call the Get method
		bill, err := billService.Get(1)

		// Assert the result
		assert.NoError(t, err)
		assert.NotNil(t, bill)
	})

	t.Run("Get bill with error on get", func(t *testing.T) {
		beforeEach()

		// Mock the necessary repository methods
		mockBillRepo.On("TakeById", mock.Anything).Return(nil, assert.AnError)

		// Call the Get method
		bill, err := billService.Get(1)

		// Assert the result
		assert.Error(t, err)
		assert.Nil(t, bill)
	})
}

func TestUpdateBill(t *testing.T) {

	var (
		mockBillRepo *mocks.IBillRepository
		billService  services.IBillService
	)

	beforeEach := func() {
		mockBillRepo = new(mocks.IBillRepository)
		billService = services.NewBillService(mockBillRepo)
	}

	t.Run("Update bill", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := &models.Bill{
			Id:          1,
			Type:        "Electricity",
			FarmGroupId: 1,
			PaidAmount:  1000,
			PaymentDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockBillRepo.On("Update", mock.Anything).Return(nil)

		// Call the Update method
		err := billService.Update(request, "testUser")

		// Assert the result
		assert.NoError(t, err)
	})

	t.Run("Update bill with error on update", func(t *testing.T) {
		beforeEach()

		// Define test data
		request := &models.Bill{
			Id:          1,
			Type:        "Electricity",
			FarmGroupId: 1,
			PaidAmount:  1000,
			PaymentDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Base: models.Base{
				DelFlag:   false,
				CreatedBy: "testUser",
				UpdatedBy: "testUser",
			},
		}

		// Mock the necessary repository methods
		mockBillRepo.On("Update", mock.Anything).Return(assert.AnError)

		// Call the Update method
		err := billService.Update(request, "testUser")

		// Assert the result
		assert.Error(t, err)
	})
}
