package bittrex

import "time"

type BTCPrice struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  struct {
		Time struct {
			Updated    string    `json:"updated"`
			UpdatedISO time.Time `json:"updatedISO"`
			Updateduk  string    `json:"updateduk"`
		} `json:"time"`
		Disclaimer string `json:"disclaimer"`
		Bpi        struct {
			USD struct {
				Code        string  `json:"code"`
				Rate        string  `json:"rate"`
				Description string  `json:"description"`
				RateFloat   float64 `json:"rate_float"`
			} `json:"USD"`
		} `json:"bpi"`
	} `json:"result"`
}