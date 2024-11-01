package main

import (
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
	shop := Barbarshop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		BarbarsDoneChan: doneChan,
		ClientChan:      clientChan,
		Open:            true,
	}

	// add barbers
	go shop.addBarbar("Frank")

	//start barbarshop as a goroutine

	//add client

	// block until the barbershop is closed

	time.Sleep(5 * time.Second)

}
