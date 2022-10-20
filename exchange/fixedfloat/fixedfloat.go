package fixedfloat

import (
	"code.cryptopower.dev/exchange/lightningswap"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	API_BASE = "https://fixedfloat.com/api/v1/"
	LIBNAME  = "fixedfloat"
)

// The work on fixedfloat is pending

func init() {
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

type FixedFloat struct {
	conf   *lightningswap.ExchangeConfig
	client *lightningswap.Client
	lightningswap.IDExchange
}

func New(conf lightningswap.ExchangeConfig) (*FixedFloat, error) {
	if conf.ApiKey == "" || conf.ApiSecret == "" {
		return nil, fmt.Errorf("%s:error: api kay and api secret must be provided", LIBNAME)
	}
	client := lightningswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		key := []byte(conf.ApiSecret)
		sig := hmac.New(sha256.New, key)
		sig.Write([]byte(body))
		signedMsg := hex.EncodeToString(sig.Sum(nil))
		r.Header.Set("X-API-SIGN", signedMsg)
		r.Header.Set("X-API-KEY", conf.ApiKey)
		return nil
	})
	return &FixedFloat{client: client, conf: &conf}, nil
}

//SetDebug set enable/disable http request/response dump
func (c *FixedFloat) SetDebug(enable bool) {
	c.conf.Debug = enable
}

func (c *FixedFloat) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {
	var form = make(url.Values)
	form.Set("fromCurrency", vars.From)
	form.Set("toCurrency", vars.To)
	form.Set("type", "fixed")
	form.Set("fromQty", fmt.Sprintf("%.8f", vars.Amount))
	var r []byte
	r, err = c.client.Do(API_BASE, "POST", "getPrice", form.Encode(), false)
	if err != nil {
		return res, err
	}
	fmt.Println(string(r), err)
	return
}

func (c *FixedFloat) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	return res, fmt.Errorf("not supported")
}
func (c *FixedFloat) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	return
}
func (c *FixedFloat) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {
	return
}
func (c *FixedFloat) CreateOrder(vars lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {
	return
}

//UpdateOrder accepts orderID value and more if needed per lib
func (c *FixedFloat) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	return
}
func (c *FixedFloat) CancelOrder(orderID string) (res string, err error) {
	return
}

//OrderInfo accepts orderID value and more if needed per lib
func (c *FixedFloat) OrderInfo(orderID string) (res lightningswap.OrderInfoResult, err error) {
	return
}
func (c *FixedFloat) EstimateAmount(vars interface{}) (res lightningswap.EstimateAmount, err error) {
	return
}

//GetLocalStatus translate local status to idexchange status id
func GetLocalStatus(status string) (iStatus int) {
	// closed, confirming, exchanging, expired, failed, finished, refunded, sending, verifying, waiting
	status = strings.ToLower(status)
	switch status {
	case "wait":
		return 2
	case "confirmation":
		return 3
	case "confirmed":
		return 4
	case "exchanging":
		return 9
	case "sending", "sending_confirmation":
		return 10
	case "success":
		return 1
	case "overdue":
		return 7
	case "error":
		return 11
	case "refunded":
		return 5
	default:
		return 0
	}
}