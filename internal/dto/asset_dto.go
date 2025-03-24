package dto

type AssetInputDto struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Value           float64 `json:"value"`
	AcquisitionDate string  `json:"acquisition_date"`
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
