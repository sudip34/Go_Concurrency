package main

import (
	"fmt"
	"strings"
)

// RECEIVE only channel: `<-chan`
// SEND only channel: `chan<-`
func shout(ping <-chan string, pong chan<- string) {
	for {
		s, ok := <-ping

		if !ok {
			// do some thing
			// may be interrapt
		}

		// ping channel sending the string to variable s with `<- channel`
		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(s)) //pong channel receives string here  with `channel <-`
	}
}

func main() {

	//create two channels
	ping := make(chan string)
	pong := make(chan string)

	// run the shout function behind main with go GoRoutine
	go shout(ping, pong)

	// user input
	fmt.Println("Type something adn press ENTER (enter Q to quit)")

	for {
		//print a promt
		fmt.Print("->")

		//get user input
		var userInput string
		_, _ = fmt.Scanln(&userInput)

		if userInput == strings.ToLower("q") {
			break
		}

		ping <- userInput
		//wait for a response
		response := <-pong
		fmt.Println("Response:", response)
	}

	fmt.Println("Done with all. Closeing all chennels")
	close(ping)
	close(pong)

}
