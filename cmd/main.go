package main

import (
	"TPForum/internal/app/server"
	"TPForum/internal/pkg/config"
)

func main() {
	server.RunServer(config.Get().Main.Port)
}
