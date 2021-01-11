package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"video-catalog/models"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// CreateCastMember controller handles castMember creation request
func (server *Server) CreateCastMember(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	castMember := models.CastMember{}
	if err = json.Unmarshal(body, &castMember); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}
	castMember.ID = uuid.NewV4().String()

	castMember.Prepare()
	if err := castMember.Validate("create"); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	castMemberCreated, err := castMember.Create(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing request",
		})
		return
	}

	c.JSON(http.StatusCreated, castMemberCreated)
}

// GetCastMembers handles cast members list request
func (server *Server) GetCastMembers(c *gin.Context) {
	castMember := models.CastMember{}

	castMembers, err := castMember.FindAll(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, castMembers)
}

// GetCastMember handles cast member search request
func (server *Server) GetCastMember(c *gin.Context) {
	castMemberID := c.Param("id")
	if _, err := uuid.FromString(castMemberID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}
	castMember := models.CastMember{ID: castMemberID}

	err := castMember.FindByID(server.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, castMember)
}

// UpdateCastMember handles cast member update requests
func (server *Server) UpdateCastMember(c *gin.Context) {
	castMemberID := c.Param("id")

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	newCastMember := models.CastMember{}
	if err = json.Unmarshal(body, &newCastMember); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	newCastMember.ID = castMemberID
	newCastMember.Prepare()
	if err := newCastMember.Validate("update"); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}

	updatedCastMember, err := newCastMember.Update(server.DB)
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

	c.JSON(http.StatusOK, updatedCastMember)
}

// DeleteCastMember handles cast member delete requests
func (server *Server) DeleteCastMember(c *gin.Context) {
	castMemberID := c.Param("id")
	if _, err := uuid.FromString(castMemberID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
		return
	}
	castMember := models.CastMember{ID: castMemberID}

	if err := castMember.Delete(server.DB); err != nil {
		if err.Error() == "CastMember not found" {
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
