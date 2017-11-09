package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

const coinTypesN = 7

var coinTypes = [coinTypesN]int{1, 2, 5, 10, 25, 50, 100}

var coins = map[int]int{
	1:   10,
	2:   4,
	5:   2,
	10:  5,
	25:  6,
	50:  9,
	100: 3,
}

func printCoinsMap(msg string, coins map[int]int) {
	fmt.Printf("--------------\n%s\n--------------\n", msg)
	fmt.Printf("coin | amount\n")
	var coinType int
	var coinAmount int
	for i := range coinTypes {
		coinType = coinTypes[i]
		coinAmount = coins[coinType]
		if coinAmount != 0 {
			fmt.Printf("%-4d | %-4d\n", coinType, coins[coinType])
		}
	}
}

func acceptCoins(coins_ch chan int, done_processing chan bool) {
	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("INSERT MONEY> ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		coin, e := strconv.Atoi(text)

		if e != nil {
			if text == "q" || text == "quit" {
				fmt.Println("QUIT")
				close(coins_ch)
				break
			}

			if text == "a" || text == "av" || text == "available" {
				printCoinsMap("AVAILABLE COINS", coins)
				continue
			} else {
				fmt.Println("BAD INPUT")
				fmt.Println("ERROR IS: ", e)
				continue
			}
		}

		coins_ch <- coin
		<-done_processing
	}
}

func processCoins(coins_ch chan int, done_processing chan bool) {
	for coin := range coins_ch {
		originalCoin := coin

		var coinsOut = make(map[int]int)
		var coinOutAmount int
		var left int

		for i := coinTypesN - 1; i >= 0; i-- {
			// fmt.Println("-- coin ", coin)
			// fmt.Println("   step", coinTypes[i])

			coinOutAmount = coin / coinTypes[i]
			left = coin % coinTypes[i]

			// fmt.Println("   coinOutAmount ", coinOutAmount)
			// fmt.Println("   left", left)

			if coinOutAmount != 0 {
				if coins[coinTypes[i]] >= coinOutAmount {
					coinsOut[coinTypes[i]] = coinOutAmount
					coins[coinTypes[i]] -= coinOutAmount
					coin = left
				}
			}
		}

		if coin == 0 {
			printCoinsMap("YOUR MONEY SPLITTED", coinsOut)
		} else {
			fmt.Println("Sorry, cannot split ", originalCoin, ". Not enough coins.")
		}

		done_processing <- true
	}
}

func main() {
	var wg sync.WaitGroup

	coins_ch := make(chan int)
	done_processing := make(chan bool)

	fmt.Println("commands are: a - to see available coins")
	fmt.Println("              q - to quit")
	fmt.Println()

	wg.Add(1)
	go func() { defer wg.Done(); acceptCoins(coins_ch, done_processing) }()

	wg.Add(1)
	go func() { defer wg.Done(); processCoins(coins_ch, done_processing) }()

	wg.Wait()
}
