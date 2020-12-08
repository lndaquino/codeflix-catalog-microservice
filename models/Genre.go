package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Genre models a genre
type Genre struct {
	ID        string     `json:"id" valid:"uuid" gorm:"type:uuid;primary_key"`
	Name      string     `json:"name" valid:"type(string),required~Genre name is required,stringlength(3|255)~Genre name must be between 3 and 255 characters" gorm:"varchar(255);unique"`
	IsActive  *bool      `json:"is_active" valid:"-" gorm:"bool;default:true"`
	CreatedAt *time.Time `json:"created_at,omitempty" valid:"-" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" valid:"-" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" valid:"-" gorm:"autoDeleteTime"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// Validate validates basic struct
func (g *Genre) Validate() error {
	if _, err := govalidator.ValidateStruct(g); err != nil {
		return err
	}
	return nil
}

// Prepare prepares values
func (g *Genre) Prepare() {
	if g.ID == "" {
		g.ID = uuid.NewV4().String()
	}
	g.Name = html.EscapeString(strings.TrimSpace(g.Name))
}

// Create creates a new genre
func (g *Genre) Create(db *gorm.DB) (*Genre, error) {
	if err := db.Create(&g).Error; err != nil {
		return &Genre{}, err
	}

	return g, nil
}

// FindAll returns all genres in db
func (g *Genre) FindAll(db *gorm.DB) (*[]Genre, error) {
	genres := []Genre{}

	if err := db.Model(&Genre{}).Find(&genres).Error; err != nil {
		return &[]Genre{}, err
	}

	return &genres, nil
}

// FindByID searchs a genre by id
func (g *Genre) FindByID(db *gorm.DB) error {
	var err error

	err = db.Take(&g).Error
	if err != nil {
		return err
	}

	if gorm.IsRecordNotFoundError(err) {
		return errors.New("Genre not found")
	}

	return err
}

// Update updates a genre by id
func (g *Genre) Update(db *gorm.DB) (*Genre, error) {
	rows := db.Model(&g).Updates(&g).RowsAffected
	if rows == 0 {
		return &Genre{}, errors.New("Genre not found")
	}

	return g, nil
}

// Delete deletes a genre by id
func (g *Genre) Delete(db *gorm.DB) error {
	db = db.Delete(&g)
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return errors.New("Genre not found")
	}
	return nil
}
