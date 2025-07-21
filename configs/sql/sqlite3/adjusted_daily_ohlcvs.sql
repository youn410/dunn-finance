CREATE TABLE IF NOT EXISTS adjusted_daily_ohlcvs (
  yyyymmdd     TEXT NOT NULL,
  code         TEXT NOT NULL,
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
);
