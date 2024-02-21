package routes

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"text/template"
	"website/backend/telegram"
)

func HandleLogin(c echo.Context, client *redis.Client) (error, bool) {
	ok, err := UserExists(client, c.FormValue("username"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error checking if user exists: "+err.Error()), false
	}
	if !ok {
		content, err := os.ReadFile("frontend/html/user/user_not_exist.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error reading incorrect_user.html: "+err.Error()), false
		}
		return c.HTML(http.StatusOK, string(content)), false
	} else {
		ok, err = IsCorrectPassword(c, client, c.FormValue("username"), c.FormValue("password"))
		if err != nil {
			return err, false
		}
		if ok {
			return nil, true
		} else {
			content, err := os.ReadFile("frontend/html/user/incorrect_pass.html")
			if err != nil {
				return c.String(http.StatusInternalServerError, "Error reading incorrect_user.html: "+err.Error()), false
			}
			return c.HTML(http.StatusOK, string(content)), false
		}
	}
}

func HandleRegistration(c echo.Context, client *redis.Client) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	exists, err := UserExists(client, username)
	if err != nil {
		return fmt.Errorf("Error checking if user exists: %v", err)
	}
	if exists {
		tmpl, err := template.ParseFiles("frontend/html/user/user_exists.html")
		if err != nil {
			return fmt.Errorf("Error parsing user_exists.html: %v", err)
		}
		err = tmpl.Execute(c.Response().Writer, nil)
		if err != nil {
			return fmt.Errorf("Error executing user_exists.html: %v", err)
		}
		return nil
	} else {
		hashedPassword := hashingPassword(password)

		if len(username) >= 16 && len(hashedPassword) >= 16 || len(username) <= 5 && len(hashedPassword) <= 5 {
			tmpl, err := template.ParseFiles("frontend/html/user/enough_password.html")
			err = tmpl.Execute(c.Response().Writer, nil)
			if err != nil {
				return fmt.Errorf("Error executing user_exists.html: %v", err)
			}
		} else {
			err = client.Set(username, hashedPassword, 0).Err()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error saving user data to Redis"})
			}
			go SendNotification(username)
		}
	}
	return nil
}

func UserExists(client *redis.Client, username string) (bool, error) {
	exists, err := client.Exists(username).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func IsCorrectPassword(c echo.Context, client *redis.Client, username string, password string) (bool, error) {
	storedPassword, err := client.Get(username).Result()

	if err == redis.Nil {
		tmpl, err := template.ParseFiles("frontend/html/user/user_not_exist.html")
		if err != nil {
			return false, fmt.Errorf("Error parsing user_exists.html: %v", err)
		}

		err = tmpl.Execute(c.Response().Writer, nil)
		if err != nil {
			return false, fmt.Errorf("Error executing user_exists.html: %v", err)
		}
	} else if err != nil {
		return false, c.String(http.StatusInternalServerError, "Error checking user: "+err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func hashingPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "Error hashing password"
	}
	return string(hashedPassword)
}

func SendNotification(username string) {
	message := "Новый пользователь зарегистрирован: " + username
	err := telegram.SendTelegramMessage(message)
	if err != nil {
		fmt.Println("Error sending message to Telegram: ", err)
	}
}
