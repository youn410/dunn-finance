package dao_test

import (
  "reflect"
  "testing"

  _ "github.com/mattn/go-sqlite3"

  "dunn-finance/pkg/dao"
  "dunn-finance/pkg/model"
)

var createAdjustedDailyOhlcvsTableSQL = `
CREATE TABLE adjusted_daily_ohlcvs (
  yyyymmdd     TEXT,
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
  PRIMARY KEY (yyyymmdd, code),
  FOREIGN KEY (code) REFERENCES codes(code)
);`

func NewAdjustedDailyOHLCV(overrides ...func(*model.AdjustedDailyOHLCV)) *model.AdjustedDailyOHLCV {
  o := &model.AdjustedDailyOHLCV{
    Yyyymmdd:   "20250706",
    Code:       "1234",
    OpenPrice:  1000.0,
    HighPrice:  1050.0,
    LowPrice:   990.0,
    ClosePrice: 1020.0,
    DMAPrice5:  1010.0,
    DMAPrice25: 1005.0,
    DMAPrice75: 995.0,
    VMAP:       1015.0,
    Volume:     150000.0,
    VMA5:       140000.0,
    VMA25:      130000.0,
  }

  for _, fn := range overrides {
    fn(o)
  }

  return o
}

func TestAdjustedDaliyOhlcvDao_Create_Success(t *testing.T) {
  db := dao.PrepareTestDB(t)

  stockDao := dao.StockDAO{DB: db}
  stock := &model.Stock{Code: "1234", Name: "テスト会社"}
  err := stockDao.Create(stock)
  if err != nil { t.Errorf("Failed to create stock record: %v", err) }

  ohlcvDao := dao.AdjustedDailyOHLCVDAO{DB: db}

  ohlcv := NewAdjustedDailyOHLCV()
  err = ohlcvDao.Create(ohlcv)
  if err != nil { t.Errorf("Failed to create adjusted daily ohlcv record: %v", err) }
}

func TestAdjustedDaliyOhlcvDao_Find_Success(t *testing.T) {
  db := dao.PrepareTestDB(t)

  stockDao := dao.StockDAO{DB: db}
  stock := &model.Stock{Code: "1234", Name: "テスト会社"}
  err := stockDao.Create(stock)
  if err != nil { t.Errorf("Failed to create stock record: %v", err) }

  ohlcvDao := dao.AdjustedDailyOHLCVDAO{DB: db}
  actualOhlcv := NewAdjustedDailyOHLCV()
  err = ohlcvDao.Create(actualOhlcv)
  if err != nil { t.Errorf("Failed to create adjusted daily ohlcv record: %v", err) }

  expectedOhlcv, _ := ohlcvDao.Find("1234", "20250706")

  vExpected := reflect.ValueOf(expectedOhlcv).Elem()
  vActual := reflect.ValueOf(actualOhlcv).Elem()
  ohlcvType := vExpected.Type()
  for i := 0; i < vExpected.NumField(); i++ {
    fieldName := ohlcvType.Field(i).Name
    valExpected := vExpected.Field(i).Interface()
    valActual := vActual.Field(i).Interface()
    if !reflect.DeepEqual(valExpected, valActual) {
      t.Errorf("Field: %s, expected: %v, but got %v", fieldName, valExpected, valActual)
    }
  }
}

func TestAdjustedDailyOhlcvDao_FindByDateRange_Success(t *testing.T) {
  db := dao.PrepareTestDB(t)
  ohlcvDao := dao.AdjustedDailyOHLCVDAO{DB: db}

  stockDao := dao.StockDAO{DB: db}
  stock := &model.Stock{Code: "1234", Name: "テスト会社"}
  err := stockDao.Create(stock)
  if err != nil { t.Errorf("Failed to create stock record: %v", err) }

  path := "testdata/adjusted_daily_ohlcvs.csv"
  records, err := dao.LoadOhlcvCSV(path)
  if err != nil { t.Fatalf("Load CSV: %v", err) }
  for _, rec := range records {
    if err := ohlcvDao.Create(rec); err != nil {
      t.Fatalf("Insert AdjustedDailyOHLCV: %v", err)
    }
  }

  got, err := ohlcvDao.FindByDateRange("1234", "20250702", "20250704")
  if err != nil { t.Fatalf("FindByDateRange: %v", err) }
  if len(got) != 3 { t.Errorf("want 3 records, got %d", len(got)) }
}
