package context

import "time"

type TradeChanel struct {
	Timestamp int64
	Id        int64
	Price     string
	Quantity  string
	Side      string
	Symbol    string
}

type CandleCanel struct {
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Symbol    string    `json:"symbol"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Period    string    `json:"period"`
}

type RsiCanel struct {
	Symbol string  `json:"symbol"`
	Period string  `json:"period"`
	Rsi    float64 `json:"rsi"`
}
