package hospital_wl

import (
	"net/http"

	"github.com/Serbel97/ambulance-api/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implHospitalsAPI struct {
}

func NewHospitalsApi() HospitalsAPI {
	return &implHospitalsAPI{}
}

func (o *implHospitalsAPI) GetHospital(c *gin.Context) {
	v, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_service not found"})
		return
	}
	db, ok := v.(db_service.DbService[Hospital])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_service has wrong type"})
		return
	}

	hospitals, err := db.ListDocuments(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if hospitals == nil {
		hospitals = []Hospital{}
	}
	c.JSON(http.StatusOK, hospitals)
}

func (o *implHospitalsAPI) CreateHospital(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Hospital])
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}

	hospital := Hospital{}
	err := c.BindJSON(&hospital)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	if hospital.Id == "" {
		hospital.Id = uuid.New().String()
	}

	err = db.CreateDocument(c, hospital.Id, &hospital)

	switch err {
	case nil:
		c.JSON(
			http.StatusCreated,
			hospital,
		)
	case db_service.ErrConflict:
		c.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Hospital already exists",
				"error":   err.Error(),
			},
		)
	default:
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create hospital in database",
				"error":   err.Error(),
			},
		)
	}
}

func (o *implHospitalsAPI) DeleteHospital(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Hospital])
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	hospitalId := c.Param("hospitalId")
	err := db.DeleteDocument(c, hospitalId)

	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Hospital not found",
				"error":   err.Error(),
			},
		)
	default:
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete hospital from database",
				"error":   err.Error(),
			},
		)
	}
}
