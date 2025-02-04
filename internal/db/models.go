package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Link struct {
	Base
	Slug        string `json:"slug" gorm:"unique,index"`
	OriginalUrl string `gorm:"index" json:"original_url"`
	Description string `json:"description"`
}

func NewLink(slug, originalUrl string) *Link {
	return &Link{
		Slug:        slug,
		OriginalUrl: originalUrl,
	}
}

func LinkBySlug(db *gorm.DB, slug string) (*Link, error) {
	var link Link

	err := db.Where("slug = ?", slug).First(&link).Error

	if err != nil {
		return nil, err
	}

	return &link, nil
}

func LinkById(db *gorm.DB, id string) (*Link, error) {
	var link Link

	err := db.Where("id = ?", uuid.MustParse(id)).First(&link).Error

	if err != nil {
		return nil, err
	}

	return &link, nil
}
