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

func TestCreateGenre(t *testing.T) {
	gin.SetMode(gin.TestMode)
	if err := refreshGenreTable(); err != nil {
		log.Fatal(err)
	}

	x := true
	samples := []struct {
		inputJSON  string
		statusCode int
		name       string
		isActive   *bool
	}{
		{
			// basic creation
			inputJSON:  `{"name":"genre name", "is_active": true}`,
			statusCode: http.StatusCreated,
			name:       "genre name",
			isActive:   &x,
		},
		{
			// duplicated genre name forbidden
			inputJSON:  `{"name":"genre name"}`,
			statusCode: http.StatusInternalServerError,
		},
		{
			// only required fields
			inputJSON:  `{"name":"new genre"}`,
			statusCode: http.StatusCreated,
			name:       "new genre",
			isActive:   &x,
		},
		{
			// wrong name data type
			inputJSON:  `{"name":1}`,
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			// wrong is_active data type
			inputJSON:  `{"name":"valid name", "is_active": "invalid"}`,
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
		r.POST("/genre", server.CreateGenre)
		req, err := http.NewRequest(http.MethodPost, "/genre", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, v.statusCode, rr.Code)

		// parsing response body to test json response
		if v.statusCode == http.StatusCreated {
			responseGenre := models.Genre{}
			err = json.Unmarshal([]byte(rr.Body.String()), &responseGenre)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}

			assert.Equal(t, v.name, responseGenre.Name)
			assert.Equal(t, v.isActive, responseGenre.IsActive)
		}
	}
}

func TestGetGenres(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshGenreTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedGenres()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/genres", server.GetGenres)

	req, err := http.NewRequest(http.MethodGet, "/genres", nil)
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	categoriesMap := []models.Category{}
	err = json.Unmarshal([]byte(rr.Body.String()), &categoriesMap)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, 2, len(categoriesMap))
}

func TestGetGenreyByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshGenreTable()
	if err != nil {
		log.Fatal(err)
	}
	genre, err := seedOneGenre()
	if err != nil {
		log.Fatal(err)
	}

	genreSample := []struct {
		id         string
		statusCode int
		name       string
		isActive   *bool
	}{
		{
			// valid genre
			id:         genre.ID,
			statusCode: http.StatusOK,
			name:       genre.Name,
			isActive:   genre.IsActive,
		},
		{
			// invalid id parameter
			id:         "invalid id parameter",
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			// invalid genre
			id:         uuid.NewV4().String(),
			statusCode: http.StatusNotFound,
		},
	}

	for _, v := range genreSample {
		req, _ := http.NewRequest(http.MethodGet, "/genre/"+v.id, nil)
		rr := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/genre/:id", server.GetGenre)
		r.ServeHTTP(rr, req)

		assert.Equal(t, v.statusCode, rr.Code)

		if v.statusCode == http.StatusOK {
			responseGenre := models.Genre{}
			err = json.Unmarshal([]byte(rr.Body.String()), &responseGenre)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}

			assert.Equal(t, v.name, responseGenre.Name)
			assert.Equal(t, v.isActive, responseGenre.IsActive)
		}

	}
}

func TestUpdateGenre(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshGenreTable()
	if err != nil {
		log.Fatal(err)
	}
	genres, err := seedGenres()
	if err != nil {
		log.Fatal(err)
	}

	isFalse := false
	samples := []struct {
		id         string
		updateJSON string
		statusCode int
		name       string
		isActive   *bool
	}{
		{
			// valid genre full update
			id:         genres[0].ID,
			statusCode: http.StatusOK,
			updateJSON: `{"name":"updated name", "is_active": false}`,
			name:       "updated name",
			isActive:   &isFalse,
		},
		{
			// valid genre partial update
			id:         genres[0].ID,
			statusCode: http.StatusOK,
			updateJSON: `{"name":"updated name again"}`,
			name:       "updated name again",
			isActive:   &isFalse,
		},
		{
			// invalid id
			id:         "abc",
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"name":"updated name", "is_active": false}`,
		},
		{
			// genre not found
			id:         uuid.NewV4().String(),
			statusCode: http.StatusInternalServerError,
			updateJSON: `{"name":"updated name", "is_active": false}`,
		},
		{
			// missing required field
			id:         genres[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"is_active": false}`,
		},
		{
			// invalid name data type
			id:         genres[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"name":1, "is_active": false}`,
		},
		{
			// invalid is_active data type
			id:         genres[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"name":"updated name", "is_active": "false"}`,
		},
		{
			// no data
			id:         genres[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{}`,
		},
		{
			// updating to an existing name
			id:         genres[0].ID,
			statusCode: http.StatusInternalServerError,
			updateJSON: `{"name":"` + genres[1].Name + `", "description":"updated description", "is_active": false}`,
		},
	}

	for _, v := range samples {
		r := gin.Default()
		r.PUT("/genre/:id", server.UpdateGenre)

		req, err := http.NewRequest(http.MethodPut, "/genre/"+v.id, bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("Error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, v.statusCode, rr.Code)

		if v.statusCode == http.StatusOK {
			responseGenre := models.Genre{}
			err = json.Unmarshal([]byte(rr.Body.String()), &responseGenre)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}

			assert.Equal(t, v.name, responseGenre.Name)
			assert.Equal(t, v.isActive, responseGenre.IsActive)
		}
	}
}

func TestDeleteGenre(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshGenreTable()
	if err != nil {
		log.Fatal(err)
	}
	genre, err := seedOneGenre()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		id         string
		statusCode int
	}{
		{ // valid category
			id:         genre.ID,
			statusCode: http.StatusNoContent,
		},
		{ // soft deleted genre
			id:         genre.ID,
			statusCode: http.StatusNotFound,
		},
		{ // invalid genre
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
		r.DELETE("/genre/:id", server.DeleteGenre)
		req, _ := http.NewRequest(http.MethodDelete, "/genre/"+v.id, nil)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, v.statusCode, rr.Code)
	}
}

func refreshGenreTable() error {
	err := server.DB.DropTableIfExists(&models.Genre{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Genre{}).Error
	if err != nil {
		return err
	}
	log.Printf("Sucessfully refreshed table")
	return nil
}

func seedGenres() ([]models.Genre, error) {
	var err error
	x := true
	genres := []models.Genre{
		{
			ID:       uuid.NewV4().String(),
			Name:     "genre 1",
			IsActive: &x,
		},
		{
			ID:       uuid.NewV4().String(),
			Name:     "genre 2",
			IsActive: &x,
		},
	}

	for i := range genres {
		err = server.DB.Model(&models.Genre{}).Create(&genres[i]).Error
		if err != nil {
			return []models.Genre{}, err
		}
	}
	return genres, nil
}

func seedOneGenre() (models.Genre, error) {
	x := true
	genre := models.Genre{
		ID:       uuid.NewV4().String(),
		Name:     "genre name",
		IsActive: &x,
	}

	err := server.DB.Model(&models.Genre{}).Create(&genre).Error
	if err != nil {
		return models.Genre{}, err
	}
	return genre, nil
}
