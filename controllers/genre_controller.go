package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"video-catalog/models"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// CreateGenre controller handles genre creation request
func (server *Server) CreateGenre(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	genre := models.Genre{}
	if err = json.Unmarshal(body, &genre); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}
	genre.ID = uuid.NewV4().String()

	genre.Prepare()
	if err := genre.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	genreCreated, err := genre.Create(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing request",
		})
		return
	}

	c.JSON(http.StatusCreated, genreCreated)
}

// GetGenres handles genres list request
func (server *Server) GetGenres(c *gin.Context) {
	genre := models.Genre{}

	genres, err := genre.FindAll(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, genres)
}

// GetGenre handles genre search request
func (server *Server) GetGenre(c *gin.Context) {
	genreID := c.Param("id")
	if _, err := uuid.FromString(genreID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}
	genre := models.Genre{ID: genreID}

	err := genre.FindByID(server.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, genre)
}

// UpdateGenre handles genre update requests
func (server *Server) UpdateGenre(c *gin.Context) {
	genreID := c.Param("id")

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	newGenre := models.Genre{}
	if err = json.Unmarshal(body, &newGenre); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	newGenre.ID = genreID
	newGenre.Prepare()
	if err := newGenre.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	updatedGenre, err := newGenre.Update(server.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, updatedGenre)
}

// DeleteGenre handles genre delete requests
func (server *Server) DeleteGenre(c *gin.Context) {
	genreID := c.Param("id")
	if _, err := uuid.FromString(genreID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}
	genre := models.Genre{ID: genreID}

	if err := genre.Delete(server.DB); err != nil {
		if err.Error() == "Genre not found" {
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
