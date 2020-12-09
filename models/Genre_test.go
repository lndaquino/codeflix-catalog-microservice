package models_test

import (
	"testing"
	"video-catalog/models"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfGenreIsEmpty(t *testing.T) {
	genre := models.Genre{}
	err := genre.Validate()
	require.Error(t, err)
}

func TestValidateGenreFullFilled(t *testing.T) {
	x := true
	genre := models.Genre{
		ID:       uuid.NewV4().String(),
		Name:     "name",
		IsActive: &x,
	}
	genre.Prepare()
	err := genre.Validate()
	require.Nil(t, err)
}

func TestValidateGenreRequiredFieldsFilled(t *testing.T) {
	genre := models.Genre{
		ID:   uuid.NewV4().String(),
		Name: "name",
	}
	genre.Prepare()
	err := genre.Validate()
	require.Nil(t, err)
}

func TestValidateGenreIDisNotUUID(t *testing.T) {
	x := true
	genre := models.Genre{
		ID:       "id",
		Name:     "name",
		IsActive: &x,
	}
	genre.Prepare()
	err := genre.Validate()
	require.Error(t, err)
}

func TestGenreNameIsEmpty(t *testing.T) {
	x := true
	genre := models.Genre{
		ID:       uuid.NewV4().String(),
		IsActive: &x,
	}
	genre.Prepare()
	err := genre.Validate()
	require.Error(t, err)
}

func TestGenreNameIsLessThan3Characters(t *testing.T) {
	genre := models.Genre{
		ID:   uuid.NewV4().String(),
		Name: "ab",
	}
	genre.Prepare()
	err := genre.Validate()
	require.Error(t, err)
}

func TestGenreNameIsMoreThan255Characters(t *testing.T) {
	genre := models.Genre{
		ID:   uuid.NewV4().String(),
		Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}
	genre.Prepare()
	err := genre.Validate()
	require.Error(t, err)
}
