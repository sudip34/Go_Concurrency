package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"subscirption-service/data"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn" //` _ ` is used to ignore as we are not using it directly in this main bur we need to in go mod
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = ":80"

func main() {
	//connect to database
	db := initDB()
	fmt.Println("we are pinging the database")
	db.Ping()

	// create session
	session := initSession()

	//create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create waitgroup
	wg := sync.WaitGroup{}

	//set up the application config
	app := Config{
		Session:  session,
		DB:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Wait:     &wg,
		Models:   data.New(db),
	}

	// set up mail

	// listen for signals
	go app.ListenForShutdown()

	// listen for web connections
	app.serve()
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    webPort,
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting web server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to database")
	}
	return conn
}

// connectToDB tries to connect to postgres, and backs off until a connection
// is made, or we have not connected after 10 tries
func connectToDB() *sql.DB {
	counts := 0

	// dsn := "user=postgres password=password host=localhost port=5432 database=concurrency sslmode=disable connect_timeout=5"

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready...")
		} else {
			log.Print("connected to database!")
			return connection
		}

		if counts > 20 {
			return nil
		}

		log.Print("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++

		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
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
	//register user in to a session
	gob.Register(data.User{})
	//setup
	session := scs.New()
	session.Store = redisstore.New(initRadis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
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

func (app *Config) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	//perform any cleanup-task
	app.InfoLog.Println("would run clean up tasks..")

	//block until Waitgroup is empty
	app.Wait.Wait()
	app.Mailer.DoneCahn <- true

	app.InfoLog.Println("closing channels and shutting down")

	//close all mailer channel
	close(app.Mailer.MailerChan)
	close(app.Mailer.ErrorChan)
	close(app.Mailer.DoneCahn)
}

func (app *Config) createMail() Mail {
	//create channels
	errorChan := make(chan error)
	mailerChan := make(chan Message, 100)
	mailerDoneChan := make(chan bool)

	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        8025,
		Encryption:  "none",
		FromAddress: "info@sudip.com",
		FromName:    "Info",
		Wait:        app.Wait,
		ErrorChan:   errorChan,
		MailerChan:  mailerChan,
		DoneCahn:    mailerDoneChan,
	}

	return m
}
