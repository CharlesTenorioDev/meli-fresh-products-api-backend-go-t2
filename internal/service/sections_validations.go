package service

type SectionsValidation struct{}

func (v *SectionsValidation) WarehouseExistsById(id int) bool {
	return true
}

func (v *SectionsValidation) ProductTypeExistsById(id int) bool {
	return true
}
