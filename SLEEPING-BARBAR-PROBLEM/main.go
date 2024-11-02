package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// variable
var seatingCapacity = 10
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second

func main() {
	// seed our random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	//print welcome message
	color.Red("The Sleeping Barber Problem")
	color.Red("-------------------------")

	//create channels if we need any
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create the barbershop
	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		BarbersDoneChan: doneChan,
		ClientsChan:     clientChan,
		Open:            true,
	}

	// add barbers
	shop.addBarber("Frank")
	shop.addBarber("Jadu")
	shop.addBarber("Kodu")
	shop.addBarber("Frankoo")
	shop.addBarber("Ram")

	//start barbarshop as a goroutine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()

		closed <- true
	}()

	// add clients
	i := 1

	go func() {
		for {
			// get a random number with average arrival rate
			randomMillseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMillseconds)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	// block until the barbershop is closed
	<-closed
	// time.Sleep(5 * time.Second)

}
