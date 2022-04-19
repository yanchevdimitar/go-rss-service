package database

import (
	"time"

	"gorm.io/gorm"
)

const RSSTableName = "rss"

type RSS struct {
	ID        int64     `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (s *RSS) TableName() string {
	return RSSTableName
}

type RSSRepository interface {
	DeleteAll() *gorm.DB
	Create(rss RSS) *gorm.DB
	Get() (*gorm.DB, []RSS)
}

type defaultRSSRepository struct {
	db *gorm.DB
}

func (rr *defaultRSSRepository) Create(rss RSS) *gorm.DB {
	return rr.db.Create(&rss)
}

func (rr *defaultRSSRepository) DeleteAll() *gorm.DB {
	var rss RSS
	return rr.db.Where("id <> ?", "NULL").Delete(&rss)
}

func (rr *defaultRSSRepository) Get() (*gorm.DB, []RSS) {
	var rss []RSS
	return rr.db.Find(&rss), rss
}

type MySQLRSSRepository struct {
	*defaultRSSRepository
}

func NewMySQLRSSRepository(db *gorm.DB) RSSRepository {
	return &MySQLRSSRepository{&defaultRSSRepository{db}}
}
