package bittrex

import (
	"testing"
	"time"
	"fmt"
	"log"
)

func TestBittrexSubscribeOrderBook(t *testing.T) {
	bt := New("", "")
	ch := make(chan ExchangeState, 16)
	ch2 := make(chan ExchangeSummaryState, 16)
	errCh := make(chan error)

	go func() {

		for {
			var data interface{}
			select {
			case data = <-ch:
			case data = <-ch2:
			}

			fmt.Println(data)
		}

	}()
	go func() {
		errCh <- bt.SubscribeExchangeUpdate("USDT-BTC", ch, ch2, nil)
	}()

	select {
	case <-time.After(time.Second * 60):
		log.Print("timeout")
	case err := <-errCh:
		if err != nil {
			log.Print(err)
		}
	}

	println("End!")
}
