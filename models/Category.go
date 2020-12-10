package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

// Category models a category
type Category struct {
	ID          string     `json:"id" valid:"uuid" gorm:"type:uuid;primary_key"`
	Name        string     `json:"name" valid:"type(string),required~Category name is required,stringlength(3|255)~Category name must be between 3 and 255 characters" gorm:"varchar(255);unique"`
	Description string     `json:"description" valid:"type(string),stringlength(3|255)~Category name must be between 3 and 255 characters,optional" gorm:"varchar(255)" gorm:"varchar(255)"`
	IsActive    *bool      `json:"is_active" valid:"-" gorm:"bool;default:true"`
	CreatedAt   *time.Time `json:"created_at,omitempty" valid:"-" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" valid:"-" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" valid:"-" gorm:"autoDeleteTime"`
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
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
	c.Description = html.EscapeString(strings.TrimSpace(c.Description))
}

// Create creates a new category
func (c *Category) Create(db *gorm.DB) (*Category, error) {
	if err := db.Create(&c).Error; err != nil {
		return &Category{}, err
	}

	return c, nil
}

// FindAll returns all categories in db
func (c *Category) FindAll(db *gorm.DB) (*[]Category, error) {
	categories := []Category{}

	if err := db.Model(&Category{}).Find(&categories).Error; err != nil {
		return &[]Category{}, err
	}

	return &categories, nil
}

// FindByID searchs a category by id
func (c *Category) FindByID(db *gorm.DB) error {
	var err error

	err = db.Take(&c).Error
	if err != nil {
		return err
	}

	if gorm.IsRecordNotFoundError(err) {
		return errors.New("Category not found")
	}

	return err
}

// Update updates a category by id
func (c *Category) Update(db *gorm.DB) (*Category, error) {
	req := db.Model(&c).Updates(&c).Find(&c)
	if req.Error != nil {
		return &Category{}, errors.New("Internal server error")
	}
	if req.RowsAffected == 0 {
		return &Category{}, errors.New("Category not found")
	}

	return c, nil
}

// Delete deletes a category by id
func (c *Category) Delete(db *gorm.DB) error {
	db = db.Delete(&c)
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return errors.New("Category not found")
	}
	return nil
}
