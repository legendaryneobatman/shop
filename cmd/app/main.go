package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	todo "go-shop"
	"go-shop/internal/handler"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := sqlx.Open("pgx", os.Getenv("SHARED_DB_URL"))

	if err != nil {
		logrus.Fatalf("failed to connect to db %s", err.Error())
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")
	routesHandler := handler.NewHandler(db)
	routesHandler.Init(router)

	server := new(todo.Server)
	if err := server.Run(viper.GetString("port"), router); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("internal/configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	return viper.ReadInConfig()
}
