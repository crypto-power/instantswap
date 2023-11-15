package instantswap

type ExchangeConfig struct {
	Debug     bool
	ApiKey    string
	ApiSecret string
	// AffiliateId is used to earn refer coin from transaction
	AffiliateId string
	UserId      string
}

//DECENTRALIZED EXCHANGES
//QUERY

type QueryRate struct {
	Name  string
	Value string
}
type ActiveCurrs struct {
	Name       string
	Currencies []ActiveCurr
}
type ActiveCurr struct {
	ChargedFee       float64 `json:"charged_fee,string,omitempty"`
	Code             string  `json:"code,omitempty"`
	ConfirmationTime int     `json:"confirmation_time,omitempty"`
	CreatedAt        string  `json:"created_at,omitempty"`
	CurrencyType     string  `json:"currency_type,omitempty"`
	Default          bool    `json:"default,omitempty"`
	DisplayPrecision int     `json:"display_precision,omitempty"`
	Exchange         bool    `json:"exchange,omitempty"`
	Name             string  `json:"name,omitempty"`
	Precision        int     `json:"precision,omitempty"`
	Send             bool    `json:"send,omitempty"`
	UpdatedAt        string  `json:"updated_at,omitempty"`
	Website          string  `json:"website,omitempty"`
}

type QueryLimits struct {
	Max float64 `json:"max,string"`
	Min float64 `json:"min,string"`
}

// CREATE
type CreateOrder struct {
	RefundAddress  string  `json:"refund_address"`
	Destination    string  `json:"destination"`
	FromCurrency   string  `json:"from_currency"`
	OrderedAmount  float64 `json:"ordered_amount,string"`  //amount in "to" currency. you want to be received
	InvoicedAmount float64 `json:"invoiced_amount,string"` //amount in "from" currency. you will send it
	ToCurrency     string  `json:"to_currency"`
	FromNetwork    string  `json:"from_network"`
	ToNetwork      string  `json:"to_network"`
	Provider       string  `json:"Provider"` // used for some intermediate exchange

	//changenow.io
	ExtraID string `json:"extraId,omitempty"` //changenow.io requirement
	UserID  string `json:"userId,omitempty"`  //changenow.io partner requirement

	//evercoin
	Signature       string `json:"signature,omitempty"` //evercoin requirement
	UserReferenceID string `json:"userReferenceId,omitempty"`

	//changelly
	RefundExtraID string `json:"refundExtraId,omitempty"`
}
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
type CreateResult struct {
	Expires int              `json:"expires"`
	Order   CreateResultInfo `json:"order"`
}

// UPDATE
type UpdateOrderInfo struct {
	Destination   string  `json:"destination"`
	OrderedAmount float64 `json:"ordered_amount,string"`
	RefundAddress string  `json:"refund_address"`
	UUID          string  `json:"uuid"`
}
type UpdateOrder struct {
	Order UpdateOrderInfo `json:"order"`
}
type UpdateOrderResultInfo struct {
	ChargedFee     float64 `json:"charged_fee,string"`
	Destination    string  `json:"destination"`
	ExchangeRate   float64 `json:"exchange_rate,string"`
	FromCurrency   string  `json:"from_currency"`
	InvoicedAmount float64 `json:"invoiced_amount,string"`
	OrderedAmount  float64 `json:"ordered_amount,string"`
	ToCurrency     string  `json:"to_currency"`
	UUID           string  `json:"uuid"`
}
type UpdateOrderResult struct {
	Expires int                   `json:"expires"`
	Order   UpdateOrderResultInfo `json:"order"`
}

type AcceptOrderResultInfo struct {
	ChargedFee     string `json:"charged_fee"`
	Destination    string `json:"destination"`
	ExchangeRate   string `json:"exchange_rate"`
	FromCurrency   string `json:"from_currency"`
	InvoicedAmount string `json:"invoiced_amount"`
	OrderedAmount  string `json:"ordered_amount"`
	ToCurrency     string `json:"to_currency"`
	UUID           string `json:"uuid"`
	//changenow.io
	ExtraID string `json:"extraId,omitempty"` //changenow.io requirement
}
type AcceptOrderResult struct {
	DepositAddress string                `json:"deposit_address"`
	Expires        int                   `json:"expires"`
	Order          AcceptOrderResultInfo `json:"order"`
}

// OrderInfoResult
type OrderInfoResult struct {
	Expires        int
	LastUpdate     string // should be datetime object
	ReceiveAmount  float64
	TxID           string
	Status         string
	InternalStatus Status
	Confirmations  string
}

type UUID struct {
	UUID string `json:"uuid"`
}

type EstimateAmount struct {
	EstimatedAmount          float64     `json:"estimatedAmount"` //destinationCurrency
	DepositAmount            float64     `json:"depositAmount,omitempty"`
	NetworkFee               float64     `json:"networkFee,omitempty"`
	ServiceCommission        float64     `json:"serviceCommission,omitempty"`
	TransactionSpeedForecast string      `json:"transactionSpeedForecast,omitempty"`
	WarningMessage           interface{} `json:"warningMessage,omitempty"`
	FromCurrency             string      `json:"fromCurrency,omitempty"`
	ToCurrency               string      `json:"toCurrency,omitempty"`
	Signature                string      `json:"signature,omitempty"`
}

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
	Provider        string // used for some intermediate exchange
}

type Status int

const (
	OrderStatusUnknown Status = iota
	OrderStatusCompleted
	OrderStatusWaitingForDeposit
	OrderStatusDepositReceived
	OrderStatusDepositConfirmed
	OrderStatusRefunded
	OrderStatusCanceled
	OrderStatusExpired
	OrderStatusNew
	OrderStatusExchanging
	OrderStatusSending
	OrderStatusFailed
)

func (s Status) String() string {
	switch s {
	case OrderStatusCompleted:
		return "Completed"
	case OrderStatusWaitingForDeposit:
		return "Waiting for deposit"
	case OrderStatusDepositReceived:
		return "Deposit received"
	case OrderStatusDepositConfirmed:
		return "Deposit confirmed"
	case OrderStatusRefunded:
		return "Refunded"
	case OrderStatusCanceled:
		return "Canceled"
	case OrderStatusExpired:
		return "Expired"
	case OrderStatusNew:
		return "New"
	case OrderStatusExchanging:
		return "Exchanging"
	case OrderStatusSending:
		return "Sending"
	case OrderStatusFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

type Currency struct {
	Name     string
	Symbol   string
	IsFiat   bool
	IsStable bool
	Networks []string
}
