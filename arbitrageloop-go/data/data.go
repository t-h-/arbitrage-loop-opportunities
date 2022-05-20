package data

import (
	u "arbitrageloop/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type PathArbitrage struct {
	Path      []string
	Arbitrage float32
}

type CalculationResult struct {
	Mu      sync.Mutex
	Results []PathArbitrage
}

func GetExchangeRates() (*map[[2]string]float32, []string) {
	resp, err := http.Get("https://api.swissborg.io/v1/challenge/rates")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var unmarshalledBody map[string]string
	json.Unmarshal(body, &unmarshalledBody)

	exchangeRates := make(map[[2]string]float32)
	currencies := make([]string, 0)
	for currencyPair, exchangeRate := range unmarshalledBody {
		exchangeRate, err := strconv.ParseFloat(exchangeRate, 32)
		if err != nil {
			log.Fatal(err)
		}
		cs := strings.Split(currencyPair, "-")
		c1 := cs[0]
		c2 := cs[1]
		exchangeRates[[2]string{c1, c2}] = float32(exchangeRate)
		if !u.Contains(currencies, c1) {
			currencies = append(currencies, c1)
		}
	}

	return &exchangeRates, currencies
}

func GetSyntheticExchangeRates(numCurrencies int, simple bool) (*map[[2]string]float32, []string) {
	currencies := []string{}
	for i := 1; i <= numCurrencies; i++ {
		currencies = append(currencies, "C"+fmt.Sprint(i))
	}
	exchangeRates := make(map[[2]string]float32)
	increasedValEvery := 4
	for i, v1 := range currencies {
		for j, v2 := range currencies {
			key := [2]string{v1, v2}
			yek := [2]string{v2, v1}
			if v1 == v2 {
				exchangeRates[key] = float32(1.0)
			} else {
				exchangeRates[key] = 1.0
				exchangeRates[yek] = 1.0
				if !simple {
					var addend float32
					if i%increasedValEvery == 0 {
						addend = 0.1
					}
					exchangeRates[key] = 1.0 + addend
					addend = 0
					if (j+1)%increasedValEvery == 0 {
						addend = 0.1
					}
					exchangeRates[yek] = 1.0 + addend
				}
			}
		}
	}
	if simple {
		exchangeRates[[2]string{"C1", "C2"}] = 1.1
	}
	return &exchangeRates, currencies
}

// all PathArbitrage.Path slices need to be at least len() = 3
func PprintPathArbitrageArr(arr []PathArbitrage) {
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].Path[1] == arr[j].Path[1] {
			return arr[i].Path[2] < arr[j].Path[2]
		} else {
			return arr[i].Path[1] < arr[j].Path[1]
		}
	})
	for _, v := range arr {
		fmt.Println(v)
	}
}
