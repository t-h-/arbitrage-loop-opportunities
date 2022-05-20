package main

import (
	d "arbitrageloop/data"
	u "arbitrageloop/utils"
	"fmt"
	"sync"
	"time"
)

func main() {
	startConcurrencyOnLevel := 2
	maxHops := 8
	numCurrencies := 14
	synthetic := true

	var exchangeRatesS, exchangeRatesA *map[[2]string]float32
	var unvisitedCurrenciesS, unvisitedCurrenciesA []string
	if synthetic {
		exchangeRatesS, unvisitedCurrenciesS = d.GetSyntheticExchangeRates(numCurrencies, false)
		exchangeRatesA, unvisitedCurrenciesA = d.GetSyntheticExchangeRates(numCurrencies, false)
	} else {
		exchangeRatesS, unvisitedCurrenciesS = d.GetExchangeRates()
		exchangeRatesA, unvisitedCurrenciesA = d.GetExchangeRates()
	}

	fmt.Printf("===RESULTS===\nparams: startConcurrencyOnLevel %v, maxHops %v, numCurrencies %v\n", startConcurrencyOnLevel, maxHops, numCurrencies)
	startSync := time.Now()
	resS := RunArbitrageLoop(unvisitedCurrenciesS, exchangeRatesS, maxHops, startConcurrencyOnLevel, false)
	elapsedSync := time.Since(startSync)
	fmt.Println("Sync result length: ", len(resS.Results))

	startAsync := time.Now()
	resA := RunArbitrageLoop(unvisitedCurrenciesA, exchangeRatesA, maxHops, startConcurrencyOnLevel, true)
	elapsedAync := time.Since(startAsync)
	fmt.Println("Async result length:", len(resA.Results))

	fmt.Printf("Sync time elapsed:  %s\n", elapsedSync)
	fmt.Printf("Async time elapsed: %s\n", elapsedAync)
}

func RunArbitrageLoop(currencies []string, exchangeRates *map[[2]string]float32, maxNumHops int, startConcurrencyOnLevel int, async bool) *d.CalculationResult {
	beginningCur := u.Pop(&currencies, 0)
	res := &d.CalculationResult{}
	if async {
		var wg sync.WaitGroup
		backtrackAsync(currencies, []string{beginningCur}, exchangeRates, maxNumHops, startConcurrencyOnLevel, res, &wg, "root")
		wg.Wait()
	} else {
		backtrack(currencies, []string{beginningCur}, exchangeRates, maxNumHops, startConcurrencyOnLevel, res, "root")
	}
	return res
}

func backtrack(unvisitedCurrencies []string, path []string, exchangeRates *map[[2]string]float32, maxNumHops int, startConcurrencyOnLevel int, res *d.CalculationResult, routineId string) {
	//fmt.Println("Entering: ", routineId, "\n    Path:", path)
	level := len(path) // level is synonymous with number of hops.
	if level < maxNumHops {
		for i := range unvisitedCurrencies {
			nextPath, nextUnvisitedCurrencies := getNext(path, unvisitedCurrencies, i)
			arbitrage := calcPathArbitrage(nextPath, exchangeRates)
			addIfValid(res, nextPath, arbitrage)
			newRoutineId := routineId + "_" + fmt.Sprint(i+1)
			backtrack(nextUnvisitedCurrencies, nextPath, exchangeRates, maxNumHops, startConcurrencyOnLevel, res, newRoutineId)
		}
	}
}

func backtrackAsync(unvisitedCurrencies []string, path []string, exchangeRates *map[[2]string]float32, maxNumHops int, startConcurrencyOnLevel int, res *d.CalculationResult, wg *sync.WaitGroup, routineId string) {
	//fmt.Println("Entering: ", routineId, "\n    Path:", path)
	level := len(path) // level is synonymous with number of hops.
	if level < maxNumHops {
		for i := range unvisitedCurrencies {
			nextPath, nextUnvisitedCurrencies := getNext(path, unvisitedCurrencies, i)

			arbitrage := calcPathArbitrage(nextPath, exchangeRates)
			addIfValid(res, nextPath, arbitrage)

			if level == startConcurrencyOnLevel {
				newRoutineId := "go"
				wg.Add(1)
				go backtrackAsync(nextUnvisitedCurrencies, nextPath, exchangeRates, maxNumHops, startConcurrencyOnLevel, res, wg, newRoutineId)
			} else {
				newRoutineId := routineId + "_" + fmt.Sprint(i+1)
				backtrackAsync(nextUnvisitedCurrencies, nextPath, exchangeRates, maxNumHops, startConcurrencyOnLevel, res, wg, newRoutineId)
			}
		}
	}
	if routineId == "go" {
		wg.Done()
		return
	}
}

// calcPathArbitrage takes a sequence of currencies and calculates the circular arbitrage when closing the loop from
// the last currency of the path back to the initial currency of the path
func calcPathArbitrage(path []string, exchangeRates *map[[2]string]float32) float32 {
	if len(path) < 3 {
		return -1
	}
	var arbitrage float32 = 1
	for i := range path {
		iIncMod := (i + 1) % len(path)
		currencyPair := [2]string{path[i], path[iIncMod]}
		arbitrage *= (*exchangeRates)[currencyPair]
	}
	return arbitrage
}

func getNext(path []string, unvisitedCurrencies []string, i int) ([]string, []string) {
	// slice handling in go ist just outrageously NUTS!!!
	nextUnvisitedCurrencies := make([]string, len(unvisitedCurrencies))
	copy(nextUnvisitedCurrencies, unvisitedCurrencies)
	elem := u.Pop(&nextUnvisitedCurrencies, i)
	nextPath := make([]string, len(path)+1)
	copy(nextPath, path)
	nextPath[len(nextPath)-1] = elem
	return nextPath, nextUnvisitedCurrencies
}

func addIfValid(res *d.CalculationResult, path []string, arbitrage float32) {
	if arbitrage <= 1 {
		return
	}
	(*res).Mu.Lock()
	defer (*res).Mu.Unlock()
	(*res).Results = append((*res).Results, d.PathArbitrage{Path: path, Arbitrage: arbitrage})
}
