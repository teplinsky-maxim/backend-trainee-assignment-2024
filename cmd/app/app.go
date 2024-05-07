package main

import (
	"avito-backend-2024-trainee/config"
	"avito-backend-2024-trainee/internal/app"
)

func main() {
	app.Run(config.DefaultConfigPath)
}
