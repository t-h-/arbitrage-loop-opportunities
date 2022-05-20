package utils

import (
	"fmt"
)

func Pop(s *[]string, i int) string {
	elem := (*s)[i]
	(*s)[i] = (*s)[len((*s))-1]
	(*s) = (*s)[:len((*s))-1]
	return elem
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func PprintMap(mp map[[2]string]float32) {
	for k, v := range mp {
		fmt.Println(k, v)
	}
}
