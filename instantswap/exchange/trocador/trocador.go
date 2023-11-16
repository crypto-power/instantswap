package trocador

import (
	"encoding/json"
	"fmt"
	"github.com/crypto-power/instantswap/instantswap"
	"net/http"
	"net/url"
	"strings"
)

const (
	API_BASE = "https://trocador.app/api/"
	LIBNAME  = "trocador"
)

type trocador struct {
	client *instantswap.Client
	conf   *instantswap.ExchangeConfig
}

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

// SetDebug set enable/disable http request/response dump.
func (t *trocador) SetDebug(enable bool) {
	t.conf.Debug = enable
}

// New return a trocador client.
func New(conf instantswap.ExchangeConfig) (*trocador, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := instantswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		return nil
	})
	return &trocador{client: client, conf: &conf}, nil
}

func (t *trocador) currenciesMap() (map[string]instantswap.Currency, error) {
	var form = url.Values{}
	form.Set("api_key", t.conf.ApiKey)
	r, err := t.client.Do(API_BASE, "GET", "coins?"+form.Encode(), "", false)
	if err != nil {
		return nil, err
	}
	var coins []Coin
	err = parseResponseData(r, &coins)
	if err != nil {
		return nil, err
	}
	var mapCurrencies = make(map[string]instantswap.Currency)
	for _, tcdCurr := range coins {
		curr, ok := mapCurrencies[tcdCurr.Ticker]
		if ok {
			curr.Networks = append(curr.Networks, tcdCurr.Network)
		} else {
			curr = instantswap.Currency{
				Name:   tcdCurr.Name,
				Symbol: tcdCurr.Ticker,
				Networks: []string{
					tcdCurr.Network,
				},
			}
		}
		mapCurrencies[tcdCurr.Ticker] = curr
	}
	return mapCurrencies, nil
}

func (t *trocador) coin(ticker string) (*Coin, error) {
	var form = url.Values{}
	form.Set("api_key", t.conf.ApiKey)
	form.Set("ticker", strings.ToLower(ticker))
	r, err := t.client.Do(API_BASE, "GET", "coin?"+form.Encode(), "", false)
	if err != nil {
		return nil, err
	}
	var coins []Coin
	err = parseResponseData(r, &coins)
	if err != nil {
		return nil, err
	}
	if len(coins) == 0 {
		return nil, fmt.Errorf("coin not found")
	}
	return &coins[0], nil
}

func (t *trocador) GetCurrencies() (currencies []instantswap.Currency, err error) {
	mapCurrencies, err := t.currenciesMap()
	for _, curr := range mapCurrencies {
		currencies = append(currencies, curr)
	}
	return
}

func (t *trocador) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	mapCurrencies, err := t.currenciesMap()
	delete(mapCurrencies, strings.ToLower(from))
	for _, curr := range mapCurrencies {
		currencies = append(currencies, curr)
	}
	return
}

func (t *trocador) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	var r []byte
	var form = url.Values{}
	form.Set("api_key", t.conf.ApiKey)
	form.Set("ticker_from", strings.ToLower(vars.From))
	form.Set("ticker_to", strings.ToLower(vars.To))
	form.Set("network_from", vars.FromNetwork)
	form.Set("network_to", vars.ToNetwork)
	form.Set("amount_from", fmt.Sprintf("%.8f", vars.Amount))
	r, err = t.client.Do(API_BASE, "GET", "new_rate?"+form.Encode(), "", false)
	if err != nil {
		return res, err
	}
	var rate Rate
	err = parseResponseData(r, &rate)
	if err != nil {
		return res, err
	}
	coin, err := t.coin(vars.From)
	if err != nil {
		return res, err
	}
	return instantswap.ExchangeRateInfo{
		Min:             coin.Minimum,
		Max:             coin.Maximum,
		ExchangeRate:    rate.rate(),
		EstimatedAmount: rate.AmountTo,
		MaxOrder:        0,
		Signature:       rate.TradeId,
		Provider:        rate.maxProvider(),
	}, nil
}

func (t *trocador) QueryRates(vars interface{}) (res []instantswap.QueryRate, err error) {
	return res, fmt.Errorf("not supported")
}

func (t *trocador) QueryActiveCurrencies(vars interface{}) (res []instantswap.ActiveCurr, err error) {
	return res, fmt.Errorf("not supported")
}

func (t *trocador) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	return res, fmt.Errorf("not supported")
}

func (t *trocador) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	var r []byte
	var form = url.Values{}
	if len(vars.Signature) > 0 {
		form.Set("id", vars.Signature)
	}
	form.Set("api_key", t.conf.ApiKey)
	form.Set("ticker_from", strings.ToLower(vars.FromCurrency))
	form.Set("ticker_to", strings.ToLower(vars.ToCurrency))
	form.Set("network_from", vars.FromNetwork)
	form.Set("network_to", vars.ToNetwork)
	form.Set("amount_from", fmt.Sprintf("%.8f", vars.InvoicedAmount))
	form.Set("address", vars.Destination)
	form.Set("fixed", "True")
	form.Set("refund", vars.RefundAddress)
	form.Set("provider", vars.Provider)
	form.Set("refund_memo", "0")
	r, err = t.client.Do(API_BASE, "GET", "new_trade?"+form.Encode(), "", false)
	if err != nil {
		return res, err
	}
	var trade Trade
	err = parseResponseData(r, &trade)
	if err != nil {
		return res, err
	}

	return instantswap.CreateResultInfo{
		ChargedFee:     0,
		Destination:    trade.AddressUser,
		ExchangeRate:   trade.rate(),
		FromCurrency:   trade.CoinFrom,
		InvoicedAmount: trade.AmountFrom,
		OrderedAmount:  trade.AmountTo,
		ToCurrency:     trade.CoinTo,
		UUID:           trade.TradeId,
		DepositAddress: trade.AddressProvider,
		Expires:        0,
		ExtraID:        "",
		PayoutExtraID:  "",
	}, nil
}

// UpdateOrder accepts orderID value and more if needed per lib.
func (t *trocador) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}
func (t *trocador) CancelOrder(orderID string) (res string, err error) {
	return
}

// OrderInfo accepts orderID value and more if needed per lib.
func (t *trocador) OrderInfo(orderID string) (res instantswap.OrderInfoResult, err error) {
	var r []byte
	var form = url.Values{}
	form.Set("id", orderID)
	form.Set("api_key", t.conf.ApiKey)
	r, err = t.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("trade?%s", form.Encode()),
		"", false)
	if err != nil {
		return
	}
	var trade Trade
	err = parseResponseData(r, &trade)
	if err != nil {
		return
	}

	res = instantswap.OrderInfoResult{
		Expires:        0,
		LastUpdate:     "",
		ReceiveAmount:  0,
		TxID:           trade.Details.tx(),
		Status:         trade.Status,
		InternalStatus: localStatus(trade.Status),
		Confirmations:  "",
	}
	return
}

func (t *trocador) EstimateAmount(vars interface{}) (res instantswap.EstimateAmount, err error) {
	return
}

func parseResponseData(data []byte, obj interface{}) error {
	var err Error
	if json.Unmarshal(data, &err) == nil {
		if len(err.Error) > 0 {
			return fmt.Errorf(err.Error)
		}
	}
	return json.Unmarshal(data, obj)
}
