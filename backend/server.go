package backend

import (
	"github.com/labstack/echo/v4"
	"website/backend/database"
	"website/backend/parser"
	"website/backend/routes"
	"website/backend/telegram"
)

func StartServer() {
	database.InitDB()

	db, err := database.InitDBParser()
	if err != nil {
		panic("failed to connect database")
	}
	go parser.Parse(db)

	e := echo.New()

	client := database.GetClient()

	routes.AllRoutes(e, client)

	go telegram.InitTgBot()

	e.Logger.Fatal(e.Start(":1323"))
}
