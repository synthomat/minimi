package internal

import (
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

var modelsToMigrate = []interface{}{
	&Link{},
}

func runMigrations(db *gorm.DB, models []interface{}) {
	err := db.AutoMigrate(models...)

	if err != nil {
		panic("Failed to run migrations: " + err.Error())
	}
}

func NewDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("minimi.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	runMigrations(db, modelsToMigrate)

	return db
}

type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"type:timestamp;not null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"type:timestamp;null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index,null" json:"deleted_at"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New()
	return nil
}
