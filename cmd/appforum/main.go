package main

import "github.com/Van-programan/Forum_GO/internal/app"

// @title Forum Service API
// @version 1.0
// @description API for forum service
// @host localhost:3101
// @BasePath /
// @schemes http
func main() {
	app.RunForumServer()
}
