package entity

import "gorm.io/gorm"

type MotherService struct {
	gorm.Model
	ID   uint64
	Name string
	// TODO: add other fields based on migration files
}
