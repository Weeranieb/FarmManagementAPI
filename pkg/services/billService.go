package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
)

type IBillService interface {
	Create(request models.AddBill, userIdentity string) (*models.Bill, error)
	Get(id int) (*models.Bill, error)
	Update(request *models.Bill, userIdentity string) error
}

type billServiceImp struct {
	BillRepo repositories.IBillRepository
}

func NewBillService(billRepo repositories.IBillRepository) IBillService {
	return &billServiceImp{
		BillRepo: billRepo,
	}
}

func (sv billServiceImp) Create(request models.AddBill, userIdentity string) (*models.Bill, error) {
	var err error
	// validate request
	if err := request.Validation(); err != nil {
		return nil, err
	}

	// check bill if exist
	newBill := &models.Bill{}
	request.Transfer(newBill)
	newBill.UpdatedBy = userIdentity
	newBill.CreatedBy = userIdentity

	// create bill
	newBill, err = sv.BillRepo.Create(newBill)
	if err != nil {
		return nil, err
	}

	return newBill, nil
}

func (sv billServiceImp) Get(id int) (*models.Bill, error) {
	return sv.BillRepo.TakeById(id)
}

func (sv billServiceImp) Update(request *models.Bill, userIdentity string) error {
	// update bill
	request.UpdatedBy = userIdentity
	if err := sv.BillRepo.Update(request); err != nil {
		return err
	}
	return nil
}
