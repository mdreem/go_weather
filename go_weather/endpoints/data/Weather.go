package data

type Weather struct {
	Temperature float64 `json:"temperature"`
	Pressure    int     `json:"pressure"`
	CityName    string  `json:"cityName"`
	CityId      string  `json:"cityId"`
}
