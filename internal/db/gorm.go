package db

import (
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

var modelsToMigrate = []interface{}{
	&Link{},
	&Setting{},
}

func runMigrations(db *gorm.DB, models []interface{}) {
	err := db.AutoMigrate(models...)

	if err != nil {
		panic("Failed to run migrations: " + err.Error())
	}
}

func NewDB(fileName string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	runMigrations(db, modelsToMigrate)

	return db
}

type Base struct {
	Id        uuid.UUID      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null" json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) error {
	base.Id = uuid.New()
	return nil
}
