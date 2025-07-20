package dto

type AdjustedDailyOHLCV struct {
  Yyyymmdd   string  `json:"yyyymmdd"`
  Code       string  `json:"code"`
  OpenPrice  float64 `json:"open_price"`
  HighPrice  float64 `json:"high_price"`
  LowPrice   float64 `json:"low_price"`
  ClosePrice float64 `json:"close_price"`
  DMAPrice5  float64 `json:"dma_price_5"`
  DMAPrice25 float64 `json:"dma_price_25"`
  DMAPrice75 float64 `json:"dma_price_75"`
  VMAP       float64 `json:"vmap"`
  Volume     float64 `json:"volume"`
  VMA5       float64 `json:"vma_5"`
  VMA25      float64 `json:"vma_25"`
}
