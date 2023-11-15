package trocador

import (
	"github.com/crypto-power/instantswap/instantswap"
	"strings"
	"time"
)

type Coin struct {
	Name    string  `json:"name"`
	Ticker  string  `json:"ticker"`
	Network string  `json:"network"`
	Memo    bool    `json:"memo"`
	Image   string  `json:"image"`
	Minimum float64 `json:"minimum"`
	Maximum float64 `json:"maximum"`
}

type Error struct {
	Error string `json:"error"`
}

//- trade_id: the trade ID with us;
//- date: date and time of creation;
//- ticker_from: ticker of the coin to be sold;
//- ticker_to: ticker of the coin to be bought;
//- coin_from: name of coin to be sold;
//- coin_to: name of coin to be bought;
//- network_from: network of coin to be sold;
//- network_to: network of coin to be bought;
//- amount_from: amount of coin to be sold;
//- amount_to: amount of coin to be bought;
//- provider: exchange with the best rate;
//- fixed: rate type for the best rate, True for fixed and False for floating;
//- status: status of the trade;
//- quotes: list of all the other quotes generated, with their KYC rating and waste(spread) in percentage;
//- payment: True or False, depending if it is a standard swap or payment;

type Rate struct {
	TradeId     string  `json:"trade_id"`
	Date        string  `json:"date"`
	TickerFrom  string  `json:"ticker_from"`
	TickerTo    string  `json:"ticker_to"`
	CoinFrom    string  `json:"coin_from"`
	CoinTo      string  `json:"coin_to"`
	NetworkFrom string  `json:"network_from"`
	NetworkTo   string  `json:"network_to"`
	AmountFrom  float64 `json:"amount_from"`
	AmountTo    float64 `json:"amount_to"`
	Provider    string  `json:"provider"`
	Fixed       bool    `json:"fixed"`
	Status      string  `json:"status"`
	Quotes      struct {
		Quotes []Quote `json:"quotes"`
	} `json:"quotes"`
	Payment bool `json:"payment"`
}

func (r *Rate) maxProvider() string {
	for _, q := range r.Quotes.Quotes {
		if q.Waste == 0 {
			return q.Provider
		}
	}
	return r.Provider
}

func (r *Rate) rate() float64 {
	return r.AmountTo / r.AmountFrom
}

type Quote struct {
	Provider  string  `json:"provider"`
	KycRating string  `json:"kycrating"`
	LogPolicy string  `json:"logpolicy"`
	Insurance int     `json:"insurance"`
	Fixed     string  `json:"fixed"`
	AmountTo  float64 `json:"amount_to,string"`
	Waste     float64 `json:"waste,string"`
	Eta       float64 `json:"eta"`
}

type Trade struct {
	TradeId             string    `json:"trade_id"`
	Date                time.Time `json:"date"`
	TickerFrom          string    `json:"ticker_from"`
	TickerTo            string    `json:"ticker_to"`
	CoinFrom            string    `json:"coin_from"`
	CoinTo              string    `json:"coin_to"`
	NetworkFrom         string    `json:"network_from"`
	NetworkTo           string    `json:"network_to"`
	AmountFrom          float64   `json:"amount_from"`
	AmountTo            float64   `json:"amount_to"`
	Provider            string    `json:"provider"`
	Fixed               bool      `json:"fixed"`
	Payment             bool      `json:"payment"`
	Status              string    `json:"status"`
	AddressProvider     string    `json:"address_provider"`
	AddressProviderMemo string    `json:"address_provider_memo"`
	AddressUser         string    `json:"address_user"`
	AddressUserMemo     string    `json:"address_user_memo"`
	RefundAddress       string    `json:"refund_address"`
	RefundAddressMemo   string    `json:"refund_address_memo"`
	Password            string    `json:"password"`
	IdProvider          string    `json:"id_provider"`
	Quotes              struct {
		Support struct {
			TxUrl      string `json:"tx_url"`
			SupportUrl string `json:"support_url"`
			Type       string `json:"type"`
			Tos        string `json:"tos"`
		} `json:"support"`
		ExpiresAt string `json:"expiresAt"`
	} `json:"quotes"`
	Details          TradeDetail `json:"details"`
	AffiliatePartner string      `json:"affiliate_partner"`
}

type TradeDetail struct {
	Webhook            string      `json:"webhook"`
	Hashout            interface{} `json:"hashout"`
	MarketrateCreation float64     `json:"marketrate_creation"`
	AmountBtc          float64     `json:"amount_btc"`
	Support            struct {
		TxUrl      string `json:"tx_url"`
		SupportUrl string `json:"support_url"`
		Type       string `json:"type"`
		Tos        string `json:"tos"`
	} `json:"support"`
	ExpiresAt string `json:"expiresAt"`
}

func (t *TradeDetail) tx() string {
	if t.Hashout == nil {
		return ""
	}
	return t.Hashout.(string)
}

func (t *Trade) rate() float64 {
	return t.AmountTo / t.AmountFrom
}

//- new: you have rates, but did not create the swap yet;
//- waiting: you created the swap but no deposit was detected;
//- confirming: deposit was detected and is yet to be confirmed;
//- sending: deposit confirmed and provider is sending the coins;
//- finished: there is already a payment hash to the user;
//- failed: something might have happened to the swap, please contact support;
//- expired: payment time expired;
//- halted: some issue happened with the swap, please contact support;
//- refunded: exchange claims to have refunded the user;

var (
	statusMap = map[string]instantswap.Status{
		"new":        instantswap.OrderStatusNew,
		"waiting":    instantswap.OrderStatusWaitingForDeposit,
		"confirming": instantswap.OrderStatusWaitingForDeposit,
		"sending":    instantswap.OrderStatusSending,
		"finished":   instantswap.OrderStatusCompleted,
		"failed":     instantswap.OrderStatusFailed,
		"expired":    instantswap.OrderStatusExpired,
		"halted":     instantswap.OrderStatusFailed,
		"refunded":   instantswap.OrderStatusRefunded,
	}
)

func localStatus(status string) instantswap.Status {
	localStatus, ok := statusMap[strings.ToLower(status)]
	if ok {
		return localStatus
	}
	return instantswap.OrderStatusUnknown
}
