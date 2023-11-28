package store

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Option struct {
	gorm.Model
	Name        string `gorm:"index"`
	Description string
	Type        string
	Default     string
	Example     string
	Sources     []Source `gorm:"many2many:option_sources"`
}

type Source struct {
	gorm.Model
	URL string
}

type Store struct {
	db *gorm.DB
}

var SearchLimit = 10

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

func (s Store) FindOptions(ctx context.Context, name string) ([]Option, error) {
	options := []Option{}
	// TODO: make sure struct field is used here ? What if it changes
	result := s.db.WithContext(ctx).Limit(SearchLimit).Where("name LIKE ?", "%"+name+"%").Find(&options)
	return options, result.Error
}

func (s Store) GetLastAddedTime(ctx context.Context) (time.Time, error) {
	options := Option{}
	result := s.db.WithContext(ctx).Last(&options)
	return options.CreatedAt, result.Error
}
