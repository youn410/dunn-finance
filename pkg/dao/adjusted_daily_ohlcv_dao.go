package dao

import (
  "dunn-finance/pkg/model"
  "dunn-finance/pkg/database"
)

type AdjustedDailyOHLCVDAO struct {
  DB database.DBConnector
}

func (dao *AdjustedDailyOHLCVDAO) Create(ohlcv *model.AdjustedDailyOHLCV) error {
  _, err := dao.DB.Exec(
    `
    INSERT INTO adjusted_daily_ohlcvs (
      yyyymmdd,
      code,
      open_price,
      high_price,
      low_price,
      close_price,
      dma_price_5,
      dma_price_25,
      dma_price_75,
      vmap,
      volume,
      vma_5,
      vma_25
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `,
    ohlcv.Yyyymmdd,
    ohlcv.Code,
    ohlcv.OpenPrice,
    ohlcv.HighPrice,
    ohlcv.LowPrice,
    ohlcv.ClosePrice,
    ohlcv.DMAPrice5,
    ohlcv.DMAPrice25,
    ohlcv.DMAPrice75,
    ohlcv.VMAP,
    ohlcv.Volume,
    ohlcv.VMA5,
    ohlcv.VMA25,
  )

  return err
}

func (dao *AdjustedDailyOHLCVDAO) Find(code string, yyyymmdd string) (*model.AdjustedDailyOHLCV, error) {
  row := dao.DB.QueryRow(`
    SELECT
      yyyymmdd,
      code,
      open_price,
      high_price,
      low_price,
      close_price,
      dma_price_5,
      dma_price_25,
      dma_price_75,
      vmap,
      volume,
      vma_5,
      vma_25
    FROM adjusted_daily_ohlcvs
    WHERE code = ? AND yyyymmdd = ?
  `, code, yyyymmdd)

  var ohlcv model.AdjustedDailyOHLCV
  err := row.Scan(
    &ohlcv.Yyyymmdd,
    &ohlcv.Code,
    &ohlcv.OpenPrice,
    &ohlcv.HighPrice,
    &ohlcv.LowPrice,
    &ohlcv.ClosePrice,
    &ohlcv.DMAPrice5,
    &ohlcv.DMAPrice25,
    &ohlcv.DMAPrice75,
    &ohlcv.VMAP,
    &ohlcv.Volume,
    &ohlcv.VMA5,
    &ohlcv.VMA25,
  )
  if err != nil { return nil, err }

  return &ohlcv, nil
}

func (dao *AdjustedDailyOHLCVDAO) FindByDateRange(code string, fromYyyymmdd string, toYyyymmdd string) ([]*model.AdjustedDailyOHLCV, error) {
  rows, err := dao.DB.Query(`
    SELECT
      yyyymmdd,
      code,
      open_price,
      high_price,
      low_price,
      close_price,
      dma_price_5,
      dma_price_25,
      dma_price_75,
      vmap,
      volume,
      vma_5,
      vma_25
    FROM adjusted_daily_ohlcvs
    WHERE code = ? AND yyyymmdd BETWEEN ? AND ?
    ORDER BY yyyymmdd
  `, code, fromYyyymmdd, toYyyymmdd)
  if err != nil { return nil, err }
  defer rows.Close()

  var results []*model.AdjustedDailyOHLCV
  for rows.Next() {
    var ohlcv model.AdjustedDailyOHLCV
    err := rows.Scan(
      &ohlcv.Yyyymmdd,
      &ohlcv.Code,
      &ohlcv.OpenPrice,
      &ohlcv.HighPrice,
      &ohlcv.LowPrice,
      &ohlcv.ClosePrice,
      &ohlcv.DMAPrice5,
      &ohlcv.DMAPrice25,
      &ohlcv.DMAPrice75,
      &ohlcv.VMAP,
      &ohlcv.Volume,
      &ohlcv.VMA5,
      &ohlcv.VMA25,
    )
    if err != nil { return nil, err }
    results = append(results, &ohlcv)
  }

  return results, nil
}
