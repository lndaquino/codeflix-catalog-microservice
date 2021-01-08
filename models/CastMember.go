package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

// TypeDirector defines Director value
const TypeDirector int = 1

// TypeActor defines Actor value
const TypeActor int = 2

// CastMember models a cast member
type CastMember struct {
	ID        string     `json:"id" valid:"uuid" gorm:"type:uuid;primary_key"`
	Name      string     `json:"name" valid:"type(string),required~CastMember name is required,stringlength(3|255)~Category name must be between 3 and 255 characters" gorm:"varchar(255);unique"`
	Type      int        `json:"type" valid:"type(int),required~CastMember type is required"`
	CreatedAt *time.Time `json:"created_at,omitempty" valid:"-" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" valid:"-" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" valid:"-" gorm:"autoDeleteTime"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// Validate validates basic struct
func (c *CastMember) Validate() error {
	log.Println("Validating CastMember...")
	log.Println(c)
	if _, err := govalidator.ValidateStruct(c); err != nil {
		return err
	}

	log.Println("CastMember validated sucessfully!")
	// if c.Type != TypeDirector && c.Type != TypeDirector {
	// 	return errors.New("Invalid CastMember Type")
	// }

	return nil
}

// Prepare prepares values
func (c *CastMember) Prepare() {
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
}

// Create creates a new cast member
func (c *CastMember) Create(db *gorm.DB) (*CastMember, error) {
	log.Println("Chegou aqui")
	if err := db.Create(&c).Error; err != nil {
		return &CastMember{}, err
	}

	return c, nil
}

// FindAll returns all cast members in db
func (c *CastMember) FindAll(db *gorm.DB) (*[]CastMember, error) {
	castMembers := []CastMember{}

	if err := db.Model(&CastMember{}).Find(&castMembers).Error; err != nil {
		return &[]CastMember{}, err
	}

	return &castMembers, nil
}

// FindByID searchs a CastMember by id
func (c *CastMember) FindByID(db *gorm.DB) error {
	var err error

	err = db.Take(&c).Error
	if err != nil {
		return err
	}

	if gorm.IsRecordNotFoundError(err) {
		return errors.New("CastMember not found")
	}

	return err
}

// Update updates a CastMember by id
func (c *CastMember) Update(db *gorm.DB) (*CastMember, error) {
	req := db.Model(&c).Updates(&c).Find(&c)
	if req.Error != nil {
		return &CastMember{}, errors.New("Internal server error")
	}
	if req.RowsAffected == 0 {
		return &CastMember{}, errors.New("CastMember not found")
	}

	return c, nil
}

// Delete deletes a CastMember by id
func (c *CastMember) Delete(db *gorm.DB) error {
	db = db.Delete(&c)
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return errors.New("CastMember not found")
	}
	return nil
}
