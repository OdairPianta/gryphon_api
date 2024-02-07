package requests

type CreateChartRadar struct {
	BackgroundColor string `json:"backgroundColor"  binding:"required" example:"#ffffff"`
	Width           string `json:"width" example:"900px"`
	Height          string `json:"height" example:"500px"`
	Legend          struct {
		Show bool     `json:"show"  binding:"required" example:"true"`
		Data []string `json:"data"  binding:"required" example:"TARGET,RESULT"`
	} `json:"legend"  binding:"required"`
	MultiSeries []struct {
		Name string `json:"name"  binding:"required" example:"TARGET"`
		Data []struct {
			Name  string    `json:"name"  binding:"required" example:"Level"`
			Value []float64 `json:"value"  binding:"required""`
		} `json:"data"  binding:"required"`
		ItemStyle struct {
			Color string `json:"color"  binding:"required" example:"rgba(47, 85, 220, 0.2)"`
		} `json:"itemStyle"  binding:"required"`
	} `json:"MultiSeries"  binding:"required"`
	Indicators []struct {
		Name  string  `json:"name"  binding:"required" example:"Level 1"`
		Min   float32 `json:"min"  binding:"required" example:"0"`
		Max   float32 `json:"max"  binding:"required" example:"100"`
		Color string  `json:"color"  binding:"required" example:"#000000"`
	} `json:"Indicators"  binding:"required"`
	SplitNumber int `json:"splitNumber"  binding:"required" example:"5"`
	SplitLine   struct {
		Show      bool `json:"show"  binding:"required" example:"true"`
		LineStyle struct {
			Opacity float32 `json:"opacity"  binding:"required" example:"0.5"`
			Color   string  `json:"color"  binding:"required" example:"#595757"`
			Type    string  `json:"type"  binding:"required" example:"solid"`
		} `json:"lineStyle"  binding:"required"`
	} `json:"splitLine"  binding:"required"`
}
