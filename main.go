package main

import (
	"os"

	"github.com/OdairPianta/gryphon_api/docs"
	"github.com/OdairPianta/gryphon_api/models"
	"github.com/OdairPianta/gryphon_api/routes"
)

// swagger embed files
// @contact.name   API Support
// @contact.url    https://spotec.app/contato/
// @contact.email  contact@spotec.app
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	models.InitApp()

	docs.SwaggerInfo.Title = "Gryphon API"
	docs.SwaggerInfo.Description = "Gryphon API is a REST API that allows send files and manage files."
	docs.SwaggerInfo.Version = "1.1.2"
	docs.SwaggerInfo.Host = os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}

	routes.HandleRequest()
}
