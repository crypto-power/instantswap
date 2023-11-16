# instantswap

instantswap is a library for instantly trading cryptocurrency.
It uses some instant exchanges api to do trading. 

### Installing

First, use go get to install the latest version of the library.

```
go get github.com/crypto-power/instantswap
```

Next, include instantswap in your application and load more exchange to trade:
In the bellow example. I am loading [flyp.me](https://flyp.me/)

```
import (
    "github.com/crypto-power/instantswap"
    _ "github.com/crypto-power/instantswap/exchange/flypme"
)
```
Now we are supporting exchanges: [changelly](https://changelly.com/), [changenow](https://changenow.io/), 
[coinswitch](https://coinswitch.co/), [fixedfloat](https://fixedfloat.com/), [flypme](https://flyp.me/),
[godex](https://godex.io/), [shapeshift](https://shapeshift.com/), [trocador](https://trocador.app/)


Then you can initial your exchange client:
```
exchange, err := lightningswap.NewExchange("flypme", instantswap.ExchangeConfig{
    Debug:     false,
    ApiKey:    "",
    ApiSecret: "",
})
```
flypme does not require apiKey to access its api. But some other exchanges does.
If you want to use there exchanges you have to get your own apiKey and pass it
to the `ExchangeConfig` params.

### Trading

To start trading you must call: GetExchangeRateInfo.
```
rateInfo, err := exchange.GetExchangeRateInfo(instantswap.ExchangeRateRequest{
    From:   "BTC",
    To:     "DCR",
    Amount: 5,
})
```
`rateInfo` includes the information which will be used to submit an order. 
The full information here: 
``` go
type ExchangeRateInfo struct {
	// Min is the smallest amount will be accepted by the exchange
	Min float64
	// Max is the maximum amount will be accepted by the exchange
	// return Max = 0 means: there are not limited amount
	Max             float64
	ExchangeRate    float64
	EstimatedAmount float64
	MaxOrder        float64
	Signature       string
}
```
`Min` and `Max` are the amount range which will be valid to submit an order. 
If the value of `Max` == 0, the exchange does not have the maximum amount limit.  
`Signature` will be used when you submit your order. Some exchanges used it 
why some others did not.


```go
order, err := exchange.CreateOrder(instantswap.CreateOrder{
    RefundAddress:   "your_btc_address", // if the trading fail, the exchange will refund here
    Destination:     "your_dcr_address", // your received dcr address
    FromCurrency:    "BTC",
    OrderedAmount:   0, // use OrderedAmount or InvoicedAmount
    InvoicedAmount:  0.5,
    ToCurrency:      "DCR",
    ExtraID:         "",
    Signature:       rateInfo.Signature,
    UserReferenceID: "",
    RefundExtraID:   "",
})
```

An order information will be returned. it includes:
```go
type CreateResultInfo struct {
	ChargedFee     float64 `json:"charged_fee,string,omitempty"`
	Destination    string  `json:"destination,omitempty"`
	ExchangeRate   float64 `json:"exchange_rate,string,omitempty"`
	FromCurrency   string  `json:"from_currency,omitempty"`
	InvoicedAmount float64 `json:"invoiced_amount,string,omitempty"`
	OrderedAmount  float64 `json:"ordered_amount,string,omitempty"`
	ToCurrency     string  `json:"to_currency,omitempty"`
	UUID           string  `json:"uuid,omitempty"`
	DepositAddress string
	Expires        int    `json:"expires,omitempty"`
	ExtraID        string `json:"extraId,omitempty"` //changenow.io requirement //changelly payinExtraId value
	PayoutExtraID  string `json:"payoutExtraId,omitempty"`
}
```

Here you must send the currency to `CreateResultInfo.DepositAddress`
The exchange will send the corresponding currency to `CreateOrder.Destination`.
You have to call
```go
exchange.OrderInfo(order.UUID)
```
to know the order's status and get txID to verify the transaction.
