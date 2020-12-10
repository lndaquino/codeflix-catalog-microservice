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

func TestCreateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	if err := refreshCategoryTable(); err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			// basic creation
			inputJSON:  `{"name":"category name", "description":"category description", "is_active": true}`,
			statusCode: http.StatusCreated,
		},
		{
			// duplicated category name forbidden
			inputJSON:  `{"name":"category name"}`,
			statusCode: http.StatusInternalServerError,
		},
		{
			// only required fields
			inputJSON:  `{"name":"new category"}`,
			statusCode: http.StatusCreated,
		},
		{
			// wrong name data type
			inputJSON:  `{"name":1}`,
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			// wrong description data type
			inputJSON:  `{"name":"valid name", "description": 1}`,
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			// wrong is_active data type
			inputJSON:  `{"name":"valid name", "description": "valid description", "is_active": "invalid"}`,
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
		r.POST("/category", server.CreateCategory)
		req, err := http.NewRequest(http.MethodPost, "/category", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, v.statusCode, rr.Code)

		// parsing response body to test json response
		if v.statusCode == http.StatusCreated {
			responseCategory := models.Category{}
			err = json.Unmarshal([]byte(rr.Body.String()), &responseCategory)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			// log.Println(responseCategory)
		}
	}
}

func TestGetCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshCategoryTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedCategories()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/categories", server.GetCategories)

	req, err := http.NewRequest(http.MethodGet, "/categories", nil)
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

func TestGetCategoryByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshCategoryTable()
	if err != nil {
		log.Fatal(err)
	}
	category, err := seedOneCategory()
	if err != nil {
		log.Fatal(err)
	}

	categorySample := []struct {
		id         string
		statusCode int
	}{
		{
			// valid category
			id:         category.ID,
			statusCode: http.StatusOK,
		},
		{
			// invalid id parameter
			id:         "invalid id parameter",
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			// invalid category
			id:         uuid.NewV4().String(),
			statusCode: http.StatusNotFound,
		},
	}

	for _, v := range categorySample {
		req, _ := http.NewRequest(http.MethodGet, "/category/"+v.id, nil)
		rr := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/category/:id", server.GetCategory)
		r.ServeHTTP(rr, req)

		assert.Equal(t, v.statusCode, rr.Code)

	}
}

func TestUpdateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshCategoryTable()
	if err != nil {
		log.Fatal(err)
	}
	categories, err := seedCategories()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		id         string
		updateJSON string
		statusCode int
	}{
		{
			// valid category full update
			id:         categories[0].ID,
			statusCode: http.StatusOK,
			updateJSON: `{"name":"updated name", "description":"updated description", "is_active": false}`,
		},
		{
			// valid category partial update
			id:         categories[0].ID,
			statusCode: http.StatusOK,
			updateJSON: `{"name":"updated name again"}`,
		},
		{
			// invalid id
			id:         "abc",
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"name":"updated name", "description":"updated description", "is_active": false}`,
		},
		{
			// category not found
			id:         uuid.NewV4().String(),
			statusCode: http.StatusNotFound,
			updateJSON: `{"name":"updated name", "description":"updated description", "is_active": false}`,
		},
		{
			// missing required field
			id:         categories[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"description":"updated description", "is_active": false}`,
		},
		{
			// invalid name data type
			id:         categories[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"name":1, "description":"updated description", "is_active": false}`,
		},
		{
			// invalid description data type
			id:         categories[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"name":"updated name", "description":1, "is_active": false}`,
		},
		{
			// invalid is_active data type
			id:         categories[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{"name":"updated name", "description":"updated description", "is_active": "false"}`,
		},
		{
			// no data
			id:         categories[0].ID,
			statusCode: http.StatusUnprocessableEntity,
			updateJSON: `{}`,
		},
		{
			// updating to an existing name
			id:         categories[0].ID,
			statusCode: http.StatusInternalServerError,
			updateJSON: `{"name":"` + categories[1].Name + `", "description":"updated description", "is_active": false}`,
		},
	}

	for _, v := range samples {
		r := gin.Default()
		r.PUT("/category/:id", server.UpdateCategory)

		req, err := http.NewRequest(http.MethodPut, "/category/"+v.id, bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("Error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, v.statusCode, rr.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshCategoryTable()
	if err != nil {
		log.Fatal(err)
	}
	category, err := seedOneCategory()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		id         string
		statusCode int
	}{
		{ // valid category
			id:         category.ID,
			statusCode: http.StatusNoContent,
		},
		{ // soft deleted category
			id:         category.ID,
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
		r.DELETE("/category/:id", server.DeleteCategory)
		req, _ := http.NewRequest(http.MethodDelete, "/category/"+v.id, nil)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, v.statusCode, rr.Code)
	}
}
