package store

import (
	"context"

	"gorm.io/gorm"
)

type Option struct {
	gorm.Model
	Name        string `gorm:"index"`
	Description string
	Type        string
	Default     string
	Example     string
	Sources     []Source `gorm:"foreignKey:ID"`
}

type Source struct {
	gorm.Model
	URL string
}

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) Store {
	return Store{
		db: db,
	}
}

func (s Store) AddOptions(ctx context.Context, options []*Option) error {
	result := s.db.WithContext(ctx).Create(options)
	s.db.WithContext(ctx).Save(options)
	return result.Error
}
