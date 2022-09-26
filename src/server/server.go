package server

import (
	"context"
	"log"
	"os"

	"github.com/bmstu-rsoi/rsoi-2022-lab1-ci-cd-mRrvz/src/db"
	"github.com/bmstu-rsoi/rsoi-2022-lab1-ci-cd-mRrvz/src/repo"
	"github.com/bmstu-rsoi/rsoi-2022-lab1-ci-cd-mRrvz/src/server/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	postgresRepo repo.PostgresRepository
)

func init() {
	ctx := context.Background()

	postgresClient, err := db.NewPostgresClient(ctx, &db.Options{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASS"),
		DB:       os.Getenv("POSTGRES_DB"),
	})

	if err != nil {
		log.Fatalf("Failed to create postgres client, error is: %s", err)
	}

	postgresRepo = repo.NewPostgresRepo(*postgresClient)

	postgresRepo.CreateTable()
}

func SetupServer() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"token", "Origin", "X-Requested-With", "Content-Type", "Accept"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE"}

	r.Use(cors.New(config))
	r.Use(func(c *gin.Context) {
		c.Set("postgres_repo", postgresRepo)
		c.Next()
	})

	v1 := r.Group("/api/v1")
	{
		persons := v1.Group("/persons")
		{
			persons.GET("/:id", controllers.GetPerson)
			persons.GET("", controllers.GetAllPersons)
			persons.POST("", controllers.CreatePerson)
			persons.PATCH("/:id", controllers.UpdatePerson)
			persons.DELETE("/:id", controllers.DeletePerson)
		}
	}

	return r
}
