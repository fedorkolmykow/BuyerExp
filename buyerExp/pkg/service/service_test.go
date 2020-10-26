package service

import (
	"fmt"
	m "github.com/fedorkolmykow/avitoauto/pkg/api"
	"math"
	"strconv"
	"testing"
	)

func TestShortenSuccess(t *testing.T){
	var reverse int
	input := 19158

	short:= shorten(input)
	for i, r := range short{
		for j, kr:= range m.KeyRunes{
			if kr == r{
				reverse += j * int(math.Pow(float64(len(m.KeyRunes)), float64(len(short) - i - 1)))
				fmt.Println(strconv.Itoa(reverse))
			}
		}
	}
	if input != reverse{
		t.Errorf("unexpected result:\n%d\nexpected:\n%d ", reverse, input)
	}

}
l = 