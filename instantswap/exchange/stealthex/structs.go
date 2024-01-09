package stealthex

import "time"

type Currency struct {
	Symbol            string      `json:"symbol"`
	Network           string      `json:"network"`
	HasExtraId        bool        `json:"has_extra_id"`
	ExtraId           string      `json:"extra_id"`
	Name              string      `json:"name"`
	WarningsFrom      []string    `json:"warnings_from"`
	WarningsTo        []string    `json:"warnings_to"`
	ValidationAddress string      `json:"validation_address"`
	ValidationExtra   interface{} `json:"validation_extra"`
	AddressExplorer   string      `json:"address_explorer"`
	TxExplorer        string      `json:"tx_explorer"`
	Image             string      `json:"image"`
}

type Estimate struct {
	EstimatedAmount float64 `json:"estimated_amount,string"`
	RateId          string  `json:"rate_id"`
}

type Range struct {
	MinAmount float64 `json:"min_amount,string"`
	MaxAmount float64 `json:"max_amount,string"`
}

type OrderRequest struct {
	CurrencyFrom  string  `json:"currency_from"`
	CurrencyTo    string  `json:"currency_to"`
	AddressTo     string  `json:"address_to"`
	ExtraIdTo     string  `json:"extra_id_to"`
	AmountFrom    float64 `json:"amount_from,omitempty,string"`
	AmountTo      float64 `json:"amount_to,omitempty,string"`
	RateId        string  `json:"rate_id"`
	Referral      string  `json:"referral"`
	Fixed         bool    `json:"fixed"`
	Provider      string  `json:"provider"`
	RefundAddress string  `json:"refund_address"`
	RefundExtraId string  `json:"refund_extra_id"`
}

type Order struct {
	Id             string              `json:"id"`
	Type           string              `json:"type"`
	Timestamp      time.Time           `json:"timestamp"`
	UpdatedAt      time.Time           `json:"updated_at"`
	CurrencyFrom   string              `json:"currency_from"`
	CurrencyTo     string              `json:"currency_to"`
	AmountFrom     float64             `json:"amount_from,string"`
	ExpectedAmount string              `json:"expected_amount"`
	AmountTo       float64             `json:"amount_to,string"`
	PartnerFee     interface{}         `json:"partner_fee"`
	AddressFrom    string              `json:"address_from"`
	AddressTo      string              `json:"address_to"`
	ExtraIdFrom    string              `json:"extra_id_from"`
	ExtraIdTo      string              `json:"extra_id_to"`
	TxFrom         string              `json:"tx_from"`
	TxTo           string              `json:"tx_to"`
	Status         string              `json:"status"`
	Currencies     map[string]Currency `json:"currencies"`
	RefundAddress  string              `json:"refund_address"`
	RefundExtraId  string              `json:"refund_extra_id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
