package hospital_wl

import (
	"net/http"

	"github.com/Serbel97/ambulance-api/internal/db_service"
	"github.com/gin-gonic/gin"
)

type hospitalUpdater = func(
	ctx *gin.Context,
	hospital *Hospital,
) (updatedHospital *Hospital, responseContent interface{}, status int)

func updateHospitalFunc(ctx *gin.Context, updater hospitalUpdater) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
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
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	hospitalId := ctx.Param("hospitalId")

	hospital, err := db.FindDocument(ctx, hospitalId)

	switch err {
	case nil:
		// continue
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Hospital not found",
				"error":   err.Error(),
			},
		)
		return
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to load hospital from database",
				"error":   err.Error(),
			})
		return
	}

	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to cast hospital from database",
				"error":   "Failed to cast hospital from database",
			})
		return
	}

	updatedHospital, responseObject, status := updater(ctx, hospital)

	if updatedHospital != nil {
		err = db.UpdateDocument(ctx, hospitalId, updatedHospital)
	} else {
		err = nil // redundant but for clarity
	}

	switch err {
	case nil:
		if responseObject != nil {
			ctx.JSON(status, responseObject)
		} else {
			ctx.AbortWithStatus(status)
		}
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Hospital was deleted while processing the request",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to update hospital in database",
				"error":   err.Error(),
			})
	}
}
