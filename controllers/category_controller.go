package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"video-catalog/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	category.Prepare()

	if err := category.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	categoryCreated, err := category.CreateCategory(server.DB)
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

	categories, err := category.FindAllCategories(server.DB)
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
	if _, err := uuid.Parse(categoryID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}
	category := models.Category{ID: categoryID}

	err := category.FindCategoryByID(server.DB)
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

	updatedCategory, err := newCategory.UpdateCategory(server.DB)
	if err != nil {
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
	if _, err := uuid.Parse(categoryID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}
	category := models.Category{ID: categoryID}

	if err := category.DeleteCategory(server.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
