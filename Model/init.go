package Model

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

var Client *redis.Client
var db *sql.DB

const (
	USERNAME = "root"
	PASSWORD = "lxkdroot"
	NETWORK = "tcp"
	SERVER = "127.0.0.1"
	PORT = 3306
	DATABASE = "blbl"
)

var RedisConn *redis.Options = &redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
}

func Init() {
	var err error
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err = sql.Open("mysql", conn)
	if err != nil {
		log.Fatal("Database Connected Failed!")
	}
	fmt.Println("Mysql Open Succeed!")
	Client = redis.NewClient(RedisConn)
	_, err = Client.Ping().Result()
	if err != nil {
		log.Fatal("Redis ping Failed")
	}
	fmt.Println("Redis Open Succeed!")
}

func Close() {
	var err error
	if err = db.Close(); err != nil {
		log.Fatal("Database Close Failed!")
	}
	if err = Client.Close(); err != nil {
		log.Fatal("Redis Close Failed!")
	}
}




