package controllers

//go:generate go run main.go

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/OdairPianta/gryphon_api/http/requests"
	"github.com/OdairPianta/gryphon_api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/snapshot-chromedp/render"
)

// @Summary		Generate a radar chart
// @Description	Generate a radar chart
// @Tags			charts
// @Accept		json
// @Produce		json
// @Param			chart	body		requests.CreateChartRadar	true	"Request body to generate a radar chart"
// @Success		200		{object}	models.File
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/charts/radar/create [post]
// @Security Bearer
func CreateChartRadar(context *gin.Context) {
	// Validate input
	var input requests.CreateChartRadar
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user, error = models.GetUserByToken(context)
	if error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Erro ao obter usuário!"})
		return
	}

	var indicators = make([]*opts.Indicator, 0)
	for i := 0; i < len(input.Indicators); i++ {
		var indicator = input.Indicators[i]
		indicators = append(indicators, &opts.Indicator{Name: indicator.Name, Min: indicator.Min, Max: indicator.Max, Color: indicator.Color})
	}

	radar := charts.NewRadar()
	radar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			BackgroundColor: input.BackgroundColor,
			Width:           input.Width,
			Height:          input.Height,
		}),
		charts.WithRadarComponentOpts(opts.RadarComponent{
			Indicator:   indicators,
			SplitNumber: input.SplitNumber,
			SplitLine: &opts.SplitLine{
				Show: input.SplitLine.Show,
				LineStyle: &opts.LineStyle{
					Opacity: input.SplitLine.LineStyle.Opacity,
					Color:   input.SplitLine.LineStyle.Color,
					Type:    input.SplitLine.LineStyle.Type,
				},
			},
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: input.Legend.Show,
			Data: input.Legend.Data,
		}),
	)

	for i := 0; i < len(input.MultiSeries); i++ {
		var multiSerie = input.MultiSeries[i]
		items := make([]opts.RadarData, 0)
		for j := 0; j < len(multiSerie.Data); j++ {
			var dataItem = multiSerie.Data[j]
			items = append(items, opts.RadarData{Value: dataItem.Value, Name: dataItem.Name})
		}
		radar.AddSeries(multiSerie.Name, items, charts.WithItemStyleOpts(opts.ItemStyle{Color: multiSerie.ItemStyle.Color}))
	}

	radar.SetSeriesOptions(
		charts.WithLineStyleOpts(opts.LineStyle{}),
		charts.WithAreaStyleOpts(opts.AreaStyle{}),
	)

	var filePathHtml = models.GenerateRandonFileName("html")
	var filePathPng = models.GenerateRandonFileName("png")

	f, err := os.Create("temp/" + filePathHtml)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating file: " + "temp/" + filePathHtml + " " + err.Error()})
		return
	}
	radar.Render(f)
	readFile, err := os.ReadFile("temp/" + filePathHtml)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading file: " + "temp/" + filePathHtml + " " + err.Error()})
		return
	}

	err = makeChartSnapshot(readFile, "temp/"+filePathPng)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating snapshot: " + err.Error()})
		return
	}

	byteFile, err := os.ReadFile("temp/" + filePathPng)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading file: " + "temp/" + filePathPng + " " + err.Error()})
		return
	}

	publicURL, err := models.SaveAsS3(byteFile, "png", user.AwsAccessKeyId, user.AwsSecretAccessKey, user.AwsRegion, user.AwsBucket, filePathPng)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving file in S3: " + err.Error()})
		return
	}

	model := models.File{
		PublicURL: publicURL,
	}

	context.JSON(http.StatusOK, model)
}

func makeChartSnapshot(content []byte, image string) error {
	path, file := filepath.Split(image)
	suffix := filepath.Ext(file)[1:]
	fileName := file[0 : len(file)-len(suffix)-1]

	config := &render.SnapshotConfig{
		RenderContent: content,
		Path:          path,
		FileName:      fileName,
		Suffix:        suffix,
		Quality:       1,
		KeepHtml:      false,
	}
	return render.MakeSnapshot(config)
}
