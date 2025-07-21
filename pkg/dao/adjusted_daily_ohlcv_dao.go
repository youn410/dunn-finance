package dao

import (
  "database/sql"

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
    ON CONFLICT(yyyymmdd, code) DO UPDATE SET
      open_price   = excluded.open_price,
      high_price   = excluded.high_price,
      low_price    = excluded.low_price,
      close_price  = excluded.close_price,
      dma_price_5  = excluded.dma_price_5,
      dma_price_25 = excluded.dma_price_25,
      dma_price_75 = excluded.dma_price_75,
      vmap         = excluded.vmap,
      volume       = excluded.volume,
      vma_5        = excluded.vma_5,
      vma_25       = excluded.vma_25
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

  var (
    openPrice, highPrice, lowPrice, closePrice sql.NullFloat64
    dmaPrice5, dmaPrice25, dmaPrice75 sql.NullFloat64
    vmap, volume, vma5, vma25 sql.NullFloat64
  )
  var ohlcv model.AdjustedDailyOHLCV

  err := row.Scan(
    &ohlcv.Yyyymmdd,
    &ohlcv.Code,
    &openPrice,
    &highPrice,
    &lowPrice,
    &closePrice,
    &dmaPrice5,
    &dmaPrice25,
    &dmaPrice75,
    &vmap,
    &volume,
    &vma5,
    &vma25,
  )
  if err != nil { return nil, err }

  toPointer := func (n sql.NullFloat64) *float64 {
    if n.Valid { return &n.Float64 }
    return nil
  }
  ohlcv.OpenPrice = toPointer(openPrice)
  ohlcv.HighPrice = toPointer(highPrice)
  ohlcv.LowPrice = toPointer(lowPrice)
  ohlcv.ClosePrice = toPointer(closePrice)
  ohlcv.DMAPrice5 = toPointer(dmaPrice5)
  ohlcv.DMAPrice25 = toPointer(dmaPrice25)
  ohlcv.DMAPrice75 = toPointer(dmaPrice75)
  ohlcv.VMAP = toPointer(vmap)
  ohlcv.Volume = toPointer(volume)
  ohlcv.VMA5 = toPointer(vma5)
  ohlcv.VMA25 = toPointer(vma25)

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
