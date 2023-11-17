package main

import (
	"github.com/OdairPianta/gryphon_api/models"
	"github.com/OdairPianta/gryphon_api/routes"
)

// swagger embed files

// @title           Gryphon API
// @version         0.1.1
// @description     Gryphon API is a REST API that allows send files and manage files.

// @contact.name   API Support
// @contact.url    https://spotec.app/contato/
// @contact.email  contact@spotec.app

// @host      localhost:8083
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	models.ConnectDatabase()
	routes.HandleRequest()
}
