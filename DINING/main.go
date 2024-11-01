package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

//The Dining Philosophers problem is well known in computer science circles.
// Five philosophers, numbered from 0 through 4, live in a house where the talbe
// is laid for them: each philoshper has their own place at the table.
// Their only difficult - breakes those of spagehetti wich has to be eaten
// with two forks. There are two forks next to each plate, so that presents no
// difficulty. As a consequence, however, this means that no two neighbours
// may be eating simutaneously, since there are five philosophers and five forks.
//
//This is a simple implementation of Dijksta's solution to the "Dining
// Philosophers" dilema.

// Philosopher is a struct wich store informatnon about a philosopher.
type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// list of all philosophers
var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Socretes", leftFork: 0, rightFork: 1},
	{name: "Aristotle", leftFork: 1, rightFork: 2},
	{name: "Pascal", leftFork: 2, rightFork: 3},
	{name: "Locke", leftFork: 3, rightFork: 4},
}

// define some variable
var hunger = 3 // how many times does a person eat?
var eatTime = 1 * time.Second
var thinkTime = 1 * time.Second

var sleepTime = 1 * time.Second // time to printout
// var dinnerFinished = make([]string, 5)
var orderMutex sync.Mutex
var dinnerFinished []string

func main() {
	//printout a welcome message
	fmt.Println("Dining Philosopher problems")
	fmt.Println("---------------------------")
	fmt.Println("The table is empty")

	// start the meal
	dine()

	// print out finished message
	time.Sleep(sleepTime)
	fmt.Println("The table is empty")

}

func dine() {
	wg := &sync.WaitGroup{} // wait to make sure that all philosopher have finished dinner
	wg.Add(len(philosophers))

	seated := &sync.WaitGroup{} // wait to make sure that all philosopher are seated on the table
	seated.Add(len(philosophers))

	// forks is a map of all 5 forks. Forks are assigned using the fields leftFork and rightFork
	// in the Philosopher type. Each fork, then, can be found using the index (an integer), and each
	// fork has a unique mutex.
	var forks = make(map[int]*sync.Mutex)

	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	//start the meal
	for i := 0; i < len(philosophers); i++ {
		// finre of a goroutine for the current philosopher
		go diningProblem(philosophers[i], wg, forks, seated)
	}

	// Wait for the philosophers to finish. This blocks until the wait group is 0.
	wg.Wait()
	color.Green("LEAVING the LOOP")

	for index, philosopher := range dinnerFinished {
		fmt.Printf("%s has fished dinner at the position %d \n", philosopher, index+1)
	}

}

func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("%s is seated at the table.\n", philosopher.name)

	// Decremented the seated WaitGroup by one.
	seated.Done()

	//Wait until everyone is seated.
	seated.Wait()

	// eat three times
	for i := hunger; i > 0; i-- {
		//get a lock on the  fork
		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right Frok\n", philosopher.name)

			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left Frok \n", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left Frok \n", philosopher.name)

			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right Frok\n", philosopher.name)
		}
		// By the time we get to this line, the philosopher has a lock (mutex) on both forks.
		fmt.Printf("\t%s has both forks and is eating\n", philosopher.name)
		time.Sleep(eatTime)

		fmt.Printf("\t%s is thinking\n", philosopher.name)
		time.Sleep(thinkTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("\t%s is put down fork", philosopher.name)

	}
	orderMutex.Lock()
	dinnerFinished = append(dinnerFinished, philosopher.name)
	orderMutex.Unlock()

	fmt.Println(philosopher.name, "is satisfied")
	fmt.Println(philosopher.name, "left the table")
}
