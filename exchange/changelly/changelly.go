package changelly

import (
	"code.cryptopower.dev/exchange/lightningswap"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	API_BASE = "https://api.changelly.com/" // API endpoint
	LIBNAME  = "changelly"
)

func init() {
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

// New return a Changelly api client
func New(conf lightningswap.ExchangeConfig) (*Changelly, error) {
	client := lightningswap.NewClient(LIBNAME, &conf)
	return &Changelly{
		client: client,
		conf:   &conf,
	}, nil
}

//Changelly represent a Changelly client
type Changelly struct {
	client *lightningswap.Client
	conf   *lightningswap.ExchangeConfig
	lightningswap.IDExchange
}

//SetDebug set enable/disable http request/response dump
func (c *Changelly) SetDebug(enable bool) {
	c.conf.Debug = enable
}
func handleErr(r json.RawMessage) (err error) {
	var errorVal jsonError
	if err = json.Unmarshal(r, &errorVal); err != nil {
		return err
	}
	if errorVal.Message != "" {
		var errorStr string
		errorStr = LIBNAME + ":error: " + errorVal.Message
		err = errors.New(errorStr)
		return
	}
	return nil
}

//CalculateExchangeRate get estimate on the amount for the exchange
func (c *Changelly) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {
	limits, err := c.QueryLimits(vars.From, vars.To)
	if err != nil {
		err = errors.New(err.Error())
		return
	}
	time.Sleep(time.Second * 1)
	estimate, err := c.EstimateAmount(vars)
	if err != nil {
		err = errors.New(err.Error())
		return
	}

	rate := vars.Amount / estimate.EstimatedAmount

	res = lightningswap.ExchangeRateInfo{
		ExchangeRate:    rate,
		Min:             limits.Min,
		Max:             limits.Max,
		EstimatedAmount: estimate.EstimatedAmount,
	}

	return
}

//EstimateAmount get estimate on the amount for the exchange
func (c *Changelly) EstimateAmount(vars lightningswap.ExchangeRateRequest) (res lightningswap.EstimateAmount, err error) {
	amountStr := strconv.FormatFloat(vars.Amount, 'f', 8, 64)
	nonce := strconv.FormatInt(time.Now().Unix(), 10)
	params := map[string]string{"from": strings.ToLower(vars.From), "to": strings.ToLower(vars.To), "amount": amountStr}
	tmpPayload := jsonRequest{
		ID:      "estimateAmount" + nonce,
		JSONRPC: "2.0",
		Method:  "getExchangeAmount",
		Params:  params,
	}

	payload, err := json.Marshal(tmpPayload)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	r, err := c.client.Do(API_BASE, "POST", "", string(payload), true)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	if response.Error != nil {
		err = handleErr(response.Error)
		if err != nil {
			return
		}
	}

	var tmpAmountStr string
	if err = json.Unmarshal(response.Result, &tmpAmountStr); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	exchangeAmount, err := strconv.ParseFloat(tmpAmountStr, 64)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	res = lightningswap.EstimateAmount{
		EstimatedAmount: exchangeAmount,
	}

	return
}

//QueryRates (list of pairs LTC-BTC, BTC-LTC, etc)
func (c *Changelly) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

//QueryActiveCurrencies get all active currencies
func (c *Changelly) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

//QueryLimits Get Exchange Rates (from, to)
func (c *Changelly) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {
	nonce := strconv.FormatInt(time.Now().Unix(), 10)
	params := map[string]string{"from": strings.ToLower(fromCurr), "to": strings.ToLower(toCurr)}
	tmpPayload := jsonRequest{
		ID:      "queryLimits" + nonce,
		JSONRPC: "2.0",
		Method:  "getMinAmount",
		Params:  params,
	}

	payload, err := json.Marshal(tmpPayload)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	r, err := c.client.Do(API_BASE, "POST", "", string(payload), true)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	if response.Error != nil {
		err = handleErr(response.Error)
		if err != nil {
			return
		}
	}

	var tmpMinAmountStr string
	if err = json.Unmarshal(response.Result, &tmpMinAmountStr); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	minAmount, err := strconv.ParseFloat(tmpMinAmountStr, 64)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	res = lightningswap.QueryLimits{
		Min: minAmount,
	}
	return
}

//CreateOrder create an instant exchange order
func (c *Changelly) CreateOrder(orderInfo lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {
	nonce := strconv.FormatInt(time.Now().Unix(), 10)
	amountStr := strconv.FormatFloat(orderInfo.InvoicedAmount, 'f', 8, 64)
	params := map[string]string{
		"from":          strings.ToLower(orderInfo.FromCurrency),
		"to":            strings.ToLower(orderInfo.ToCurrency),
		"address":       orderInfo.Destination,
		"extraId":       orderInfo.ExtraID,
		"amount":        amountStr,
		"refundAddress": orderInfo.RefundAddress,
		"refundExtraId": orderInfo.RefundExtraID,
	}
	tmpPayload := jsonRequest{
		ID:      "createOrder" + nonce,
		JSONRPC: "2.0",
		Method:  "createTransaction",
		Params:  params,
	}
	if orderInfo.InvoicedAmount == 0.0 {
		err = errors.New(LIBNAME + ":error:createorder invoiced amount is 0")
		return
	}
	payload, err := json.Marshal(tmpPayload)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	if c.conf.ApiKey == "" {
		err = errors.New(LIBNAME + ":error: APIKEY is blank")
	}

	r, err := c.client.Do(API_BASE, "POST", "", string(payload), true)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	if response.Error != nil {
		err = handleErr(response.Error)
		if err != nil {
			return
		}
	}

	var tmp CreateResult
	if err = json.Unmarshal(response.Result, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	res = lightningswap.CreateResultInfo{
		UUID:           tmp.UUID,
		Destination:    tmp.PayoutAddress,
		FromCurrency:   tmp.CurrencyFrom,
		ToCurrency:     tmp.CurrencyTo,
		DepositAddress: tmp.PayinAddress,
		ChargedFee:     tmp.ChangellyFee,
		ExtraID:        tmp.PayinExtraID,
		PayoutExtraID:  tmp.PayoutExtraID,
	}
	return
}

//UpdateOrder not available for this exchange
func (c *Changelly) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	err = errors.New(LIBNAME + ":error:update not available for this exchange")
	return
}

//CancelOrder not available for this exchange
func (c *Changelly) CancelOrder(oId string) (res string, err error) {
	err = errors.New(LIBNAME + ":error:cancel not available for this exchange")
	return
}

//OrderInfo get information on orderid/uuid
func (c *Changelly) OrderInfo(orderID string) (res lightningswap.OrderInfoResult, err error) {
	nonce := strconv.FormatInt(time.Now().Unix(), 10)
	params := map[string]string{"id": orderID}
	tmpPayload := jsonRequest{
		ID:      "orderInfo" + nonce,
		JSONRPC: "2.0",
		Method:  "getTransactions",
		Params:  params,
	}

	payload, err := json.Marshal(tmpPayload)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	if c.conf.ApiKey == "" {
		err = errors.New(LIBNAME + ":error: APIKEY is blank")
	}
	r, err := c.client.Do(API_BASE, "POST", "", string(payload), true)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	if response.Error != nil {
		err = handleErr(response.Error)
		if err != nil {
			return
		}
	}

	var tmp []OrderInfoResult
	if err = json.Unmarshal(response.Result, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var finalOrderInfo OrderInfoResult
	for _, v := range tmp {
		if v.UUID == orderID {
			finalOrderInfo = v
		}
	}
	if finalOrderInfo == (OrderInfoResult{}) {
		err = errors.New(LIBNAME + ":error: order info could not be found")
		return
	}
	res = lightningswap.OrderInfoResult{
		ReceiveAmount:  finalOrderInfo.AmountTo,
		Confirmations:  finalOrderInfo.PayinConfirmations,
		TxID:           finalOrderInfo.PayoutHash,
		Status:         finalOrderInfo.Status,
		InternalStatus: GetLocalStatus(finalOrderInfo.Status),
	}
	return
}

//GetLocalStatus translate local status to idexchange status id
func GetLocalStatus(status string) (iStatus lightningswap.Status) {
	status = strings.ToLower(status)
	switch status {
	case "finished":
		return lightningswap.OrderStatusCompleted
	case "waiting":
		return lightningswap.OrderStatusExchanging
	case "confirming":
		return lightningswap.OrderStatusDepositReceived
	case "refunded":
		return lightningswap.OrderStatusRefunded
	case "expired":
		return lightningswap.OrderStatusExpired
	case "new":
		return lightningswap.OrderStatusNew
	case "exchanging":
		return lightningswap.OrderStatusExchanging
	case "sending":
		return lightningswap.OrderStatusSending
	case "failed":
		return lightningswap.OrderStatusFailed
	default:
		return lightningswap.OrderStatusUnknown
	}
}