package sbisec

import (
  "encoding/csv"
  "fmt"
  "log"
  "strconv"
  "strings"
)

type MACDSignal int
const (
  MACDBuy MACDSignal = iota
  MACDSell
)

type ScreenedStockBase struct{
  YYYYMMDD int
  Code string
  Name string
}
type ScreenedStockRSI struct{
  ScreenedStockBase
  Price float64
  RSI float64
}
type ScreenedStockMACDSignal struct{
  ScreenedStockBase
  Price float64
  MACDSignal MACDSignal
}
// type ScreenedStock struct{
//   YYYYMMDD int64
//   Code string
//   Name string
//   Price float64
//   RSI float64
//   MACDSignal MACDSignal
// }

func ReadRSI(screenedData []byte) (stocks []ScreenedStockRSI, err error) {
  log.Printf("[INFO] Read RSI data from screened data.")

  csvReader := csv.NewReader(strings.NewReader(string(screenedData)))
  csvReader.LazyQuotes = true

  // Skip first row
  _, err = csvReader.Read()
  if err != nil {
    return nil, fmt.Errorf("Error reading header:", err)
  }

  for {
    row, err := csvReader.Read()
    if err != nil { break }

    price, err := strconv.ParseFloat(row[3], 64)
    if err != nil {
      fmt.Println("Error converting price:", err)
      continue
    }
    rsi, err := strconv.ParseFloat(row[5], 64)
    if err != nil {
      fmt.Println("Error converting RSI:", err)
      continue
    }

    stock := ScreenedStockRSI{
      ScreenedStockBase: ScreenedStockBase{
        YYYYMMDD: 20250331,
        Code: row[0],
        Name: row[1],
      },
      Price: price,
      RSI: rsi,
    }
    stocks = append(stocks, stock)
  }

  return stocks, nil
}

func ReadMACDSignal(screenedData []byte, macdSignal MACDSignal) (stocks []ScreenedStockMACDSignal, err error) {
  log.Printf("[INFO] Read MACD Signal data from screened data.")

  csvReader := csv.NewReader(strings.NewReader(string(screenedData)))
  csvReader.LazyQuotes = true

  // Skip first row
  _, err = csvReader.Read()
  if err != nil {
    return nil, fmt.Errorf("Error reading header:", err)
  }

  for {
    row, err := csvReader.Read()
    if err != nil { break }

    price, err := strconv.ParseFloat(row[3], 64)
    if err != nil {
      fmt.Println("Error converting price:", err)
      continue
    }

    stock := ScreenedStockMACDSignal{
      ScreenedStockBase: ScreenedStockBase{
        YYYYMMDD: 20250331,
        Code: row[0],
        Name: row[1],
      },
      Price: price,
      MACDSignal: macdSignal,
    }
    stocks = append(stocks, stock)
  }

  return stocks, nil
}
