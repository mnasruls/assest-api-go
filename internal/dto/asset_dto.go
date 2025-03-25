package dto

type AssetInputDto struct {
	Name            string  `json:"name" validate:"required"`
	Type            string  `json:"type" validate:"required"`
	Value           float64 `json:"value" validate:"required"`
	AcquisitionDate string  `json:"acquisition_date" validate:"required"`
}

type AssetOutputDto struct {
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Value           float64 `json:"value"`
	AcquisitionDate string  `json:"acquisition_date"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}
