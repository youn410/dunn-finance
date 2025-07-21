package csvreader

import (
  "encoding/csv"
  "os"
  "reflect"
  "strconv"
  "strings"
  "time"

  "dunn-finance/pkg/model"
)


func stringToFloat(s string) *float64 {
  if s == "--" { return nil }
  s = strings.ReplaceAll(s, ",", "")
  f, _ := strconv.ParseFloat(s, 64)

  return &f
}


func parseDateToYyyymmdd(s string) string {
  t, _ := time.Parse("2025/07/13", s)
  return t.Format("20060102")
}

func LoadAdjustedDailyOHLCVsFromCSV(
  code string,
  path string,
  fieldMap map[int]string,
  isSkipHeader bool,
  offset int,
  limit int,
) ([]*model.AdjustedDailyOHLCV, error) {
  f, err := os.Open(path)
  if err != nil { return nil ,err }
  defer f.Close()

  r := csv.NewReader(f)
  if isSkipHeader {
    _, _ = r.Read()
  }

  var result []*model.AdjustedDailyOHLCV
  rowIndex := 0
  readCount := 0
  for {
    csvRow, err := r.Read()
    if err != nil { break }

    if rowIndex < offset {
      rowIndex++
      continue
    }
    if limit > 0 && readCount >= limit { break }


    ohlcv := &model.AdjustedDailyOHLCV{Code: code}
    ohlcvVal := reflect.ValueOf(ohlcv).Elem()

    for i, csvVal := range csvRow {
      fieldName := fieldMap[i]
      ohlcvField := ohlcvVal.FieldByName(fieldName)
      if !ohlcvField.IsValid() || !ohlcvField.CanSet() { continue }

      if fieldName == "Yyyymmdd" {
        t, err := time.Parse("2006/01/02", csvVal)
        if err == nil { ohlcvField.SetString(t.Format("20060102")) }

        continue
      }

      switch ohlcvField.Kind() {
        case reflect.String:
          ohlcvField.SetString(csvVal)
        case reflect.Ptr:
          if ohlcvField.Type().Elem().Kind() == reflect.Float64 {
            ohlcvField.Set(reflect.ValueOf(stringToFloat(csvVal)))
          }
      }
    }

    result = append(result, ohlcv)
    rowIndex++
    readCount++
  }

  return result, nil
}
