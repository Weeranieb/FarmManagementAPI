package services_test

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories/mocks"
	"boonmafarm/api/pkg/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_Merchant(t *testing.T) {
	var (
		mockMerchantRepo *mocks.IMerchantRepository
		mockMerchant     services.IMerchantService
	)

	beforeEach := func() {
		mockMerchantRepo = new(mocks.IMerchantRepository)
		mockMerchant = services.NewMerchantService(mockMerchantRepo)
	}

	t.Run("Create merchant success", func(t *testing.T) {
		beforeEach()

		mockMerchantRepo.On("FirstByQuery", "\"ContactNumber\" = ? AND \"Name\" = ? AND \"DelFlag\" = ?", "contact", "name", false).Return(nil, nil)
		mockMerchantRepo.On("Create", mock.Anything).Return(&models.Merchant{
			Id:            1,
			Name:          "name",
			ContactNumber: "contact",
		}, nil)

		result, err := mockMerchant.Create(models.AddMerchant{
			Name:          "name",
			ContactNumber: "contact",
		}, "test")

		mockMerchantRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Create merchant failed by validation", func(t *testing.T) {
		beforeEach()

		result, err := mockMerchant.Create(models.AddMerchant{
			Name:          "",
			ContactNumber: "",
		}, "test")

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create merchant failed by merchant already exist", func(t *testing.T) {
		beforeEach()

		mockMerchantRepo.On("FirstByQuery", "\"ContactNumber\" = ? AND \"Name\" = ? AND \"DelFlag\" = ?", "contact", "name", false).Return(&models.Merchant{
			Id:            1,
			Name:          "name",
			ContactNumber: "contact",
		}, nil)

		result, err := mockMerchant.Create(models.AddMerchant{
			Name:          "name",
			ContactNumber: "contact",
		}, "test")

		mockMerchantRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create merchant failed by merchant already exist", func(t *testing.T) {
		beforeEach()

		mockMerchantRepo.On("FirstByQuery", "\"ContactNumber\" = ? AND \"Name\" = ? AND \"DelFlag\" = ?", "contact", "name", false).Return(&models.Merchant{
			Id:            1,
			Name:          "name",
			ContactNumber: "contact",
		}, nil)

		result, err := mockMerchant.Create(models.AddMerchant{
			Name:          "name",
			ContactNumber: "contact",
		}, "test")

		mockMerchantRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create merchant failed by repository error by FirstByQuery", func(t *testing.T) {
		beforeEach()

		mockMerchantRepo.On("FirstByQuery", "\"ContactNumber\" = ? AND \"Name\" = ? AND \"DelFlag\" = ?", "contact", "name", false).Return(nil, assert.AnError)

		result, err := mockMerchant.Create(models.AddMerchant{
			Name:          "name",
			ContactNumber: "contact",
		}, "test")

		mockMerchantRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Create merchant failed by repository error by Create", func(t *testing.T) {
		beforeEach()

		mockMerchantRepo.On("FirstByQuery", "\"ContactNumber\" = ? AND \"Name\" = ? AND \"DelFlag\" = ?", "contact", "name", false).Return(nil, nil)
		mockMerchantRepo.On("Create", mock.Anything).Return(nil, assert.AnError)

		result, err := mockMerchant.Create(models.AddMerchant{
			Name:          "name",
			ContactNumber: "contact",
		}, "test")

		mockMerchantRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestGet_Merchant(t *testing.T) {
	var (
		mockMerchantRepo *mocks.IMerchantRepository
		mockMerchant     services.IMerchantService
	)

	beforeEach := func() {
		mockMerchantRepo = new(mocks.IMerchantRepository)
		mockMerchant = services.NewMerchantService(mockMerchantRepo)
	}

	t.Run("Get merchant success", func(t *testing.T) {
		beforeEach()

		mockMerchantRepo.On("TakeById", 1).Return(
			&models.Merchant{
				Id:            1,
				Name:          "name",
				ContactNumber: "contact",
			}, nil)

		result, err := mockMerchant.Get(1)

		mockMerchantRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Get merchant failed by repository error", func(t *testing.T) {
		beforeEach()

		mockMerchantRepo.On("TakeById", 1).Return(
			nil, assert.AnError)

		result, err := mockMerchant.Get(1)

		mockMerchantRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdate_Merchant(t *testing.T) {
	var (
		mockMerchantRepo *mocks.IMerchantRepository
		mockMerchant     services.IMerchantService
	)

	beforeEach := func() {
		mockMerchantRepo = new(mocks.IMerchantRepository)
		mockMerchant = services.NewMerchantService(mockMerchantRepo)
	}

	t.Run("Update merchant success", func(t *testing.T) {
		beforeEach()

		mockMerchantRepo.On("Update", mock.Anything).Return(nil)

		err := mockMerchant.Update(&models.Merchant{
			Id:            1,
			Name:          "name",
			ContactNumber: "contact",
		}, "test")

		mockMerchantRepo.AssertExpectations(t)
		assert.Nil(t, err)
	})

	t.Run("Update merchant failed by repository error", func(t *testing.T) {
		beforeEach()

		mockMerchantRepo.On("Update", mock.Anything).Return(assert.AnError)

		err := mockMerchant.Update(&models.Merchant{
			Id:            1,
			Name:          "name",
			ContactNumber: "contact",
		}, "test")

		mockMerchantRepo.AssertExpectations(t)
		assert.NotNil(t, err)
	})
}
