package routes

import (
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"website/backend/database"
	"website/backend/gmail"
)

func AllRoutes(e *echo.Echo, client *redis.Client) *echo.Echo {
	registrationUser(e, client)
	loginUser(e)
	err := homePage(e, client)
	if err != nil {
		return nil
	}
	err = discrptionArticle(e)
	if err != nil {
		return nil
	}
	err = aboutPage(e)
	if err != nil {
		return nil
	}
	err = contactsPage(e)
	if err != nil {
		return nil
	}

	return e
}

func registrationUser(e *echo.Echo, client *redis.Client) {

	e.POST("/registration", func(c echo.Context) error {
		err := HandleRegistration(c, client)
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusFound, "/")
	})

	e.GET("/registration", func(c echo.Context) error {
		content, err := os.ReadFile("frontend/html/register.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error reading register.html: "+err.Error())
		}
		return c.HTML(http.StatusOK, string(content))
	})
}

func loginUser(e *echo.Echo) {

	e.GET("/", func(c echo.Context) error {
		content, err := os.ReadFile("frontend/html/login.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error reading login.html: "+err.Error())
		}
		return c.HTML(http.StatusOK, string(content))
	})

}

func homePage(e *echo.Echo, client *redis.Client) error {
	e.POST("/home", func(c echo.Context) error {
		err, ok := HandleLogin(c, client)
		if err != nil {
			return err
		}
		if ok {
			err = database.PutTitle(c)
			if err != nil {
				return err
			}
		}
		return nil
	})

	e.GET("/home", func(c echo.Context) error {
		HandleLogin(c, client)
		database.PutTitle(c)
		return nil
	})
	return nil
}

func discrptionArticle(e *echo.Echo) error {
	e.GET("/article/:id", func(c echo.Context) error {
		database.InitDBParser()
		err := database.PutDescription(c)
		if err != nil {
			return err
		}

		return nil
	})
	return nil
}

func aboutPage(e *echo.Echo) error {

	e.GET("/about", func(c echo.Context) error {
		content, err := os.ReadFile("frontend/html/about.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error reading about.html: "+err.Error())
		}
		return c.HTML(http.StatusOK, string(content))
	})
	return nil
}

func contactsPage(e *echo.Echo) error {

	e.GET("/contacts", func(c echo.Context) error {
		content, err := os.ReadFile("frontend/html/contacts.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error reading contacts.html: "+err.Error())
		}
		return c.HTML(http.StatusOK, string(content))
	})

	e.POST("/contacts", func(c echo.Context) error {
		email := c.FormValue("email")
		message := c.FormValue("message")
		name := c.FormValue("name")

		err := gmail.SendMessageGmail(name, email, message)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusFound, "/")
	})

	return nil
}
