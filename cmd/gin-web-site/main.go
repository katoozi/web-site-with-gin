package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/katoozi/gin-web-site/configs"
	"github.com/katoozi/gin-web-site/internal/app/website"
	"github.com/spf13/viper"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func init() {
	configs.SetDefaultValues()

	viper.SetConfigName("config")
	viper.SetConfigFile("configs/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config.yaml file: %v", err)
	}

	gin.SetMode(gin.DebugMode)
}

func main() {

	// database configurations
	databaseConfig := fetchDatabaseConfig()
	dbConnectionStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s sslmode=disable",
		databaseConfig.User,
		databaseConfig.DatabaseName,
		databaseConfig.Password,
	)
	db, err := sqlx.Connect("postgres", dbConnectionStr)
	if err != nil {
		log.Fatalf("Connect to db Failed: %v", err)
	}
	website.MigrateTables(db)

	r := gin.Default()
	r.Static("/static", "./web/assets")
	// r.StaticFS("/more_static", http.Dir("my_file_system"))
	// r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	website.RegisterTemplateFuncs(r)

	// load html files
	// r.LoadHTMLGlob("./web/templates/components/*")
	r.LoadHTMLGlob("./web/templates/*.html")
	//r.LoadHTMLFiles("templates/template1.html", "templates/template2.html")

	website.RegisterRoutes(r)

	// fetch server configs from config.yaml file
	serverConfig := fetchServerConfig()
	r.Run(serverConfig.GetAddr())
	fmt.Println("Start Listning...")
}

func fetchServerConfig() *configs.ServerConfig {
	serverConfig := &configs.ServerConfig{}
	viper.UnmarshalKey("server", &serverConfig)
	return serverConfig
}

func fetchDatabaseConfig() *configs.DatabaseConfig {
	databaseConfig := &configs.DatabaseConfig{}
	viper.UnmarshalKey("database", &databaseConfig)
	return databaseConfig
}
