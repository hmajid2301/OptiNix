package store

import (
	"context"

	"gorm.io/gorm"
)

type Option struct {
	gorm.Model
	// TODO: do we need a unqiue contraint ? or do we just search and return both options
	// Name        string `gorm:"uniqueIndex"`
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

// TODO: move migrate to somewhere else
func New(db *gorm.DB) (Store, error) {
	store := Store{}

	err := db.AutoMigrate(&Option{}, &Source{})
	if err != nil {
		return store, err
	}

	store.db = db
	return store, nil
}

func (s Store) AddOptions(ctx context.Context, options []*Option) error {
	result := s.db.WithContext(ctx).Create(options)
	s.db.WithContext(ctx).Save(options)
	return result.Error
}
