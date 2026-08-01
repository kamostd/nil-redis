package main

import "nil-redis/internal/core/app"

// @title           Leaderboard API
// @version         1.0
// @description     This is a project for getting acquainted with Redis.

// @host      localhost:8080
// @BasePath  /api/

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	app.App().Run()
}
