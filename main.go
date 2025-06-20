package main

import (
	"smart-scene-app-api/cmd"
	_ "smart-scene-app-api/docs"
)

// @title           Smart Scene App API
// @version         1.0
// @description     API Server for Smart Scene Application
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	cmd.Execute()
}
