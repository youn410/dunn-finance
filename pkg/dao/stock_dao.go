package dao

import (
  "dunn-finance/pkg/model"
  "dunn-finance/pkg/database"
)

type StockDAO struct {
  DB database.DBConnector
}

func (dao *StockDAO) Create(stock *model.Stock) error {
  _, err := dao.DB.Exec(
    "INSERT INTO stocks (code, name) VALUES (?, ?)",
    stock.Code, stock.Name,
  )

  return err
}

func (dao *StockDAO) Find(code string) (*model.Stock, error) {
  row := dao.DB.QueryRow("SELECT code, name FROM stocks WHERE code = ?", code)
  var stock model.Stock
  if err := row.Scan(&stock.Code, &stock.Name); err != nil {
    return nil, err
  }

  return &stock, nil
}
