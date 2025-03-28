package dbmanager

import (
  "database/sql"
  "fmt"
  "log"
  "sync"

  _ "github.com/mattn/go-sqlite3"

  "dunn-finance/pkg/browser/sites/sbisec"
)

type DBManager struct {
  db *sql.DB
}

var instance *DBManager
var once sync.Once

func GetDBManager(dbFilePath string) *DBManager {
  once.Do(func() {
    db, err := sql.Open("sqlite3", dbFilePath)
    if err != nil {
      log.Fatalf("Failed to open database: %v", err)
    }
    instance = &DBManager{db: db}
  })

  return instance
}

func (m *DBManager) SelectStocksBuyMACDSignal(signal sbisec.MACDSignal, yyyymmdd int) (codes []string, err error) {
  var buySignal = 0
  var sellSignal = 0
  if signal == sbisec.MACDBuy {
    buySignal = 1
  }
  if signal == sbisec.MACDSell {
    sellSignal = 1
  }
  selectSQL := fmt.Sprintf(`
    SELECT code FROM macd_signal
    WHERE buy = %d AND sell = %d AND yyyymmdd = %d
  `, buySignal, sellSignal, yyyymmdd)
  log.Println(fmt.Sprintf("[DEBUG] SQL: %s", selectSQL))

  rows, err := m.db.Query(selectSQL)
  if err != nil {
    return nil, fmt.Errorf("SQL Error '%s': %w", selectSQL, err)
  }
  defer rows.Close()

  for rows.Next() {
    var code string
    err := rows.Scan(&code)
    if err != nil {
      fmt.Println("Scan Error:", err)
      continue
    }

    // fmt.Printf("Code: %s\n", code)
    codes = append(codes, code)
  }

  return codes, nil
}

func (m *DBManager) InsertScreenedStocks(stocks []sbisec.ScreenedStockRSI) error {
  insertSql := `
    INSERT INTO rsi (yyyymmdd, code, name, price, rsi)
    VALUES (?, ?, ?, ?, ?)
  `
  stmt, err := m.db.Prepare(insertSql)
  if err != nil {
    return fmt.Errorf("Error preparing SQL: %w", err)
  }
  defer stmt.Close()

  for _, stock := range stocks {
    _, err := stmt.Exec(stock.YYYYMMDD, stock.Code, stock.Name, stock.Price, stock.RSI)
    if err != nil {
      fmt.Println("Error inserting screened stock:", err)
      continue
    }
  }

  return nil
}

func (m *DBManager) InsertScreenedStocksMACDSignal(stocks []sbisec.ScreenedStockMACDSignal) error {
  insertSql := `
    INSERT INTO macd_signal(yyyymmdd, code, name, price, buy, sell)
    VALUES (?, ?, ?, ?, ?, ?)
  `
  stmt, err := m.db.Prepare(insertSql)
  if err != nil {
    return fmt.Errorf("Error preparing SQL: %w", err)
  }
  defer stmt.Close()

  for _, stock := range stocks {
    var sellSignal = 0
    var buySignal = 0
    if stock.MACDSignal == sbisec.MACDBuy {
      buySignal = 1
    }
    if stock.MACDSignal == sbisec.MACDSell {
      sellSignal = 1
    }

    _, err := stmt.Exec(stock.YYYYMMDD, stock.Code, stock.Name, stock.Price, buySignal, sellSignal)
    if err != nil {
      fmt.Println("Error inserting screened stock:", err)
      continue
    }
  }

  return nil
}

func (m *DBManager) Close() {
  if err := m.db.Close(); err != nil {
    log.Printf("Failed to close database: %v", err)
  }
}
