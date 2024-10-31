package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

type Income struct {
	Source string
	Amount int
}

func main() {
	// variable for bankbalance
	var bankBalance int
	var balance sync.Mutex

	//print out starting values
	fmt.Printf("Initial account balalnce $%d.00", bankBalance)
	fmt.Println()

	// define weekly revenue

	incomes := []Income{
		{Source: "Maing job", Amount: 500},
		{Source: "Gifts", Amount: 10},
		{Source: "Parttime job", Amount: 50},
		{Source: "Investmens", Amount: 100},
	}

	wg.Add(len(incomes))

	// loop through 52 weeks and print out how much is  made: keep a running total

	for index, income := range incomes {
		go func(week int, income Income) {
			defer wg.Done()
			for week := 1; week <= 52; week++ {
				balance.Lock()
				temp := bankBalance
				temp += income.Amount
				bankBalance = temp
				balance.Unlock()
				fmt.Printf("On week %d, you earned $%d.00 from %s\n", week, bankBalance, income.Source)
			}
		}(index, income)
	}
	wg.Wait()

	// print out final balance
	fmt.Printf("Total income is $%d.00", bankBalance)

}
