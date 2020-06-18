package openweather

type WeatherResponse struct {
	Main MainData `json:"main"`
	Name string   `json:"name"`
	Id   int64    `json:"id"`
}

type MainData struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  float64 `json:"pressure"`
	Humidity  float64 `json:"humidity"`
}
