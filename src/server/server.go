package server

import (
	"context"
	"log"

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
		Host:     "ec2-18-214-35-70.compute-1.amazonaws.com",                         /* os.Getenv("POSTGRES_HOST"), */
		Port:     "5432",                                                             /* os.Getenv("POSTGRES_PORT"), */
		User:     "kygszpllzydzcl",                                                   /* os.Getenv("POSTGRES_USER"), */
		Password: "e2fae91b45875b8e34b2fb98abffa2d7a8fe5f8e283fd7d7a2ce3975ec3a48a3", /* os.Getenv("POSTGRES_PASS"), */
		DB:       "d2vcb8m0neje4o",                                                   /* os.Getenv("POSTGRES_DB"), */
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
