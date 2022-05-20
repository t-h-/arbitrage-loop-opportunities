package main

import (
	d "arbitrageloop/data"
	u "arbitrageloop/utils"
	"fmt"
	"testing"
)

func TestPop(t *testing.T) {
	s := []string{"0", "1", "2"}
	elem := u.Pop(&s, 0)
	fmt.Println(s)
	fmt.Println(elem)
}

func TestGetExchangeRates(t *testing.T) {
	ers, curs := d.GetExchangeRates()
	fmt.Println(ers)
	fmt.Println(curs)
}

func TestGetSyntheticExchangeRates(t *testing.T) {
	ers, curs := d.GetSyntheticExchangeRates(6, true)
	u.PprintMap(*ers)
	fmt.Println(curs)
}
