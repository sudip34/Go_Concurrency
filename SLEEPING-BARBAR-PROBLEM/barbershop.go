package main

import (
	"time"

	"github.com/fatih/color"
)

type Barbarshop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbarsDoneChan chan bool
	ClientChan      chan string
	Open            bool
}

func (shop *Barbarshop) addBarbar(barber string) {
	isSleeping := false
	color.Yellow("%s goes to the waiting room to check for clients.", barber)

	for {
		// if there are no clients, the barber goes to sleep
		if len(shop.ClientChan) == 0 {
			color.Yellow("There is nothing to do, so %s takes a nap.", barber)
			isSleeping = true
		}
		client, shopOpen := <-shop.ClientChan

		if shopOpen {
			if isSleeping {
				color.Yellow("%s wakes %s up.\n", client, barber)
				isSleeping = false
			}
			// cur hair
			shop.curHair(barber, client)
		} else {
			// shop is closed and send barber home and close this goroutine
			shop.sendBarberHome(barber)
			return
		}
	}
}

func (shop *Barbarshop) curHair(barber, client string) {
	color.Green("%s is cutting %s's hair", barber, client)
	// cutting the hair
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is finished cutting %s's hair", barber, client)
}

func (shop *Barbarshop) sendBarberHome(barber string) {
	color.Cyan("%s is going home.", barber)
	shop.BarbarsDoneChan <- true
}
