package data

type Weather struct {
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	CityName    string  `json:"cityName"`
	CityId      int64   `json:"cityId"`
}
