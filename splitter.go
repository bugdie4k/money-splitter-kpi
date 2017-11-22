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
	2:   5,
	5:   12,
	10:  5,
	25:  10,
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

type Data struct {
	coin int
	splitOn int
}

func acceptCoins(coins_ch chan Data, done_processing chan bool) {
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

		fmt.Print("SPLIT ON> ")
		text2, _ := reader.ReadString('\n')
		text2 = text2[:len(text2)-1]
		splitOn, e := strconv.Atoi(text2)

		coins_ch <- Data{coin, splitOn}
		<-done_processing
	}
}

func processCoins(coins_ch chan Data, done_processing chan bool) {
	for data := range coins_ch {
		coin := data.coin
		splitOn := data.splitOn
		
		originalCoin := coin

		var coinsOut = make(map[int]int)
		var coinOutAmount int
		var left int

		coinOutAmount = coin / splitOn
		left = coin % splitOn

		if coinOutAmount != 0 {
			if coins[splitOn] >= coinOutAmount {
				coinsOut[splitOn] = coinOutAmount
				coins[splitOn] -= coinOutAmount
				coin = left
			} else {
				coinsOut[splitOn] = coins[splitOn]
				coin -= coins[splitOn] * splitOn
				coins[splitOn] = 0
			}
		}

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

	coins_ch := make(chan Data)
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
