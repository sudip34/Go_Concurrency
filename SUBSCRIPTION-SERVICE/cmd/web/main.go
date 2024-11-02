package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn" //` _ ` is used to ignore as we are not using it directly in this main bur we need to in go mod
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {
	//connect to database
	db := initDB()
	db.Ping()

	// create session
	session := initSession()

	// create channels

	// create waitgroup

	//set up the application config

	// set up mail

	// listen for web connections

}

func initDB() *sql.DB {
	conn := connectDB()
	if conn == nil {
		log.Panic("Can't connect to database")
	}
	return conn
}

func connectDB() *sql.DB {
	counts := 0

	dsn := os.Getenv("DSN")

	for {
		connection, err := OpenDB(dsn)
		if err != nil {
			log.Println("Database not yet ready ...")
		} else {
			log.Println("Connect to databse")
			return connection
		}

		if counts > 10 {
			return nil
		}

		log.Print("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++
		continue
	}

}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initSession() *scs.SessionManager {
	session := scs.New()
	session.Store = redisstore.New(initRadis())

}

func initRadis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))

		},
	}
	return redisPool
}
