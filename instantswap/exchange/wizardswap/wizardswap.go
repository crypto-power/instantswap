// the work on wizardswap is incompleted.
// I tried to create an order(exchange in wizardswap) but the server
// was always return "/tfalse" in the response body.
// I tried to contact support but they do not response yet
// wizardswap is leaved behide temporary

package wizardswap

import (
	"encoding/json"
	"fmt"
	"github.com/crypto-power/instantswap/instantswap"
	"net/http"
	"strings"
)

const (
	API_BASE = "https://www.wizardswap.io/api/"
	LIBNAME  = "wizardswap"
)

type wizardswap struct {
	client *instantswap.Client
	conf   *instantswap.ExchangeConfig
}

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

// SetDebug set enable/disable http request/response dump.
func (w *wizardswap) SetDebug(enable bool) {
	w.conf.Debug = enable
}

// New return a wizardswap client.
func New(conf instantswap.ExchangeConfig) (*wizardswap, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := instantswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		return nil
	})
	return &wizardswap{client: client, conf: &conf}, nil
}

func (w *wizardswap) GetCurrencies() (currencies []instantswap.Currency, err error) {
	r, err := w.client.Do(API_BASE, http.MethodGet, "currency", "", false)
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
			Networks: nil,
		}
	}
	return currencies, nil
}

func (w *wizardswap) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	from = strings.ToLower(from)
	r, err := w.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("pairs/%s", from), "", false)
	if err != nil {
		return nil, err
	}
	var pairs []string
	err = parseResponseData(r, &pairs)
	if err != nil {
		return nil, err
	}
	for _, toCurr := range pairs {
		if toCurr == from {
			continue
		}
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

func (w *wizardswap) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	f := map[string]string{
		"currency_from": strings.ToLower(vars.From),
		"currency_to":   strings.ToLower(vars.To),
		"amount_from":   fmt.Sprintf("%.8f", vars.Amount),
		"api_key":       w.conf.ApiKey,
	}
	data, _ := json.Marshal(f)
	r, err := w.client.Do(API_BASE, http.MethodPost, "estimate", string(data), false)
	if err != nil {
		return res, err
	}
	var estimate Estimate
	err = parseResponseData(r, &estimate)
	if err != nil {
		return res, err
	}
	res.EstimatedAmount = estimate.EstimatedAmount
	res.ExchangeRate = estimate.EstimatedAmount / vars.Amount
	return res, nil
}

func (w *wizardswap) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	return res, fmt.Errorf("not supported")
}

func (w *wizardswap) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	return res, fmt.Errorf("not supported yet")
}

func (w *wizardswap) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}
func (w *wizardswap) CancelOrder(orderID string) (res string, err error) {
	return
}

func (w *wizardswap) OrderInfo(orderID string, extraIds ...string) (res instantswap.OrderInfoResult, err error) {
	return res, fmt.Errorf("not supported yet")
}

func (w *wizardswap) EstimateAmount(vars interface{}) (res instantswap.EstimateAmount, err error) {
	return
}

func parseResponseData(data []byte, obj interface{}) error {
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
