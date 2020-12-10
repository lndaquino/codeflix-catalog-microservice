package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"video-catalog/models"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// CreateCategory controller handles category creation request
func (server *Server) CreateCategory(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	category := models.Category{}
	if err = json.Unmarshal(body, &category); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}
	category.ID = uuid.NewV4().String()

	category.Prepare()
	if err := category.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	categoryCreated, err := category.Create(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing request",
		})
		return
	}

	c.JSON(http.StatusCreated, categoryCreated)
}

// GetCategories handles categories list request
func (server *Server) GetCategories(c *gin.Context) {
	category := models.Category{}

	categories, err := category.FindAll(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategory handles category search request
func (server *Server) GetCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if _, err := uuid.FromString(categoryID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}
	category := models.Category{ID: categoryID}

	err := category.FindByID(server.DB)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory handles category update requests
func (server *Server) UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	newCategory := models.Category{}
	if err = json.Unmarshal(body, &newCategory); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	newCategory.ID = categoryID
	newCategory.Prepare()
	if err := newCategory.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	updatedCategory, err := newCategory.Update(server.DB)
	if err != nil {
		if err.Error() == "Internal server error" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

// DeleteCategory handles category delete requests
func (server *Server) DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if _, err := uuid.FromString(categoryID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}
	category := models.Category{ID: categoryID}

	if err := category.Delete(server.DB); err != nil {
		if err.Error() == "Category not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
