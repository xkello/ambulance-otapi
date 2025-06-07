package hospital_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implHospitalRolesAPI struct {
}

func NewHospitalRolesApi() HospitalRolesAPI {
	return &implHospitalRolesAPI{}
}

func (o *implHospitalRolesAPI) GetRoles(c *gin.Context) {
	updateHospitalFunc(c, func(
		c *gin.Context,
		hospital *Hospital,
	) (updatedHospital *Hospital, responseContent interface{}, status int) {
		result := hospital.PredefinedRoles
		if result == nil {
			result = []Role{}
		}
		return nil, result, http.StatusOK
	})
}
