package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/katoozi/gin-web-site/internal/pkg/auth"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"

	_ "github.com/lib/pq" // register postgresql driver
)

// Init will intiate the db and load settings from yaml file
func initialPostgres() *sqlx.DB {
	// generate the postgres connect address
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	user := viper.GetString("database.user")
	pass := viper.GetString("database.pass")
	dbName := viper.GetString("database.db.name")
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		user,
		dbName,
		pass,
	)

	// connect to postgresql
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Connect to db Failed: %v", err)
	}
	auth.MigrateTables(db)
	return db
}

func initialRedis() *redis.Client {
	// redis required env variables
	redisDB := viper.GetInt("redis.db")
	redisPass := viper.GetString("redis.pass")

	// determine redis connection type
	redisConfig := viper.GetString("redis.type")
	if redisConfig == "sentinel" {
		// use redis sentinel for high availability and failover
		redisAddrs := viper.GetString("redis.sentinels")
		redisMasterName := viper.GetString("redis.sentinels.master.name")

		sentinelAddrs := strings.Split(redisAddrs, ",")
		client := redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    redisMasterName,
			SentinelAddrs: sentinelAddrs,
			Password:      redisPass,
			DB:            redisDB,
		})
		_, err := client.Ping().Result()
		if err != nil {
			log.Fatalf("Error while connect to redis sentinel: %v\n", err)
		}
		return client

	} else if redisConfig == "simple" {
		// config package for using single redis instance without sentinel or cluster capability
		redisHost := viper.GetString("redis.host")
		redisPort := viper.GetString("redis.port")
		client := redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + redisPort,
			Password: redisPass,
			DB:       redisDB,
		})
		_, err := client.Ping().Result()
		if err != nil {
			log.Fatalf("Error while connect to redis: %v\n", err)
		}
		return client
	}

	log.Fatalf("redis configuration type not found!!!. set the TEST_PROJECT_REDIS_TYPE variable")
	return nil
}

func initialRabbitmq() *amqp.Connection {
	// get rabbitmq configuration variables
	rabbitUser := viper.GetString("rabbitmq.user")
	rabbitPass := viper.GetString("rabbitmq.pass")
	rabbitHost := viper.GetString("rabbitmq.host")
	rabbitPort := viper.GetString("rabbitmq.port")

	// generate server address
	rabbitMQServer := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPass, rabbitHost, rabbitPort)

	conn, err := amqp.Dial(rabbitMQServer)
	if err != nil {
		log.Fatalf("Failed to open rabbitmq connection. %s", err)
	}

	return conn
}
