package constants

type BillType string

const (
	BillTypeWorker      BillType = "WORKER"
	BillTypeElectricity BillType = "ELECTRICITY"
	BillTypeOther       BillType = "OTHER"
)
