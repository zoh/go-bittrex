package bittrex

import (
	"testing"
	"fmt"
)

func TestBittrexSubscribeOrderBook(t *testing.T) {
	bt := New("", "")
	ch := make(chan ExchangeEvent, 16)
	//errCh := make(chan error)
	markets := []string{
		"BTC-ETC",
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

			if counter > 3 {
				println("stop and refresh")
				markets = []string{
					"USDT-BTC",
				}
				stop <- true
			}
			counter++
		}

	}()

	err := bt.SubscribeExchangeUpdate(markets, "", ch, stop, )
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
