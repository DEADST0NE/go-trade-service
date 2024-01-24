package broker

type SsParams struct {
	Symbols []string `json:"symbols"`
}

type SsSubMessage struct {
	Method string   `json:"method"`
	Ch     string   `json:"ch"`
	Params SsParams `json:"params"`
}

type CoreOptions struct {
	Env string `json:"env"`
}

type CoreSubMessage struct {
	Method  string      `json:"method"`
	Options CoreOptions `json:"options"`
}

type CoreMsgTrades struct {
	RequestId   string  `json:"requestId"`
	ClientId    string  `json:"clientId"`
	Environment string  `json:"environment"`
	Timestamp   string  `json:"timestamp"`
	Price       float64 `json:"price"`
	Quantity    float64 `json:"quantity"`
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"`
	order_id    string
	Id          string `json:"id"`
}

type SsMsgTrades struct {
	Ch       string       `json:"ch"`
	Update   *SsTradeList `json:"update,omitempty"`
	Snapshot *SsTradeList `json:"snapshot,omitempty"`
}

type SsTradeList map[string][]SsTrade

type SsTrade struct {
	T int64  `json:"t"`
	I int64  `json:"i"`
	P string `json:"p"`
	Q string `json:"q"`
	S string `json:"s"`
}
