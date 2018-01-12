package bittrex

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"github.com/zoh/signalr_bittrex"
)

type OrderUpdate struct {
	Orderb
	Type int
}

type Fill struct {
	Orderb
	Price     decimal.Decimal `json:",omitempty"`
	OrderType string
	Timestamp jTime
}

// ExchangeState contains fills and order book updates for a market.
type ExchangeState struct {
	MarketName string
	Nounce     int
	Buys       []OrderUpdate
	Sells      []OrderUpdate
	Fills      []Fill
	Initial    bool
}

type ExchangeSummaryState struct {
	Deltas []MarketSummary
	Nounce int
}

// doAsyncTimeout runs f in a different goroutine
//	if f returns before timeout elapses, doAsyncTimeout returns the result of f().
//	otherwise it returns "operation timeout" error, and calls tmFunc after f returns.
func doAsyncTimeout(f func() error, tmFunc func(error), timeout time.Duration) error {
	errs := make(chan error)
	go func() {
		err := f()
		select {
		case errs <- err:
		default:
			if tmFunc != nil {
				tmFunc(err)
			}
		}
	}()
	//if timeout == 0 {
	//	timeout = time.Hour * 12
	//}
	select {
	case err := <-errs:
		return err
	case <-time.After(timeout):
		return errors.New("operation timeout")
	}
}

func sendStateAsync(dataCh chan<- ExchangeEvent, st ExchangeEvent) {
	select {
	case dataCh <- st:
	default:
	}
}

func subForMarket(client *signalr.Client, markets []string, initialQueryState string) (json.RawMessage, error) {
	for _, m := range markets {
		_, err := client.CallHub(WS_HUB, "SubscribeToExchangeDeltas", m)
		if err != nil {
			log.Println("Error:", err)
		}
	}

	if initialQueryState != "" {
		return client.CallHub(WS_HUB, "QueryExchangeState", initialQueryState)
	} else {
		return nil, nil
	}
}

func parseStates(messages []json.RawMessage, dataCh chan<- ExchangeEvent) {
	for _, msg := range messages {
		var st ExchangeState
		if err := json.Unmarshal(msg, &st); err != nil {
			continue
		}
		sendStateAsync(dataCh, ExchangeEvent{
			Method: UpdateExchangeState,
			T:      time.Now(),
			State:  st,
		})
	}
}

func parseSummaryState(messages []json.RawMessage, dataCh chan<- ExchangeEvent) {
	for _, msg := range messages {
		var st ExchangeSummaryState
		if err := json.Unmarshal(msg, &st); err != nil {
			log.Println(err)
			continue
		}
		select {
		case dataCh <- ExchangeEvent{
			Method:       UpdateSummaryState,
			T:            time.Now(),
			SummaryState: st}:
		default:
		}
	}
}

type ExchangeEventMethod string

const UpdateExchangeState ExchangeEventMethod = "updateExchangeState"
const UpdateSummaryState ExchangeEventMethod = "updateSummaryState"
const QueryExchangeState ExchangeEventMethod = "queryExchangeState"

type ExchangeEvent struct {
	Method ExchangeEventMethod
	T      time.Time

	State        ExchangeState
	SummaryState ExchangeSummaryState
}

// SubscribeExchangeUpdate subscribes for updates of the market.
// Updates will be sent to dataCh.
// To stop subscription, send to, or close 'stop'.
func (b *Bittrex) SubscribeExchangeUpdate(
	markets []string,
	initialMarketState string,
	dataCh chan<- ExchangeEvent,
	stop <-chan bool,
) error {
	const timeout = 10 * time.Minute
	client := signalr.NewWebsocketClient()
	client.OnClientMethod = func(hub string, method string, messages []json.RawMessage) {
		if hub != WS_HUB {
			return
		}
		switch method {
		case "updateExchangeState":
			parseStates(messages, dataCh)

		case "updateSummaryState":
			parseSummaryState(messages, dataCh)
		}
	}
	err := doAsyncTimeout(func() error {
		return client.Connect("https", WS_BASE, []string{WS_HUB})
	}, func(err error) {
		if err == nil {
			client.Close()
		}
	}, timeout)
	if err != nil {
		return err
	}
	defer client.Close()

	var msg json.RawMessage
	err = doAsyncTimeout(func() error {
		var err error
		msg, err = subForMarket(client, markets, initialMarketState)
		return err
	}, nil, timeout)
	if err != nil {
		return err
	}

	if msg != nil {
		var st ExchangeState
		if err = json.Unmarshal(msg, &st); err != nil {
			return err
		}
		st.Initial = true
		st.MarketName = initialMarketState
		sendStateAsync(dataCh, ExchangeEvent{
			Method: QueryExchangeState,
			T:      time.Now(),
			State:  st,
		})
	}

	println("Wait!")
	select {
	case <-stop:
		println("Stop!")
		return nil
	case <-client.DisconnectedChannel:
		return DisconnectedChannel
	}
	return nil
}

var DisconnectedChannel = errors.New("Disconnect channel signalR")
