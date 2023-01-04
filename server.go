package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thepaeth/assessment/pkg/expenses"
)

func main() {
	expenses.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/expenses", expenses.CreateExpense)

	log.Fatal(e.Start(os.Getenv("PORT")))

	// fmt.Println("Please use server.go for main file")
	// fmt.Println("start at port:", os.Getenv("PORT"))
}
