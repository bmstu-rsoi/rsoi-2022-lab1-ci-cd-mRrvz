package main

import (
	"log"

	"github.com/bmstu-rsoi/rsoi-2022-lab1-ci-cd-mRrvz/src/server"
)

func main() {
	log.Fatalln(server.SetupServer().Run())
}
