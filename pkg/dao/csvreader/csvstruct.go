package csvreader

import (
  "time"
)

// 日付,始値,高値,安値,終値,5日平均,25日平均,75日平均,VWAP,出来高,5日平均,25日平均
// 2024/12/30,"9,430","9,440","9,257","9,264","9,260.40","9,093.32","8,335.81","9,310.6120","2,486,400","2,689,920.00","4,661,964.00"

type ChartDataBase struct {
  Open float64
  High float64
  Low float64
  Close float64
}

type DailyChartData struct {
  ChartDataBase
  Date time.Time
}
var DailyChartCSVHeaderMapping = map[string]string{
  "日付": "Date",
  "始値": "Open",
  "高値": "High",
  "安値": "Low",
  "終値": "Close",
}
