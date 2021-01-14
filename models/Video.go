package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RatingList enum
var RatingList = []string{"L", "10", "12", "14", "16", "18"}

// Video models a video
type Video struct {
	ID           string     `json:"id" valid:"uuid" gorm:"type:uuid;primary_key"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	YearLaunched int        `json:"year_launched"`
	Opened       *bool      `json:"opened" gorm:"default:false"`
	Rating       string     `json:"rating"`
	Duration     int        `json:"duration"`
	CreatedAt    *time.Time `json:"created_at,omitempty" valid:"-" gorm:"autoCreateTime"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" valid:"-" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" valid:"-" gorm:"autoDeleteTime"`
}

// Prepare cleans string fields
func (v *Video) Prepare() {
	v.Title = html.EscapeString(strings.TrimSpace(v.Title))
	v.Description = html.EscapeString(strings.TrimSpace(v.Description))
}

// Validate validates Video struct
func (v *Video) Validate(action string) error {
	if _, err := uuid.Parse(v.ID); err != nil {
		return errors.New("Invalid id")
	}

	lenTitle := len(v.Title)
	lenDescription := len(v.Description)
	descriptionWords := len(strings.Fields(v.Description))

	switch action {
	case "create":
		if err := validateTitle(lenTitle); err != nil {
			return err
		}
		if err := validateDescription(lenDescription, descriptionWords); err != nil {
			return err
		}
		if err := validateYearLaunched(v.YearLaunched); err != nil {
			return err
		}
		if err := validateRating(v.Rating); err != nil {
			return err
		}
		if err := validateDuration(v.Duration); err != nil {
			return err
		}

	case "update":
		if lenTitle == 0 && lenDescription == 0 && v.YearLaunched == 0 && v.Rating == "" && v.Duration == 0 {
			return errors.New("Video must update at least one field")
		} else {
			if lenTitle != 0 {
				if err := validateTitle(lenTitle); err != nil {
					return err
				}
			}

			if lenDescription != 0 {
				if err := validateDescription(lenDescription, descriptionWords); err != nil {
					return err
				}
			}

			if v.YearLaunched != 0 {
				if err := validateYearLaunched(v.YearLaunched); err != nil {
					return err
				}
			}

			if v.Rating != "" {
				if err := validateRating(v.Rating); err != nil {
					return err
				}
			}

			if v.Duration != 0 {
				if err := validateDuration(v.Duration); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func validateTitle(lenTitle int) error {
	if lenTitle < 3 || lenTitle > 255 {
		return errors.New("Title length must be between 3 and 255 characters")
	}
	return nil
}

func validateDescription(lenDescription, descriptionWords int) error {
	if lenDescription < 15 || descriptionWords < 10 {
		return errors.New("Description must have at least 10 words and 15 characters")
	}
	return nil
}

func validateYearLaunched(yearLaunched int) error {
	if yearLaunched < 1895 || yearLaunched > time.Now().Year() {
		return errors.New("Year launched must be between 1895 and current year")
	}
	return nil
}

func validateRating(rating string) error {
	foundRating := false
	for _, ratingOption := range RatingList {
		if ratingOption == rating {
			foundRating = true
		}
	}
	if !foundRating {
		return errors.New("Rating must be a valid value")
	}
	return nil
}

func validateDuration(duration int) error {
	if duration < 1 {
		return errors.New("Duration must be greater than 0")
	}
	return nil
}
