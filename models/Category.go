package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Category models a category
type Category struct {
	ID          string    `json:"id" valid:"uuid" gorm:"type:uuid;primary_key"`
	Name        string    `json:"name" valid:"type(string),required~Category name is required,stringlength(3|255)~Category name must be between 3 and 255 characters" gorm:"varchar(255);unique"`
	Description string    `json:"description" valid:"type(string),optional" gorm:"varchar(255)"`
	IsActive    *bool     `json:"is_active" valid:"-" gorm:"bool;default:true"`
	CreatedAt   time.Time `json:"created_at" valid:"-" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" valid:"-" gorm:"autoUpdateTime"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// Validate validates basic struct
func (c *Category) Validate() error {
	if _, err := govalidator.ValidateStruct(c); err != nil {
		return err
	}
	return nil
}

// Prepare prepares values
func (c *Category) Prepare() {
	if c.ID == "" {
		c.ID = uuid.NewV4().String()
	}
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
	c.Description = html.EscapeString(strings.TrimSpace(c.Description))
}

// CreateCategory creates a new category
func (c *Category) CreateCategory(db *gorm.DB) (*Category, error) {
	if err := db.Debug().Create(&c).Error; err != nil {
		return &Category{}, err
	}

	return c, nil
}

// FindAllCategories returns all categories in db
func (c *Category) FindAllCategories(db *gorm.DB) (*[]Category, error) {
	categories := []Category{}

	if err := db.Debug().Model(&Category{}).Find(&categories).Error; err != nil {
		return &[]Category{}, err
	}

	return &categories, nil
}

// FindCategoryByID searchs a category by id
func (c *Category) FindCategoryByID(db *gorm.DB) error {
	var err error

	err = db.Debug().Take(&c).Error
	if err != nil {
		return err
	}

	if gorm.IsRecordNotFoundError(err) {
		return errors.New("Category not found")
	}

	return err
}

// UpdateCategory updates a category by id
func (c *Category) UpdateCategory(db *gorm.DB) (*Category, error) {
	var err error

	rows := db.Debug().Model(&c).Updates(&c).RowsAffected
	log.Println(err)
	if rows == 0 {
		return &Category{}, errors.New("Category not found")
	}

	return c, nil
}

// DeleteCategory deletes a category by id
func (c *Category) DeleteCategory(db *gorm.DB) error {
	db = db.Debug().Delete(&c)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
