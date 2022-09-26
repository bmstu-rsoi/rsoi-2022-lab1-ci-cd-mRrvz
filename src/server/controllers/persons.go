package controllers

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"strconv"

	"github.com/bmstu-rsoi/rsoi-2022-lab1-ci-cd-mRrvz/src/db"
	"github.com/bmstu-rsoi/rsoi-2022-lab1-ci-cd-mRrvz/src/repo"
	"github.com/bmstu-rsoi/rsoi-2022-lab1-ci-cd-mRrvz/src/server/models"
	"github.com/gin-gonic/gin"
)

func GetPerson(c *gin.Context) {
	personId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Failed to get id from gin.Context parameters: %s", err)
		c.AbortWithStatusJSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	repo, ok := c.MustGet("postgres_repo").(repo.PostgresRepository)
	if !ok {
		log.Errorf("Failed to get repo from gin.Context")
		c.AbortWithStatusJSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	person, err := getPerson(repo, personId)
	if err != nil {
		log.Errorf("Failed to get person with id = %d, err = %s", personId, err)
		c.AbortWithStatusJSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	if person == nil {
		c.JSON(404, models.ErrorResponse{Message: "not found"})
		return
	}

	c.JSON(200, models.PersonResponse{
		Id:      person.PersonID,
		Name:    person.Name,
		Age:     person.Age,
		Address: person.Address,
		Work:    person.Work,
	})
}

func GetAllPersons(c *gin.Context) {
	repo, ok := c.MustGet("postgres_repo").(repo.PostgresRepository)
	if !ok {
		log.Errorf("Failed to get repo from gin.Context")
		c.AbortWithStatusJSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	persons, err := getAllPersons(repo)
	if err != nil {
		log.Errorf("Failed to get all persons, err = %s", err)
		c.AbortWithStatusJSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	response := make([]models.PersonResponse, 0)
	for _, el := range persons {
		r := models.PersonResponse{
			Id:      el.PersonID,
			Name:    el.Name,
			Age:     el.Age,
			Address: el.Address,
			Work:    el.Work,
		}

		response = append(response, r)
	}

	c.JSON(200, response)
}

func CreatePerson(c *gin.Context) {
	repo, ok := c.MustGet("postgres_repo").(repo.PostgresRepository)
	if !ok {
		log.Errorf("Failed to get repo from gin.Context")
		c.AbortWithStatusJSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	var req models.PersonRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("failed to parse request body, error is: %s", err)
		c.AbortWithStatusJSON(400, models.ValidationErrorResponse{
			Message: "Failed to parse request body",
			Errors:  models.Errors{AdditionalProperties: fmt.Sprintf("%s", err)},
		})
		return
	}

	person := db.Person{
		Name:    req.Name,
		Age:     req.Age,
		Address: req.Address,
		Work:    req.Work,
	}

	personId, err := createPerson(repo, person)
	if err != nil {
		log.Errorf("Failed to create person, err = %s", err)
		c.JSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/persons/%d", personId))
	c.JSON(201, personId)
}

func UpdatePerson(c *gin.Context) {
	personId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Failed to get id from gin.Context parameters: %s", err)
		c.AbortWithStatusJSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	repo, ok := c.MustGet("postgres_repo").(repo.PostgresRepository)
	if !ok {
		log.Errorf("Failed to get repo from gin.Context")
		c.AbortWithStatusJSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	var req models.PersonRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("failed to parse request body, error is: %s", err)
		c.AbortWithStatusJSON(400, models.ValidationErrorResponse{
			Message: "Failed to parse request body",
			Errors:  models.Errors{AdditionalProperties: fmt.Sprintf("%s", err)},
		})
		return
	}

	person := db.Person{
		PersonID: int32(personId),
		Name:     req.Name,
		Age:      req.Age,
		Address:  req.Address,
		Work:     req.Work,
	}

	if err := updatePerson(repo, person); err != nil {
		log.Errorf("Failed to update person with id = %d, err = %s", personId, err)
		c.JSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	// todo
	if 0 == 1 {
		c.JSON(404, models.ErrorResponse{Message: "not found"})
		return
	}

	updatedPerson, err := getPerson(repo, personId)
	if err != nil {
		log.Errorf("Failed to get person with id = %d, err = %s", personId, err)
		c.AbortWithStatusJSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	c.JSON(200, models.PersonResponse{
		Id:      updatedPerson.PersonID,
		Name:    updatedPerson.Name,
		Age:     updatedPerson.Age,
		Address: updatedPerson.Address,
		Work:    updatedPerson.Work,
	})
}

func DeletePerson(c *gin.Context) {
	personId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Failed to get id from gin.Context parameters: %s", err)
		c.JSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	repo, ok := c.MustGet("postgres_repo").(repo.PostgresRepository)
	if !ok {
		log.Errorf("Failed to get repo from gin.Context")
		c.JSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	if err := deletePerson(repo, personId); err != nil {
		log.Errorf("Failed to remove person with id = %d, err = %s", personId, err)
		c.JSON(500, models.ErrorResponse{Message: "internal error"})
		return
	}

	c.JSON(204, nil)
}

func getPerson(repo repo.PostgresRepository, personID int) (*db.Person, error) {
	person, err := repo.Get(personID)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func getAllPersons(repo repo.PostgresRepository) ([]db.Person, error) {
	person, err := repo.GetAll()
	if err != nil {
		return nil, err
	}

	return person, nil
}

func deletePerson(repo repo.PostgresRepository, personID int) error {
	if err := repo.Delete(personID); err != nil {
		return err
	}

	return nil
}

func createPerson(repo repo.PostgresRepository, person db.Person) (int32, error) {
	return repo.Create(person)
}

func updatePerson(repo repo.PostgresRepository, person db.Person) error {
	if err := repo.Update(person, int(person.PersonID)); err != nil {
		return err
	}

	return nil
}
