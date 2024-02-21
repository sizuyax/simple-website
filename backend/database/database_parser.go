package database

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"text/template"
)

var db *gorm.DB

type ParsedText struct {
	gorm.Model
	Title       string `gorm:"primaryKey"`
	Description string
	Image       string
}

func InitDBParser() (*gorm.DB, error) {
	dsn := "host=db_pars user=parser password=parser dbname=db_pars port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&ParsedText{})
	if err != nil {
		panic("failed to migrate database")
	}

	return db, nil
}

func PutTitle(c echo.Context) error {
	var title []ParsedText
	if err := db.Find(&title).Error; err != nil {
		return err
	}

	// Открытие файла HTML-шаблона
	tmpl, err := template.ParseFiles("frontend/html/home_page.html")
	if err != nil {
		return err
	}

	return tmpl.Execute(c.Response().Writer, title)
}

func PutDescription(c echo.Context) error {
	articleID := c.Param("id")

	var description []ParsedText
	var image []ParsedText

	if err := db.Where("id = ?", articleID).First(&description).Error; err != nil {
		return err
	}
	if err := db.Where("id = ?", articleID).First(&image).Error; err != nil {
		return err
	}

	tmpl, err := template.ParseFiles("frontend/html/article.html")
	if err != nil {
		return err
	}

	return tmpl.Execute(c.Response().Writer, description)
}
