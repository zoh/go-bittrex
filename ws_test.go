package bittrex

import (
	"testing"
	"fmt"
)

func TestBittrexSubscribeOrderBook(t *testing.T) {
	bt := New("", "")
	ch := make(chan ExchangeEvent, 16)

	ms, err := bt.GetMarkets()
	if err != nil {
		panic(err)
	}

	// first 100 markets
	var markets []string
	for i := 0; i < len(ms); i++ {
		m := ms[i]
		markets = append(markets, m.MarketName)
	}

RESTART:
	stop := make(chan bool)

	go func() {
		counter := 0
		for {
			select {
			case data := <-ch:
				fmt.Println(data.Method, data.State.MarketName, data.T)
			}

			counter++
		}

	}()

	err = bt.SubscribeExchangeUpdate(markets, "", ch, stop, )
	fmt.Println("err", err)

	//select {
	//case <-time.After(time.Second * 60):
	//	log.Print("timeout")
	//case err := <-errCh:
	//	if err != nil {
	//		log.Print(err)
	//	}
	//}

	goto RESTART
	println("End!")
}
