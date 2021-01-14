package tests

import (
	"testing"
	"time"
	"video-catalog/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := models.Video{}
	err := video.Validate("create")
	require.Error(t, err)
}

func TestValidateVideoFullFilled(t *testing.T) {
	isTrue := true
	video := models.Video{
		ID:           uuid.New().String(),
		Title:        "Film title",
		Description:  "Film description with many details, sinopse, cast and marketing descriptions",
		YearLaunched: 2020,
		Opened:       &isTrue,
		Rating:       "L",
		Duration:     90,
	}
	video.Prepare()
	err := video.Validate("create")
	require.Nil(t, err)
}

func TestValidateVideoIDIsNotUUID(t *testing.T) {
	isTrue := true
	video := models.Video{
		ID:           "invalid id",
		Title:        "Film title",
		Description:  "Film description with many details, sinopse, cast and marketing descriptions",
		YearLaunched: 2020,
		Opened:       &isTrue,
		Rating:       "L",
		Duration:     90,
	}
	video.Prepare()
	err := video.Validate("create")
	require.Error(t, err)
}

func TestValidateVideoTitle(t *testing.T) {
	isTrue := true
	video := models.Video{
		ID:           uuid.New().String(),
		Description:  "Film description with many details, sinopse, cast and marketing descriptions",
		YearLaunched: 2020,
		Opened:       &isTrue,
		Rating:       "L",
		Duration:     90,
	}
	video.Prepare()
	err := video.Validate("create")
	require.Error(t, err)

	video.Title = "ab"
	video.Prepare()
	err = video.Validate("create")
	require.Error(t, err)

	video.Title = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	video.Prepare()
	err = video.Validate("create")
	require.Error(t, err)
}

func TestValidateVideoDescription(t *testing.T) {
	isTrue := true
	video := models.Video{
		ID:           uuid.New().String(),
		Title:        "Film title",
		YearLaunched: 2020,
		Opened:       &isTrue,
		Rating:       "L",
		Duration:     90,
	}
	video.Prepare()
	err := video.Validate("create")
	require.Error(t, err)

	video.Description = "-Film details-"
	video.Prepare()
	err = video.Validate("create")
	require.Error(t, err)

	video.Description = "a b c d e f a b    c"
	video.Prepare()
	err = video.Validate("create")
	require.Error(t, err)
}

func TestValidateVideoYearLaunched(t *testing.T) {
	isTrue := true
	video := models.Video{
		ID:          uuid.New().String(),
		Title:       "Film title",
		Description: "Film description with many details, sinopse, cast and marketing descriptions",
		Opened:      &isTrue,
		Rating:      "L",
		Duration:    90,
	}
	video.Prepare()
	err := video.Validate("create")
	require.Error(t, err)

	video.YearLaunched = 1894
	video.Prepare()
	err = video.Validate("create")
	require.Error(t, err)

	video.YearLaunched = time.Now().Year() + 1
	video.Prepare()
	err = video.Validate("create")
	require.Error(t, err)
}

func TestValidateVideoRating(t *testing.T) {
	invalidSamples := []string{"X", "Livre", "99", "0", ""}
	validSamples := []string{"L", "10", "12", "14", "16", "18"}

	isTrue := true
	video := models.Video{
		ID:           uuid.New().String(),
		Title:        "Film title",
		Description:  "Film description with many details, sinopse, cast and marketing descriptions",
		YearLaunched: 2020,
		Opened:       &isTrue,
		Duration:     90,
	}
	video.Prepare()
	err := video.Validate("create")
	require.Error(t, err)

	for _, sample := range invalidSamples {
		video.Rating = sample
		video.Prepare()
		err = video.Validate("create")
		require.Error(t, err)
	}

	for _, sample := range validSamples {
		video.Rating = sample
		video.Prepare()
		err = video.Validate("create")
		require.Nil(t, err)
	}
}

func TestValidateVideoDuration(t *testing.T) {
	isTrue := true
	video := models.Video{
		ID:           uuid.New().String(),
		Title:        "Film title",
		Description:  "Film description with many details, sinopse, cast and marketing descriptions",
		YearLaunched: 2020,
		Opened:       &isTrue,
		Rating:       "L",
	}
	video.Prepare()
	err := video.Validate("create")
	require.Error(t, err)

	video.Duration = 0
	video.Prepare()
	err = video.Validate("create")
	require.Error(t, err)

	video.Duration = 90
	video.Prepare()
	err = video.Validate("create")
	require.Nil(t, err)
}

func TestValidateUpdate(t *testing.T) {
	isTrue := true
	video := models.Video{
		ID:           uuid.New().String(),
		Title:        "valid title",
		Description:  "Film description with many details, sinopse, cast and marketing descriptions",
		YearLaunched: 2020,
		Opened:       &isTrue,
		Rating:       "L",
		Duration:     90,
	}
	video.Prepare()
	err := video.Validate("update")
	require.Nil(t, err)

	video = models.Video{}
	video.Prepare()
	err = video.Validate("update")
	require.Error(t, err)

	video.ID = "invalid id"
	video.Prepare()
	err = video.Validate("update")
	require.Error(t, err)

	video.ID = uuid.New().String()
	video.Title = "ab"
	video.Prepare()
	err = video.Validate("update")
	require.Error(t, err)
	video.Title = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	video.Prepare()
	err = video.Validate("update")
	require.Error(t, err)
	video.Title = "valid title"
	video.Prepare()
	err = video.Validate("update")
	require.Nil(t, err)

	video.Description = "-Film details-"
	video.Prepare()
	err = video.Validate("update")
	require.Error(t, err)
	video.Description = "a b c d e f a b    c "
	video.Prepare()
	err = video.Validate("update")
	require.Error(t, err)
	video.Description = "Film description with many details, sinopse, cast and marketing descriptions"
	video.Prepare()
	err = video.Validate("update")
	require.Nil(t, err)

	video.YearLaunched = 1894
	video.Prepare()
	err = video.Validate("update")
	require.Error(t, err)
	video.YearLaunched = time.Now().Year() + 1
	video.Prepare()
	err = video.Validate("create")
	require.Error(t, err)
	video.YearLaunched = 2020
	video.Prepare()
	err = video.Validate("update")
	require.Nil(t, err)

	invalidSamples := []string{"X", "Livre", "99", "0"}
	validSamples := []string{"L", "10", "12", "14", "16", "18"}
	for _, sample := range invalidSamples {
		video.Rating = sample
		video.Prepare()
		err = video.Validate("update")
		require.Error(t, err)
	}
	for _, sample := range validSamples {
		video.Rating = sample
		video.Prepare()
		err = video.Validate("update")
		require.Nil(t, err)
	}

	video.Duration = 0
	video.Prepare()
	err = video.Validate("update")
	require.Nil(t, err)
	video.Duration = 1
	video.Prepare()
	err = video.Validate("update")
	require.Nil(t, err)
}
