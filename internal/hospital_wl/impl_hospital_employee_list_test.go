package hospital_wl

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Serbel97/ambulance-api/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type HospitalWlSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[Hospital]
}

func TestHospitalWlSuite(t *testing.T) {
	suite.Run(t, new(HospitalWlSuite))
}

type DbServiceMock[DocType interface{}] struct {
	mock.Mock
}

func (this *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (m *DbServiceMock[DocType]) ListDocuments(ctx context.Context) ([]DocType, error) {
	args := m.Called(ctx)
	return args.Get(0).([]DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := this.Called(ctx, id)
	return args.Get(0).(*DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := this.Called(ctx, id)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := this.Called(ctx)
	return args.Error(0)
}

func (suite *HospitalWlSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[Hospital]{}

	var _ db_service.DbService[Hospital] = suite.dbServiceMock

	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Hospital{
				Id: "test-hospital",
				EmployeeList: []EmployeeListEntry{
					{
						Id: "test-entry",
					},
				},
			},
			nil,
		)
}

func (suite *HospitalWlSuite) Test_UpdateWl_DbServiceUpdateCalled() {
	suite.dbServiceMock.On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	json := `{
        "id": "test-entry"
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "hospitalId", Value: "test-hospital"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("POST", "/api/hospital/test-hospital/employeelist/test-entry", strings.NewReader(json))

	sut := &implHospitalEmployeeListAPI{} //TODO

	sut.UpdateEmployeeListEntry(ctx)
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-hospital", mock.Anything)
}
