package bittrex

import (
	"time"
	"strings"
	"fmt"
	"net/http"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"io/ioutil"
	"errors"
	"encoding/json"
	"net/url"
)

func makeSign(apiSecret string, url string) string {
	mac := hmac.New(sha512.New, []byte(apiSecret))
	mac.Write([]byte( url ))
	sig := hex.EncodeToString(mac.Sum(nil))

	return sig
}

// do prepare and process HTTP request to Bittrex API
func (c *client) do20(method string, ressource string, payload string, authNeeded bool) (response []byte, err error) {
	connectTimer := time.NewTimer(c.httpTimeout)

	var rawurl string
	if strings.HasPrefix(ressource, "http") {
		rawurl = ressource
	} else {
		rawurl = fmt.Sprintf("%s%s/%s", API_BASE, API_VERSION, ressource)
	}

	req, err := http.NewRequest(method, rawurl, strings.NewReader(payload))
	if err != nil {
		return
	}
	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Add("Accept", "application/json")

	// Auth
	if authNeeded {
		cookie := http.Cookie{
			Name:  "__RequestVerificationToken",
			Value: "...",
		}
		req.AddCookie(&cookie)
		req.AddCookie(&http.Cookie{
			Name: ".AspNet.ApplicationCookie",
			Value: "....",
		})
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	//fmt.Println(fmt.Sprintf("reponse %s", response), err)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}
	return response, err
}

/*
MarketName:BTC-COVAL
orderId:c2e7a1df-51b1-40b7-8d2b-157a5d4e7cec
*/

func (b *Bittrex) GetOrderHistory2(marketName string) {
	data := url.Values{}
	data.Set("MarketName", marketName)

	r, err := b.client.do20("POST", "https://bittrex.com/api/v2.0/auth/market/GetOrderHistory", data.Encode(), true)
	if err != nil {
		return
	}

	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}

	fmt.Println("=>>", response, err)
}
