package tests

import (
	"testing"
	"video-catalog/models"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfCategoryIsEmpty(t *testing.T) {
	category := models.Category{}
	err := category.Validate()
	require.Error(t, err)
}

func TestValidateCategoryFullFilled(t *testing.T) {
	x := true
	category := models.Category{
		ID:          uuid.NewV4().String(),
		Name:        "name",
		Description: "description",
		IsActive:    &x,
	}
	category.Prepare()
	err := category.Validate()
	require.Nil(t, err)
}

func TestValidateCategoryRequiredFieldsFilled(t *testing.T) {
	category := models.Category{
		ID:   uuid.NewV4().String(),
		Name: "name",
	}
	category.Prepare()
	err := category.Validate()
	require.Nil(t, err)
}

func TestValidateCategoryIDisNotUUID(t *testing.T) {
	x := true
	category := models.Category{
		ID:          "id",
		Name:        "name",
		Description: "description",
		IsActive:    &x,
	}
	category.Prepare()
	err := category.Validate()
	require.Error(t, err)
}

func TestCategoryNameIsEmpty(t *testing.T) {
	x := true
	category := models.Category{
		ID:          uuid.NewV4().String(),
		Description: "description",
		IsActive:    &x,
	}
	category.Prepare()
	err := category.Validate()
	require.Error(t, err)
}

func TestCategoryNameIsLessThan3Characters(t *testing.T) {
	category := models.Category{
		ID:   uuid.NewV4().String(),
		Name: "ab",
	}
	category.Prepare()
	err := category.Validate()
	require.Error(t, err)
}

func TestCategoryNameIsMoreThan255Characters(t *testing.T) {
	category := models.Category{
		ID:   uuid.NewV4().String(),
		Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}
	category.Prepare()
	err := category.Validate()
	require.Error(t, err)
}

func TestCategoryDescriptionIsLessThan3Characters(t *testing.T) {
	category := models.Category{
		ID:          uuid.NewV4().String(),
		Name:        "name",
		Description: "ab",
	}
	category.Prepare()
	err := category.Validate()
	require.Error(t, err)
}

func TestCategoryDescriptionIsMoreThan255Characters(t *testing.T) {
	category := models.Category{
		ID:          uuid.NewV4().String(),
		Name:        "name",
		Description: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}
	category.Prepare()
	err := category.Validate()
	require.Error(t, err)
}
