package dao_test

import (
  "testing"

  _ "github.com/mattn/go-sqlite3"

  // "dunn-finance/pkg/database"
  "dunn-finance/pkg/dao"
  "dunn-finance/pkg/model"
)

// var manager = &database.DBManager{ Driver: "sqlite3", DSN: ":memory:", }

func TestStockDao_Create_Success(t *testing.T) {
  db := dao.PrepareTestDB(t)
  // db := dao.TestManager.GetDBInstance()
  // defer db.Close()
  // db.Exec(createStocksSQL)

  stockDao := dao.StockDAO{DB: db}
  stock := &model.Stock{Code: "1234", Name: "テスト会社"}

  err := stockDao.Create(stock)
  if err != nil { t.Errorf("Failed to create stock record: %v", err) }
}

func TestStockDao_Find_Success(t *testing.T) {
  db := dao.PrepareTestDB(t)
  // db := dao.TestManager.GetDBInstance()
  // db := manager.GetDBInstance()
  // defer db.Close()
  // db.Exec(createStocksSQL)
  stockDao := dao.StockDAO{DB: db}

  code := "1234"
  name := "テスト会社"
  stock := &model.Stock{Code: code, Name: name}
  err := stockDao.Create(stock)
  if err != nil { t.Fatalf("Failed to create stock record: %v", err)}

  actualStock, _ := stockDao.Find(code)
  if code != actualStock.Code { t.Errorf("got %s, want %s", actualStock.Code, code) }
  if name != actualStock.Name { t.Errorf("got %s, want %s", actualStock.Name, name) }
}
