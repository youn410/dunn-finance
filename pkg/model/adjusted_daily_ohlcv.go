package model

type AdjustedDailyOHLCV struct {

  Yyyymmdd   string
  Code       string
  OpenPrice  float64
  HighPrice  float64
  LowPrice   float64
  ClosePrice float64
  DMAPrice5  float64
  DMAPrice25 float64
  DMAPrice75 float64
  VMAP       float64
  Volume     float64
  VMA5       float64
  VMA25      float64
}
