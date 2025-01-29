package internal

import (
	"gorm.io/gorm"
)

type Link struct {
	Base
	Slug        string `json:"slug" gorm:"unique,index"`
	OriginalUrl string `json:"original_url"`
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
