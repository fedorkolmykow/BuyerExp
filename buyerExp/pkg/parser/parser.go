package parser

import (
	"github.com/fedorkolmykow/avitoexp/pkg/api"
	"math/rand"
	"time"
)

type Parser interface {
	ParsePrices(notices []api.Notice) (ChangedNotices []api.Notice, err error)
}

type parser struct {

}

type mockParser struct {
}

func (p *parser) ParsePrices(notices []api.Notice) (ChangedNotices []api.Notice, err error){
	return
}

func (m *mockParser) ParsePrices(notices []api.Notice) (ChangedNotices []api.Notice, err error){
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, n := range notices{
		if r.Intn(100) % 2 == 0{
			n.Price = r.NormFloat64() * 10000
			ChangedNotices = append(ChangedNotices, n)
		}
	}
	return
}

func NewParser() Parser{
	//p := &parser{}
	p := &mockParser{}
	return p
}