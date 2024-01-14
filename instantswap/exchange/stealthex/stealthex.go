package stealthex

import (
	"encoding/json"
	"fmt"
	"github.com/crypto-power/instantswap/instantswap"
	"net/http"
	"strings"
)

const (
	API_BASE = "https://api.stealthex.io/api/v2/"
	LIBNAME  = "stealthex"
)

type stealthex struct {
	client *instantswap.Client
	conf   *instantswap.ExchangeConfig
}

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

// SetDebug set enable/disable http request/response dump.
func (s *stealthex) SetDebug(enable bool) {
	s.conf.Debug = enable
}

// New return a stealthex client.
func New(conf instantswap.ExchangeConfig) (*stealthex, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := instantswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		return nil
	})
	return &stealthex{client: client, conf: &conf}, nil
}

func (s *stealthex) GetCurrencies() (currencies []instantswap.Currency, err error) {
	r, err := s.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("currency?api_key=%s&fixed=boolean", s.conf.ApiKey), "", false)
	if err != nil {
		return nil, err
	}
	var sCurrs []Currency
	err = parseResponseData(r, &sCurrs)
	if err != nil {
		return nil, err
	}
	currencies = make([]instantswap.Currency, len(sCurrs))
	for i, curr := range sCurrs {
		currencies[i] = instantswap.Currency{
			Name:     curr.Name,
			Symbol:   curr.Symbol,
			IsFiat:   false,
			IsStable: false,
			Networks: []string{
				curr.Network,
			},
		}
	}
	return currencies, nil
}

func (s *stealthex) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	r, err := s.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("pairs/%s?api_key=%s", strings.ToLower(from), s.conf.ApiKey), "", false)
	if err != nil {
		return nil, err
	}
	var pairs []string
	err = parseResponseData(r, &pairs)
	if err != nil {
		return nil, err
	}
	for _, toCurr := range pairs {
		currencies = append(currencies, instantswap.Currency{
			Name:     "",
			Symbol:   toCurr,
			IsFiat:   false,
			IsStable: false,
			Networks: nil,
		})
	}
	return currencies, nil
}

func (s *stealthex) estimateAmount(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	r, err := s.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("estimate/%s/%s?api_key=%s&fixed=true&amount=%.8f",
			strings.ToLower(vars.From), strings.ToLower(vars.To), s.conf.ApiKey, vars.Amount), "", false)
	if err != nil {
		return res, err
	}
	var estimate Estimate
	err = parseResponseData(r, &estimate)
	if err != nil {
		return res, err
	}
	res.EstimatedAmount = estimate.EstimatedAmount
	res.ExchangeRate = vars.Amount / estimate.EstimatedAmount
	res.Signature = estimate.RateId
	return res, nil
}

func (s *stealthex) getRange(vars instantswap.ExchangeRateRequest) (*Range, error) {
	body, err := s.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("range/%s/%s?api_key=%s&fixed=true",
			strings.ToLower(vars.From), strings.ToLower(vars.To), s.conf.ApiKey), "", false)
	if err != nil {
		return nil, err
	}
	var r Range
	err = parseResponseData(body, &r)
	return &r, err
}

func (s *stealthex) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	res, err = s.estimateAmount(vars)
	if err != nil {
		return res, err
	}
	r, err := s.getRange(vars)
	if err != nil {
		res.Min = r.MinAmount
		res.Max = r.MaxAmount
	}
	return res, nil
}

func (s *stealthex) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	return res, fmt.Errorf("not supported")
}

func (s *stealthex) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	var req = OrderRequest{
		CurrencyFrom:  vars.FromCurrency,
		CurrencyTo:    vars.ToCurrency,
		AddressTo:     vars.Destination,
		AmountFrom:    vars.InvoicedAmount,
		RateId:        vars.Signature,
		RefundAddress: vars.RefundAddress,
		RefundExtraId: vars.RefundExtraID,
		Provider:      vars.Provider,
		Fixed:         true,
	}
	body, _ := json.Marshal(req)
	r, err := s.client.Do(API_BASE, http.MethodPost, fmt.Sprintf("exchange?api_key=%s", s.conf.ApiKey), string(body), false)
	if err != nil {
		return res, err
	}
	var order Order
	err = parseResponseData(r, &order)
	if err != nil {
		return res, err
	}
	res = instantswap.CreateResultInfo{
		ChargedFee:     0,
		Destination:    order.AddressTo,
		ExchangeRate:   order.AmountFrom / order.AmountTo,
		FromCurrency:   order.CurrencyFrom,
		InvoicedAmount: order.AmountFrom,
		OrderedAmount:  order.AmountTo,
		ToCurrency:     order.CurrencyTo,
		UUID:           order.Id,
		DepositAddress: order.AddressFrom,
		Expires:        0,
		ExtraID:        "",
		PayoutExtraID:  "",
	}
	return res, nil
}

func (s *stealthex) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}
func (s *stealthex) CancelOrder(orderID string) (res string, err error) {
	return
}

func (s *stealthex) OrderInfo(orderID string, extraIds ...string) (res instantswap.OrderInfoResult, err error) {
	r, err := s.client.Do(API_BASE, http.MethodGet, fmt.Sprintf("exchange/%s?api_key=%s", orderID, s.conf.ApiKey), "", false)
	if err != nil {
		return res, err
	}
	var order Order
	err = parseResponseData(r, &order)
	if err != nil {
		return res, err
	}
	res = instantswap.OrderInfoResult{
		Expires:        0,
		LastUpdate:     "",
		ReceiveAmount:  order.AmountTo,
		TxID:           order.TxTo,
		Status:         order.Status,
		InternalStatus: parseStatus(order.Status),
		Confirmations:  "",
	}
	return
}

func (s *stealthex) EstimateAmount(vars interface{}) (res instantswap.EstimateAmount, err error) {
	return
}

func parseResponseData(data []byte, obj interface{}) error {
	var err Error
	if json.Unmarshal(data, &err) == nil {
		if err.Code > 0 {
			return fmt.Errorf(err.Message)
		}
	}
	return json.Unmarshal(data, obj)
}

// waiting, confirming, exchanging, sending, finished, failed, refunded, verifying
func parseStatus(status string) instantswap.Status {
	switch status {
	case "waiting":
		return instantswap.OrderStatusWaitingForDeposit
	case "confirming":
		return instantswap.OrderStatusDepositReceived
	case "exchanging":
		return instantswap.OrderStatusExchanging
	case "sending":
		return instantswap.OrderStatusSending
	case "finished":
		return instantswap.OrderStatusCompleted
	case "failed":
		return instantswap.OrderStatusFailed
	case "refunded":
		return instantswap.OrderStatusRefunded
	case "verifying":
		return instantswap.OrderStatusDepositReceived
	}
	return instantswap.OrderStatusUnknown
}
