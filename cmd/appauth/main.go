package main

import (
	"github.com/Van-programan/Forum_GO/internal/app"
)

// @title Auth Service API
// @version 1.0
// @description API for auth service
// @host localhost:3100
// @BasePath /
// @schemes http
func main() {
	app.RunAuthServer()
}
