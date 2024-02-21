package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
	"net/http"
	"time"
	"website/backend/database"
)

func Parse(db *gorm.DB) error {
	for {
		time.Sleep(1 * time.Second)
		res, err := http.Get("https://3dnews.ru/news")
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return err
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return err
		}

		doc.Find(".article-entry .cntPrevWrapper a.entry-header").Each(func(i int, s *goquery.Selection) {
			time.Sleep(1 * time.Second)
			link, exists := s.Attr("href")
			if exists {
				link = "https://3dnews.ru" + link
				parseArticle(link, db)
			}
		})
	}
}

func parseArticle(link string, db *gorm.DB) {
	res, err := http.Get(link)
	if err != nil {
		fmt.Printf("Error fetching article %s: %v\n", link, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Error fetching article %s: status code %d\n", link, res.StatusCode)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Printf("Error parsing article %s: %v\n", link, err)
		return
	}

	// Извлекаем описание со страницы статьи
	description := doc.Find(".entry-body p").Text()
	if description == "" {
		fmt.Printf("Description not found for article %s\n", link)
		return
	}

	// Извлекаем заголовок статьи
	title := doc.Find(".entry-header h1").Text()
	if title == "" {
		fmt.Printf("Title not found for article %s\n", link)
		return
	}

	image, _ := doc.Find(".entry-body .js-mediator-article img").Attr("src")
	if image == "" {
		fmt.Printf("Image not found for article %s\n", link)
		return
	}

	var existingText database.ParsedText
	result := db.Where("title = ?", title).First(&existingText)
	if result.RowsAffected == 0 {
		db.Create(&database.ParsedText{Title: title, Description: description, Image: image})
		fmt.Printf("New text inserted: %s\n", title)
	} else {
		fmt.Printf("Text already exists: %s\n", title)
	}
}
