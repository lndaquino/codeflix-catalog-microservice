package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"video-catalog/models"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateCastMember(t *testing.T) {
	gin.SetMode(gin.TestMode)
	if err := refreshCastMemberTable(); err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON  string
		statusCode int
		Name       string
		Type       int
	}{
		{
			// basic creation
			inputJSON:  `{"name":"castMember name", "type":1}`,
			statusCode: http.StatusCreated,
			Name:       "castMember name",
			Type:       1,
		},
		{
			// duplicated castMember name forbidden
			inputJSON:  `{"name":"castMember name", "type":1}`,
			statusCode: http.StatusInternalServerError,
		},
		{
			// wrong name data type
			inputJSON:  `{"name":1, "type":1}`,
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			// wrong type data type
			inputJSON:  `{"name":"valid name", "type": "string"}`,
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			// invalid type data type
			inputJSON:  `{"name":"valid name", "type": 0}`,
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			// no data
			inputJSON:  `{}`,
			statusCode: http.StatusUnprocessableEntity,
		},
	}

	for _, v := range samples {
		r := gin.Default()
		r.POST("/cast_member", server.CreateCastMember)
		req, err := http.NewRequest(http.MethodPost, "/cast_member", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, v.statusCode, rr.Code)

		// parsing response body to test json response
		if v.statusCode == http.StatusCreated {
			responseCastMember := models.CastMember{}
			err = json.Unmarshal([]byte(rr.Body.String()), &responseCastMember)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}

			assert.Equal(t, v.Name, responseCastMember.Name)
			assert.Equal(t, v.Type, responseCastMember.Type)
		}
	}
}

func TestGetCastMembers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshCastMemberTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedCastMembers()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/cast_members", server.GetCastMembers)

	req, err := http.NewRequest(http.MethodGet, "/cast_members", nil)
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	castMembersMap := []models.CastMember{}
	err = json.Unmarshal([]byte(rr.Body.String()), &castMembersMap)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, 2, len(castMembersMap))
}

func TestGetCastMemberByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshCastMemberTable()
	if err != nil {
		log.Fatal(err)
	}
	castMember, err := seedOneCastMember()
	if err != nil {
		log.Fatal(err)
	}

	castMemberSample := []struct {
		id         string
		statusCode int
		Name       string
		Type       int
	}{
		{
			// valid castMember
			id:         castMember.ID,
			statusCode: http.StatusOK,
			Name:       castMember.Name,
			Type:       castMember.Type,
		},
		{
			// invalid id parameter
			id:         "invalid id parameter",
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			// invalid castMember
			id:         uuid.NewV4().String(),
			statusCode: http.StatusNotFound,
		},
	}

	for _, v := range castMemberSample {
		req, _ := http.NewRequest(http.MethodGet, "/cast_member/"+v.id, nil)
		rr := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/cast_member/:id", server.GetCastMember)
		r.ServeHTTP(rr, req)

		assert.Equal(t, v.statusCode, rr.Code)

		if v.statusCode == http.StatusOK {
			responseCastMember := models.CastMember{}
			err = json.Unmarshal([]byte(rr.Body.String()), &responseCastMember)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}

			assert.Equal(t, v.Name, responseCastMember.Name)
			assert.Equal(t, v.Type, responseCastMember.Type)
		}

	}
}

func TestUpdateCastMember(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshCastMemberTable()
	if err != nil {
		log.Fatal(err)
	}
	castMembers, err := seedCastMembers()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		id         string
		updateJSON string
		statusCode int
		Name       string
		Type       int
	}{
		{
			// valid castMember full update
			id:         castMembers[0].ID,
			statusCode: http.StatusOK,
			updateJSON: `{"name":"updated name", "type":2}`,
			Name:       "updated name",
			Type:       2,
		},
		{
			// valid castMember partial update
			id:         castMembers[0].ID,
			statusCode: http.StatusOK,
			updateJSON: `{"name":"updated name again"}`,
			Name:       "updated name again",
			Type:       2,
		},
		{
			// valid castMember partial update
			id:         castMembers[0].ID,
			statusCode: http.StatusOK,
			updateJSON: `{"type": 1}`,
			Name:       "updated name again",
			Type:       1,
		},
		{
			// invalid id
			id:         "abc",
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"name":"updated name", "type": 2}`,
		},
		{
			// castMember not found
			id:         uuid.NewV4().String(),
			statusCode: http.StatusInternalServerError,
			updateJSON: `{"name":"updated name", "type": 2}`,
		},
		{
			// invalid name data type
			id:         castMembers[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"name":1, "type": 2}`,
		},
		{
			// invalid type data type
			id:         castMembers[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"name":"updated name", "type": "1"}`,
		},
		{
			// wrong value range on type data type
			id:         castMembers[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"name":"updated name", "type": 3}`,
		},
		{
			// no data
			id:         castMembers[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{}`,
		},
		{
			// updating to an existing name
			id:         castMembers[0].ID,
			statusCode: http.StatusInternalServerError,
			updateJSON: `{"name":"` + castMembers[1].Name + `", "type": 2}`,
		},
	}

	for _, v := range samples {
		r := gin.Default()
		r.PUT("/cast_member/:id", server.UpdateCastMember)

		req, err := http.NewRequest(http.MethodPut, "/cast_member/"+v.id, bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("Error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, v.statusCode, rr.Code)

		if v.statusCode == http.StatusOK {
			responseCastMember := models.CastMember{}
			err = json.Unmarshal([]byte(rr.Body.String()), &responseCastMember)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}

			assert.Equal(t, v.Name, responseCastMember.Name)
			assert.Equal(t, v.Type, responseCastMember.Type)
		}
	}
}

func TestDeleteCastMember(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshCastMemberTable()
	if err != nil {
		log.Fatal(err)
	}
	castMember, err := seedOneCastMember()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		id         string
		statusCode int
	}{
		{ // valid category
			id:         castMember.ID,
			statusCode: http.StatusNoContent,
		},
		{ // soft deleted category
			id:         castMember.ID,
			statusCode: http.StatusNotFound,
		},
		{ // invalid category
			id:         uuid.NewV4().String(),
			statusCode: http.StatusNotFound,
		},
		{ // invalid id parameter
			id:         "abc",
			statusCode: http.StatusUnprocessableEntity,
		},
		{ // no id parameter
			id:         "",
			statusCode: http.StatusNotFound,
		},
	}

	for _, v := range samples {
		r := gin.Default()
		r.DELETE("/cast_member/:id", server.DeleteCastMember)
		req, _ := http.NewRequest(http.MethodDelete, "/cast_member/"+v.id, nil)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, v.statusCode, rr.Code)
	}
}

func refreshCastMemberTable() error {
	err := server.DB.DropTableIfExists(&models.CastMember{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.CastMember{}).Error
	if err != nil {
		return err
	}
	log.Printf("Sucessfully refreshed CastMember table")
	return nil
}

func seedCastMembers() ([]models.CastMember, error) {
	var err error
	castMembers := []models.CastMember{
		{
			ID:   uuid.NewV4().String(),
			Name: "cast 1",
			Type: models.TypeActor,
		},
		{
			ID:   uuid.NewV4().String(),
			Name: "cast 2",
			Type: models.TypeDirector,
		},
	}

	for i := range castMembers {
		err = server.DB.Model(&models.CastMember{}).Create(&castMembers[i]).Error
		if err != nil {
			return []models.CastMember{}, err
		}
	}
	return castMembers, nil
}

func seedOneCastMember() (models.CastMember, error) {
	castMember := models.CastMember{
		ID:   uuid.NewV4().String(),
		Name: "castMember name",
		Type: models.TypeDirector,
	}

	err := server.DB.Model(&models.CastMember{}).Create(&castMember).Error
	if err != nil {
		return models.CastMember{}, err
	}
	return castMember, nil
}
