CREATE TABLE IF NOT EXISTS daily_stocks (
  yyyymmdd INTEGER,
  code TEXT,
  openPrice REAL,
  highPrice REAL,
  lowPrice REAL,
  closePrice REAL,
  rsi REAL,
  macd REAL,
  signal REAL
)
