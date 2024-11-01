package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	// variable-name chan  Type
	data chan PizzaOrder

	quit chan chan error // chan chan : channels of channels
}

type PizzaOrder struct {
	PizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error) // makes an channel for error if something happens during quit.
	p.quit <- ch
	return <-ch // it should be nil if there is no error
}
func makePizza(pizzaNumber int) *PizzaOrder { // here pizzaNumber will be 0
	pizzaNumber++ // pizzaNUmber = 1
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order $%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		fmt.Printf("Making pizza #%d. It will take %d seconds...\n", pizzaNumber, delay)

		//delay
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** we ran out of ingredients for pizza #%d", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
		}

		p := PizzaOrder{
			PizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
		return &p
	}
	return &PizzaOrder{
		PizzaNumber: pizzaNumber,
	}
}

func pizzaria(pizzaMaker *Producer) {
	//keep track fo which pizza we making
	var i = 0

	//run for ever or until we receive a quit notification
	// try to make pizzas
	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.PizzaNumber
		}

		select {
		// we tried to make a pizza (we sent something to the data channel)
		case pizzaMaker.data <- *currentPizza:

		case quitChan := <-pizzaMaker.quit:

			//close channles
			close(pizzaMaker.data)
			close(quitChan)
			// go out of the all loops
			return
		}
	}
}

func main() {
	//seed the random number generator
	// rand.Seed(time.Now().UnixNano())
	rand.New(rand.NewSource(time.Now().UnixNano()))

	//print out a message
	color.Cyan("The Pizzaria is open for business")
	color.Cyan("---------------------------------")

	//create  a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the backgroung
	go pizzaria(pizzaJob)

	// careate and run the consumer
	for i := range pizzaJob.data {
		if i.PizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order $%d is out for delivery!", i.PizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad")
			}
		} else {
			color.Cyan("Done making pizzas")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing channel!", err)
			}
		}
	}

	//print out the ending message
	color.Cyan("-----------------")
	color.Cyan("Done for the day ")

	color.Cyan("We made %d pizzas, but failed to make %d, with %d attemps in a total.", pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		color.Red("it was a awaful day...")
	case pizzasFailed >= 6:
		color.Red("it was not a very good day...")
	case pizzasFailed >= 4:
		color.Yellow("It was a okay day...")
	case pizzasFailed >= 2:
		color.Yellow("It was a pretty good day...")
	default:
		color.Green("It was a greate day !")

	}
}
