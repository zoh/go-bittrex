package bittrex

type BittrexClient interface {
	SetDebug(enable bool)
	GetDistribution(market string) (distribution Distribution, err error)
	GetCurrencies() (currencies []Currency, err error)
	GetMarkets() (markets []Market, err error)
	GetTicker(market string) (ticker Ticker, err error)

	GetMarketSummaries() (marketSummaries []MarketSummary, err error)

	GetMarketSummary(market string) (marketSummary []MarketSummary, err error)

	GetOrderBook(market, cat string, depth int) (orderBook OrderBook, err error)

	GetOrderBookBuySell(market, cat string, depth int) (orderb []Orderb, err error)
	GetMarketHistory(market string) (trades []Trade, err error)

	BuyLimit(market string, quantity, rate float64) (uuid string, err error)
	BuyMarket(market string, quantity float64) (uuid string, err error)
	SellLimit(market string, quantity, rate float64) (uuid string, err error)

	SellMarket(market string, quantity float64) (uuid string, err error)

	CancelOrder(orderID string) (err error)

	GetOpenOrders(market string) (openOrders []Order, err error)

	GetBalances() (balances []Balance, err error)
	GetBalance(currency string) (balance Balance, err error)

	GetDepositAddress(currency string) (address Address, err error)

	Withdraw(address, currency string, quantity float64) (withdrawUuid string, err error)

	GetOrderHistory(market string) (orders []Order, err error)

	GetWithdrawalHistory(currency string) (withdrawals []Withdrawal, err error)

	GetDepositHistory(currency string) (deposits []Deposit, err error)

	GetOrder(order_uuid string) (order Order2, err error)

	// GetTicks is used to get ticks history values for a market.
	GetTicks(market string, interval string) ([]Candle, error)

	GetBTCPrice() (res BTCPrice, err error)

	SubscribeExchangeUpdate(markets []string,
		initialMarketState string,
		dataCh chan<- ExchangeEvent,
		stop <-chan bool,
	) error
}
