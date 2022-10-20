package changenow

import (
	"code.cryptopower.dev/exchange/lightningswap"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	API_BASE = "https://changenow.io/api/v1/" // API endpoint
	LIBNAME  = "changenow"
)

func init() {
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

// New return an ChangeNow client struct with IDExchange implement
func New(conf lightningswap.ExchangeConfig) (*ChangeNow, error) {
	if conf.ApiKey == "" {
		err := fmt.Errorf("APIKEY is blank")
		return nil, err
	}
	client := lightningswap.NewClient(LIBNAME, &conf)
	return &ChangeNow{client: client, conf: &conf}, nil
}

//ChangeNow represent a ChangeNow client
type ChangeNow struct {
	conf   *lightningswap.ExchangeConfig
	client *lightningswap.Client
	lightningswap.IDExchange
}

//SetDebug set enable/disable http request/response dump
func (c *ChangeNow) SetDebug(enable bool) {
	c.conf.Debug = enable
}

//CalculateExchangeRate get estimate on the amount for the exchange
func (c *ChangeNow) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {
	limits, err := c.QueryLimits(vars.From, vars.To)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	time.Sleep(time.Second * 1)
	estimate, err := c.EstimateAmount(vars)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
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
func (c *ChangeNow) EstimateAmount(vars lightningswap.ExchangeRateRequest) (res lightningswap.EstimateAmount, err error) {
	amountStr := strconv.FormatFloat(vars.Amount, 'f', 8, 64)
	r, err := c.client.Do(API_BASE, "GET",
		fmt.Sprintf("exchange-amount/%s/%s_%s?api_key=%s", amountStr, vars.From, vars.To, c.conf.ApiKey), "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmpRes EstimateAmount
	if err = json.Unmarshal(r, &tmpRes); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	res = lightningswap.EstimateAmount{
		EstimatedAmount:          tmpRes.EstimatedAmount,
		NetworkFee:               tmpRes.NetworkFee,
		ServiceCommission:        tmpRes.ServiceCommission,
		TransactionSpeedForecast: tmpRes.TransactionSpeedForecast,
		WarningMessage:           tmpRes.WarningMessage,
	}

	return
}

//QueryRates (list of pairs LTC-BTC, BTC-LTC, etc)
func (c *ChangeNow) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

//QueryActiveCurrencies get all active currencies
func (c *ChangeNow) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	//vars not used here
	r, err := c.client.Do(API_BASE, "GET", "currencies?active=true", "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmpArr []ActiveCurr
	if err = json.Unmarshal(r, &tmpArr); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	for _, v := range tmpArr {
		currType := "CRYPTO"
		if v.IsFiat {
			currType = "FIAT"
		}
		tmpItem := lightningswap.ActiveCurr{
			CurrencyType: currType,
			Name:         v.Ticker,
		}
		res = append(res, tmpItem)
	}
	return
}

//QueryLimits Get Exchange Rates (from, to)
func (c *ChangeNow) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {
	r, err := c.client.Do(API_BASE, "GET", "min-amount/"+fromCurr+"_"+toCurr, "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmp QueryLimits
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	res = lightningswap.QueryLimits{
		Min: tmp.Min,
	}
	return
}

//CreateOrder create an instant exchange order
func (c *ChangeNow) CreateOrder(orderInfo lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {
	tmpOrderInfo := CreateOrder{
		FromCurrency:      orderInfo.FromCurrency,
		ToCurrency:        orderInfo.ToCurrency,
		ToCurrencyAddress: orderInfo.Destination,
		InvoicedAmount:    orderInfo.InvoicedAmount,
		ExtraID:           orderInfo.ExtraID,
	}
	if tmpOrderInfo.InvoicedAmount == 0.0 {
		err = errors.New(LIBNAME + ":error:createorder invoiced amount is 0")
		return
	}
	payload, err := json.Marshal(tmpOrderInfo)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	r, err := c.client.Do(API_BASE, "POST", "transactions/"+c.conf.ApiKey, string(payload), false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tempItem interface{}
	err = json.Unmarshal(r, &tempItem)
	var tmp CreateResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	res = lightningswap.CreateResultInfo{
		UUID:           tmp.UUID,
		Destination:    tmp.DestinationAddress,
		ExtraID:        tmp.PayinExtraID,
		FromCurrency:   tmp.FromCurrency,
		ToCurrency:     tmp.ToCurrency,
		DepositAddress: tmp.DepositAddress,
	}
	return
}

//UpdateOrder not available for this exchange
func (c *ChangeNow) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	err = errors.New(LIBNAME + ":error:update not available for this exchange")
	return
}

//CancelOrder not available for this exchange
func (c *ChangeNow) CancelOrder(oId string) (res string, err error) {
	err = errors.New(LIBNAME + ":error:cancel not available for this exchange")
	return
}

//OrderInfo get information on orderid/uuid
func (c *ChangeNow) OrderInfo(orderID string) (res lightningswap.OrderInfoResult, err error) {
	r, err := c.client.Do(API_BASE, "GET", "transactions/"+orderID+"/"+c.conf.ApiKey, "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmp OrderInfoResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var amountRecv float64
	if tmp.Status != "finished" {
		amountRecv = tmp.ExpectedAmountReceive
	} else {
		amountRecv = tmp.AmountReceive
	}
	res = lightningswap.OrderInfoResult{
		LastUpdate:     tmp.UpdatedAt,
		ReceiveAmount:  amountRecv,
		TxID:           tmp.PayoutHash,
		Status:         tmp.Status,
		InternalStatus: GetLocalStatus(tmp.Status),
	}
	return
}

//GetLocalStatus translate local status to idexchange status id
//Possible transaction statuses
//new waiting confirming exchanging sending finished failed refunded expired
func GetLocalStatus(status string) lightningswap.Status {
	status = strings.ToLower(status)
	switch status {
	case "finished":
		return lightningswap.OrderStatusCompleted
	case "waiting":
		return lightningswap.OrderStatusWaitingForDeposit
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