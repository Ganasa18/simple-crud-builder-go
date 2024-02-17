package domain

import "gorm.io/gorm"

// Model represents the structure of your database table
type Model struct {
	gorm.Model
	Name    string
	Fields  []Field `gorm:"foreignKey:ModelID"`
	ModelID uint    // This field will serve as the foreign key

}

// Field represents the structure of a field in the model
type Field struct {
	Name    string
	Type    string
	ModelID uint
	gorm.Model
}
