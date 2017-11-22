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

type InputData struct {
	money int
	splitOn int
}

func promptForInt(promptMsg string, errorHandler func(e error, input string) int) int {
	fmt.Print(promptMsg)
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]
	retval, e := strconv.Atoi(input)

	if e == nil {
		return retval
	} else {
		return errorHandler(e, input)
	}
}

func handleInputError (e error, input string) int {
	if input == "q" || input == "quit" {
		fmt.Println("QUIT")
		return -1
	} else if input == "a" || input == "av" || input == "available" {
		printCoinsMap("AVAILABLE COINS", coins)
		return -2
	} else {
		fmt.Println("BAD INPUT")
		fmt.Println("ERROR IS: ", e)
		return -2
	}
}

func getInput(promptMsg string) int{
	for {
		retval := promptForInt(promptMsg, handleInputError)
		if retval == -1 {
			return -1
		} else if retval == -2 {
			continue
		} else {
			return retval
		}
	}
}

func acceptCoins(inputDataChan chan InputData, done_processing chan bool) {
	for {
		money := getInput("INSERT MONEY> ")

		if money == -1 {
			close(inputDataChan)
			break
		}
		
		splitOn := getInput("SPLIT ON> ")
		
		if splitOn == -1 {
			close(inputDataChan)
			break
		}

		inputDataChan <- InputData{money, splitOn}
		<-done_processing
	}
}

func processCoins(inputDataChan chan InputData, done_processing chan bool) {
	for data := range inputDataChan {
		money := data.money
		splitOn := data.splitOn
		
		originalCoin := money

		var coinsOut = make(map[int]int)
		var coinsOutAmount int
		var left int

		coinsOutAmount = money / splitOn
		left = money % splitOn

		if coinsOutAmount != 0 {
			if coins[splitOn] >= coinsOutAmount {
				coinsOut[splitOn] = coinsOutAmount
				coins[splitOn] -= coinsOutAmount
				money = left
			} else {
				coinsOut[splitOn] = coins[splitOn]
				money -= coins[splitOn] * splitOn
				coins[splitOn] = 0
			}
		}

		for i := coinTypesN - 1; i >= 0; i-- {

			coinsOutAmount = money / coinTypes[i]
			left = money % coinTypes[i]

			if coinsOutAmount != 0 {
				if coins[coinTypes[i]] >= coinsOutAmount {
					coinsOut[coinTypes[i]] = coinsOutAmount
					coins[coinTypes[i]] -= coinsOutAmount
					money = left
				}
			}
		}

		if money == 0 {
			printCoinsMap("YOUR MONEY SPLITTED", coinsOut)
		} else {
			fmt.Println("Sorry, cannot split ", originalCoin, ". Not enough coins.")
		}

		done_processing <- true
	}
}

func main() {
	var wg sync.WaitGroup

	inputDataChan := make(chan InputData)
	done_processing := make(chan bool)

	fmt.Println("commands are: a - to see available coins")
	fmt.Println("              q - to quit")
	fmt.Println()

	wg.Add(1)
	go func() { defer wg.Done(); acceptCoins(inputDataChan, done_processing) }()

	wg.Add(1)
	go func() { defer wg.Done(); processCoins(inputDataChan, done_processing) }()

	wg.Wait()
}
