package main

import (
  "flag"
  "log"

  _ "github.com/mattn/go-sqlite3"

  "dunn-finance/pkg/csvreader"
  "dunn-finance/pkg/dao"
  "dunn-finance/pkg/database"
)

var fieldMap = map[int]string{
  0: "Yyyymmdd",
  1: "OpenPrice",
  2: "HighPrice",
  3: "LowPrice",
  4: "ClosePrice",
  5: "DMAPrice5",
  6: "DMAPrice25",
  7: "DMAPrice75",
  8: "VMAP",
  9: "Volume",
  10: "VMA5",
  11: "VMA25",
}

func main() {
  log.Println("[INFO] update adjusted daily ohlcv starts.")

  code := flag.String("code", "", "stock code")
  csvPath := flag.String("csvpath", "", "Path to the CSV file")
  dbPath := flag.String("dbpath", "", "Path to the DB file")
  isSkipHeader := flag.Bool("skip-header", true, "Whether to skip the header row (default: true)")
  offset := flag.Int("offset", 0, "Number of rows to skip from the beginning")
  limit := flag.Int("limit", 100, "Maximum number of rows to read")

  flag.Parse()

  if *code == "" { log.Fatal("[ERROR] Please specify the stock code -code") }
  if *csvPath == "" { log.Fatal("[ERROR] Please specify the path to CSV file using -csvpath") }
  if *dbPath == "" { log.Fatal("[ERROR] Please specify the path to DB file using -dbpath") }

  log.Printf("[INFO] code: %s, CSV path: %s, skip header: %t, offset: %d, limit: %d\n", *code, *csvPath, *isSkipHeader, *offset, *limit)

  dbManager := &database.DBManager{ Driver: "sqlite3", DSN: *dbPath }
  db := dbManager.GetDBInstance()
  defer db.Close()

  ohlcvDao := dao.AdjustedDailyOHLCVDAO{DB: db}
  for {
    records, err := csvreader.LoadAdjustedDailyOHLCVsFromCSV(*code, *csvPath, fieldMap, *isSkipHeader, *offset, *limit)
    if err != nil { log.Fatalf("Failed to load CSV: %v", err) }
    if len(records) == 0 {
      log.Println("[INFO] Reached end of CSV")
      break
    }

    log.Printf("[INFO] Loaded %d records from offset %d\n", len(records), *offset)

    for _, record := range records {
      err := ohlcvDao.Create(record)
      if err != nil { log.Printf("Failed to insert: %+v, err: %v", record, err) }
    }

    *offset += len(records)
    *isSkipHeader = false
  }

  log.Println("[INFO] update adjusted daily ohlcv ends.")
}
