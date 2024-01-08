package stealthex

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
}

type Range struct {
	MinAmount float64 `json:"min_amount,string"`
	MaxAmount float64 `json:"max_amount,string"`
}
