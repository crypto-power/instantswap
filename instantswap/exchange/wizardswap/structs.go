package wizardswap

type Currency struct {
	Field1                   int         `json:"0"`
	Field2                   string      `json:"1"`
	Field3                   string      `json:"2"`
	Field4                   int         `json:"3"`
	Field5                   string      `json:"4"`
	Field6                   int         `json:"5"`
	Field7                   string      `json:"6"`
	Field8                   int         `json:"7"`
	Field9                   string      `json:"8"`
	Field10                  interface{} `json:"9"`
	Field11                  string      `json:"10"`
	Field12                  string      `json:"11"`
	Field13                  string      `json:"12"`
	Field14                  int         `json:"13"`
	Field15                  string      `json:"14"`
	Field16                  string      `json:"15"`
	Field17                  string      `json:"16"`
	Field18                  string      `json:"17"`
	Field19                  string      `json:"18"`
	Field20                  string      `json:"19"`
	Field21                  int         `json:"20"`
	Id                       int         `json:"id"`
	Symbol                   string      `json:"symbol"`
	Name                     string      `json:"name"`
	Decimals                 int         `json:"decimals"`
	Explorer                 string      `json:"explorer"`
	Minconf                  int         `json:"minconf"`
	Minamt                   string      `json:"minamt"`
	Enabled                  int         `json:"enabled"`
	ValidationAddress        string      `json:"validation_address"`
	ValidationExtra          interface{} `json:"validation_extra"`
	Endpoint                 string      `json:"endpoint"`
	EndpointSockopt          string      `json:"endpoint_sockopt"`
	HashblockEndpoint        string      `json:"hashblock_endpoint"`
	CurrentBlock             int         `json:"current_block"`
	CurrentBlockhash         string      `json:"current_blockhash"`
	HashblockEndpointSockopt string      `json:"hashblock_endpoint_sockopt"`
	FeeAddress               string      `json:"fee_address"`
	NetworkFee               string      `json:"network_fee"`
	Date                     string      `json:"date"`
	DateUpdated              string      `json:"date_updated"`
	Testnet                  int         `json:"testnet"`
}

type Estimate struct {
	EstimatedAmount float64 `json:"estimated_amount,string"`
}

/*currency_from	String	Base currency ticker in lowercase
currency_to	String	Quote currency ticker in lowercase
address_to	String	Recipient blockchain address
amount_from	String	Amount in base currency
refund_address	String	Refund address (optional)
extra_id_to	String	Recipient address Extra ID for currencies that require it (optional)
refund_extra_id	String	Refund Extra ID (optional)
api_key	String	User API key to earn referral fees.*/

type OrderRequest struct {
	CurrencyFrom  string  `json:"currency_from"`
	CurrencyTo    string  `json:"currency_to"`
	AddressTo     string  `json:"address_to"`
	AmountFrom    float64 `json:"amount_from,string"`
	RefundAddress string  `json:"refund_address"`
	ExtraIdTo     string  `json:"extra_id_to"`
	RefundExtraId string  `json:"refund_extra_id"`
	ApiKey        string  `json:"api_key"`
}
