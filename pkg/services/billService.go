package services

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/repositories"
	"boonmafarm/api/utils/httputil"
)

type IBillService interface {
	Create(request models.AddBill, userIdentity string) (*models.Bill, error)
	Get(id int) (*models.Bill, error)
	Update(request *models.Bill, userIdentity string) error
	TakePage(clientId, page, pageSize int, orderBy, keyword string, billType string, farmGroupId int) (*httputil.PageModel, error)
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

func (sv billServiceImp) TakePage(clientId, page, pageSize int, orderBy, keyword string, billType string, farmGroupId int) (*httputil.PageModel, error) {
	result := &httputil.PageModel{}
	var typePointer *string
	var farmGroupIdPointer *int

	if billType != "" {
		typePointer = &billType
	}

	if farmGroupId != 0 {
		farmGroupIdPointer = &farmGroupId
	}

	items, total, err := sv.BillRepo.TakePage(clientId, page, pageSize, orderBy, keyword, typePointer, farmGroupIdPointer)
	if err != nil {
		return nil, err
	}

	result.Items = items
	result.Total = total

	return result, nil
}
