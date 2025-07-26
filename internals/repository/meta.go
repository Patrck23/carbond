package repository

import "gorm.io/gorm"

type Excecute interface {
	Begin() Excecute
	Commit() error
	Rollback()
	Exec(query string, args ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
}

type GormDatabase struct {
	db *gorm.DB
}

func NewExcecute(db *gorm.DB) Excecute {
	return &GormDatabase{db: db}
}

func (g *GormDatabase) Begin() Excecute {
	return &GormDatabase{db: g.db.Begin()}
}

func (g *GormDatabase) Commit() error {
	return g.db.Commit().Error
}

func (g *GormDatabase) Rollback() {
	g.db.Rollback()
}

func (g *GormDatabase) Exec(query string, args ...interface{}) *gorm.DB {
	return g.db.Exec(query, args...)
}

func (g *GormDatabase) Create(value interface{}) *gorm.DB {
	return g.db.Create(value)
}
