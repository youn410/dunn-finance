package dao

import (
  "encoding/csv"
  "io"
  "os"
  "strconv"
  "testing"

  "dunn-finance/pkg/database"
  "dunn-finance/pkg/model"
)

var TestManager = &database.DBManager{
  Driver: "sqlite3",
  DSN: ":memory:",
}

var createStocksTableSql = "CREATE TABLE stocks (code TEXT PRIMARY KEY, name TEXT)"
var createAdjustedDailyOhlcvsTableSql = `
  CREATE TABLE adjusted_daily_ohlcvs (
    yyyymmdd     INTEGER,
    code         TEXT,
    open_price   REAL,
    high_price   REAL,
    low_price    REAL,
    close_price  REAL,
    dma_price_5  REAL,
    dma_price_25 REAL,
    dma_price_75 REAL,
    vmap         REAL,
    volume       REAL,
    vma_5        REAL,
    vma_25       REAL,
    PRIMARY KEY (code, yyyymmdd),
    FOREIGN KEY (code) REFERENCES codes(code)
  );`

func PrepareTestDB(t *testing.T) database.DBConnector {
  db := TestManager.GetDBInstance()
  t.Cleanup(func() { db.Close() })

  if  _, err := db.Exec(createStocksTableSql); err != nil {
    t.Fatalf("Failed to create test stocks table: %v", err)
  }
  if  _, err := db.Exec(createAdjustedDailyOhlcvsTableSql); err != nil {
    t.Fatalf("Failed to create test adjusted_daily_ohlcvs table: %v", err)
  }

  return db
}

func LoadOhlcvCSV(path string) ([]*model.AdjustedDailyOHLCV, error) {
  f, err := os.Open(path)
  if err != nil { return nil, err }
  defer f.Close()

  r := csv.NewReader(f)
  if _, err := r.Read(); err != nil { return nil, err }

  var records []*model.AdjustedDailyOHLCV
  for {
    rec, err := r.Read()
    if err == io.EOF { break }
    if err != nil { return nil, err }

    parseFloat := func(s string) float64 {
      v, _ := strconv.ParseFloat(s, 64)
      return v
    }

    records = append(records, &model.AdjustedDailyOHLCV{
			Yyyymmdd:   rec[0],
			Code:       rec[1],
			OpenPrice:  parseFloat(rec[2]),
			HighPrice:  parseFloat(rec[3]),
			LowPrice:   parseFloat(rec[4]),
			ClosePrice: parseFloat(rec[5]),
			DMAPrice5:  parseFloat(rec[6]),
			DMAPrice25: parseFloat(rec[7]),
			DMAPrice75: parseFloat(rec[8]),
			VMAP:       parseFloat(rec[9]),
			Volume:     parseFloat(rec[10]),
			VMA5:       parseFloat(rec[11]),
			VMA25:      parseFloat(rec[12]),
    })
  }

  return records, nil
}
