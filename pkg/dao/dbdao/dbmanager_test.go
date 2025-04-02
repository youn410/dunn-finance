package dbdao_test

import (
  "fmt"
  "errors"
  "database/sql"
  "os"
  "path/filepath"
  "strings"
  "testing"

  _ "github.com/mattn/go-sqlite3"

  "dunn-finance/pkg/dao/dbdao"
)

type mockConnector struct {
  shouldPingFail bool
}

func (m *mockConnector) Ping() error {
  if m.shouldPingFail { return errors.New("mock Ping error") }
  return nil
}

func (m *mockConnector) Close() error {
  return nil
}

func (m *mockConnector) GetRawDB() *sql.DB {
  return nil
}

func (m *mockConnector) Exec(query string, args ...any) (sql.Result, error) {
  return nil, nil
}

func (m *mockConnector) QueryRow(query string, args ...any) *sql.Row {
  return nil
}

func TestSQLite3_CreateDBInstance_Success(t *testing.T) {
  manager := &dbdao.DBManager{
    Driver: "sqlite3",
    DSN:    ":memory:",
  }

  db1 := manager.GetDBInstance()
  defer db1.Close()
  db2 := manager.GetDBInstance()

  if db1 != db2 {
    t.Errorf("Expected the same connection. But got the different connection.")
  }
}

func TestSQLite3_Reconnect_Success(t *testing.T) {
  manager := &dbdao.DBManager{
    Driver: "sqlite3",
    DSN:    ":memory:",
  }

  db1 := manager.GetDBInstance()
  err := db1.Close()
  if err != nil { t.Fatalf("Failed to close the first connection: %v", err) }

  db2 := manager.GetDBInstance()
  defer db2.Close()

  if db1 == db2 {
    t.Errorf("Expected a new connection after close. But got the same connection")
  }
  if err := db2.Ping(); err != nil {
    t.Errorf("Ping failed after reconnecting: %v", err)
  }
}

func Util_ExecWithoutArgs(db dbdao.DBConnector, sqlFilePath string) string {
  sqlBytes, _ := os.ReadFile(sqlFilePath)
  sqlString := string(sqlBytes)
  db.Exec(sqlString)

  return sqlString
}
func Util_GetDBInstance() dbdao.DBConnector {
  manager := &dbdao.DBManager{
    Driver: "sqlite3",
    DSN:    ":memory:",
  }
  db := manager.GetDBInstance()

  return db
}

func TestSQLite3_CreateTables_Success(t *testing.T) {
  manager := &dbdao.DBManager{
    Driver: "sqlite3",
    DSN:    ":memory:",
  }
  db := manager.GetDBInstance()
  defer db.Close()

  var sqlFilePath string
  var sqlString string
  var schema string
  var expectedSql string
  var actualSql string

  // codes table
  sqlFilePath = filepath.Join("..", "..", "..", "configs", "sql", "sqlite3", "codes.sql")
  sqlString = Util_ExecWithoutArgs(db, sqlFilePath)

  db.QueryRow(`SELECT sql FROM sqlite_master WHERE type='table' AND name='codes'`).Scan(&schema)
  expectedSql = strings.ReplaceAll(strings.ReplaceAll(sqlString, "IF NOT EXISTS ", ""), "\n", "")
  actualSql = strings.ReplaceAll(schema, "\n", "")
  if expectedSql != actualSql {
    t.Errorf("Failed to create codes table. got: %s, want: %s.", schema, sqlString)
  }

  // daily_stocks table
  sqlFilePath = filepath.Join("..", "..", "..", "configs", "sql", "sqlite3", "daily_stocks.sql")
  sqlString = Util_ExecWithoutArgs(db, sqlFilePath)

  db.QueryRow(`SELECT sql FROM sqlite_master WHERE type='table' AND name='daily_stocks'`).Scan(&schema)
  expectedSql = strings.ReplaceAll(strings.ReplaceAll(sqlString, "IF NOT EXISTS ", ""), "\n", "")
  actualSql = strings.ReplaceAll(schema, "\n", "")
  if expectedSql != actualSql {
    t.Errorf("Failed to create daily_stocks table. got: %s, want: %s.", schema, sqlString)
  }
}

func TestSQLite3_InsertOneCode_Success(t *testing.T) {
  db := Util_GetDBInstance()
  defer db.Close()

  var sqlFilePath string

  // codes table
  sqlFilePath = filepath.Join("..", "..", "..", "configs", "sql", "sqlite3", "codes.sql")
  Util_ExecWithoutArgs(db, sqlFilePath)

  testCode := "1234"
  testName := "某田株式会社"
  insertSql := "INSERT INTO codes (code, name) VALUES (?, ?)"
  db.Exec(insertSql, testCode, testName)

  // Assertion
  var actualCode string
  var actualName string
  selectSql := "SELECT code, name FROM codes WHERE code = ?"
  row := db.QueryRow(selectSql, testCode)
  row.Scan(&actualCode, &actualName)

  if testCode != actualCode { t.Errorf("Unexpected code. want: %s, got: %s", testCode, actualCode) }
  if testName != actualName { t.Errorf("Unexpected name. want: %s, got: %s", testName, actualName) }
}

func TestSQLite3_InsertOneDailyStock_Success(t *testing.T) {
  // Prepare
  db := Util_GetDBInstance()
  defer db.Close()

  var sqlFilePath string
  sqlFilePath = filepath.Join("..", "..", "..", "configs", "sql", "sqlite3", "codes.sql")
  Util_ExecWithoutArgs(db, sqlFilePath)
  sqlFilePath = filepath.Join("..", "..", "..", "configs", "sql", "sqlite3", "daily_stocks.sql")
  Util_ExecWithoutArgs(db, sqlFilePath)

  var insertSql string
  testCode := "1234"
  testName := "某田株式会社"
  insertSql = "INSERT INTO codes (code, name) VALUES (?, ?)"
  db.Exec(insertSql, testCode, testName)

  testYyyymmdd := 20250329
  testOpen := 100.5
  // testHigh := 110
  // testLow := 93
  // testClose := 100
  // testRsi := 30.5
  // testMacd := 40
  // testSignal := 50
  // insertSql = `INSERT INTO daily_stocks (
// yyyymmdd, code, open, high, low, close,
// rsi, macd, signal)
// VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
  // db.Exec(insertSql, testYyyymmdd, testCode, testOpen, testHigh, testLow, testClose, testRsi, testMacd, testSignal)
  insertSql = `INSERT INTO daily_stocks (yyyymmdd, code, openPrice) VALUES(?, ?, ?)`
  db.Exec(insertSql, testYyyymmdd, testCode, testOpen)

  // Assertion
  var actualYyyymmdd int
  var actualCode string
  var actualOpen float64
  // var actualHigh float64
  // var actualLow float64
  // var actualClose float64
  // var actualRsi float64
  // var actualMacd float64
  // var actualSignal float64
  selectSql := `SELECT yyyymmdd, code, openPrice FROM daily_stocks WHERE code = ?`
  row := db.QueryRow(selectSql, testCode)
   // row.Scan(&actualYyyymmdd, &actualCode, &actualOpen, &actualHigh, &actualLow, &actualClose, &actualRsi, &actualMacd, &actualSignal)
   row.Scan(&actualYyyymmdd, &actualCode, &actualOpen)
  // row := db.QueryRow(selectSql)
  // selectSql := `SELECT yyyymmdd FROM daily_stocks WHERE code = ?`
  // var count int
  // row := db.QueryRow(selectSql, testCode)
  // row.Scan(&count)
  // t.Logf("%d", count)
  if testYyyymmdd != actualYyyymmdd { t.Errorf("Unexpected yyyymmdd. want: %d, got: %d", testYyyymmdd, actualYyyymmdd) }
  if testCode != actualCode { t.Errorf("Unexpected code. want: %s, got: %s", testCode, actualCode) }
  fmt.Println(actualYyyymmdd)
  if float64(testOpen) != float64(actualOpen) { t.Errorf("Unexpected open. want: %f, got: %f", float64(testOpen), float64(actualOpen)) }
  // if float64(testHigh) != float64(actualHigh) { t.Errorf("Unexpected high. want: %f, got: %f", float64(testHigh), float64(actualHigh)) }
  // if float64(testLow) != float64(actualLow) { t.Errorf("Unexpected low. want: %f, got: %f", float64(testLow), float64(actualLow)) }
  // if float64(testClose) != float64(actualClose) { t.Errorf("Unexpected close. want: %f, got: %f", float64(testClose), float64(actualClose)) }
  // if float64(testRsi) != float64(actualRsi) { t.Errorf("Unexpected rsi. want: %f, got: %f", float64(testRsi), float64(actualRsi)) }
  // if float64(testMacd) != float64(actualMacd) { t.Errorf("Unexpected macd. want: %f, got: %f", float64(testMacd), float64(actualMacd)) }
  // if float64(testSignal) != float64(actualSignal) { t.Errorf("Unexpected signal. want: %f, got: %f", float64(testSignal), float64(actualSignal)) }
}

/*
func TestNewDBManagerWithConnector_Success(t *testing.T) {
  mockConn := &mockConnector{shouldPingFail: false}
  dbManager, err := dbdao.NewDBManagerWithConnector("mockdb", "dsn", mockConn)

  if err != nil { t.Fatalf("Unexpected error: %v", err) }
  if dbManager.DB != mockConn { t.Errorf("Expected mockConn does not used.") }
}

func TestNewDBManagerWithConnector_Failure(t *testing.T) {
  mockConn := &mockConnector{shouldPingFail: true}
  _, err := dbdao.NewDBManagerWithConnector("mockdb", "dsn", mockConn)

  if err == nil || err.Error() != "mock Ping error" {
    t.Fatalf("Ping error was not handled as expected: %v", err)
  }
}

func TestSQLite3Connection_Success(t *testing.T) {
  tmpFile, err := os.CreateTemp("testdata", "tempdb-*.sqlite")
  if err != nil { t.Fatalf("Failed to Create temporary file: %v", err) }
  defer os.Remove(tmpFile.Name())
  // t.Logf("✅ Creating temporary file succeeded")

  dsn := "file:" + tmpFile.Name() + "?cache=shared&mode=rwc"
  dbManager, err := dbdao.NewDBManager("sqlite3", dsn)
  defer dbManager.DB.Close()

  if err != nil { t.Fatalf("Unexpected error: %v", err) }
  if err := dbManager.DB.Ping(); err != nil {
    t.Fatalf("Ping failed: %v", err)
  }
}

func TestSQLite3CreateTables_Success(t *testing.T) {
  tmpFile, err := os.CreateTemp("testdata", "tempdb-*.sqlite")
  if err != nil { t.Fatalf("Failed to Create temporary file: %v", err) }
  defer os.Remove(tmpFile.Name())

  dsn := "file:" + tmpFile.Name() + "?cache=shared&mode=rwc&_foreign_keys=true"
  dbManager, _ := dbdao.NewDBManager("sqlite3", dsn)
  defer dbManager.DB.Close()

  createCodesSqlFilePath := filepath.Join("..", "..", "..", "configs", "sql", "sqlite3", "codes.sql")
  sqlBytes, err := os.ReadFile(createCodesSqlFilePath)
  if err != nil { t.Fatalf("Failed to Read file '%s': %v", createCodesSqlFilePath, err) }
  createCodesSqlText := string(sqlBytes)
  _, err = dbManager.DB.Exec(createCodesSqlText)
  if err != nil { t.Fatalf("Failed to exec sql: %v", err) }

  var schema string
  dbManager.DB.QueryRow(`SELECT sql FROM sqlite_master WHERE type='table' AND name='codes'`).Scan(&schema)
  if strings.ReplaceAll(strings.ReplaceAll(createCodesSqlText, "IF NOT EXISTS ", ""), "\n", "") != strings.ReplaceAll(schema, "\n", "") {
    t.Errorf("Unexpected codes table: %s", strings.ReplaceAll(schema, "\n", ""))
  }

  createDailyStocksSqlFilePath := filepath.Join("..", "..", "..", "configs", "sql", "sqlite3", "daily_stocks.sql")
  sqlBytes, err = os.ReadFile(createDailyStocksSqlFilePath)
  if err != nil { t.Fatalf("Failed to Read file '%s': %v", createDailyStocksSqlFilePath, err) }
  createDailyStocksSqlText := string(sqlBytes)
  _, err = dbManager.DB.Exec(createDailyStocksSqlText)
  if err != nil { t.Fatalf("Failed to exec sql: %v", err) }

  dbManager.DB.QueryRow(`SELECT sql FROM sqlite_master WHERE type='table' AND name='daily_stocks'`).Scan(&schema)
  if strings.ReplaceAll(strings.ReplaceAll(createDailyStocksSqlText, "IF NOT EXISTS ", ""), "\n", "") != strings.ReplaceAll(schema, "\n", "") {
    t.Errorf("Unexpected codes table: %s", strings.ReplaceAll(schema, "\n", ""))
  }
}

func TestSQLite3InsertCode_Success(t *testing.T) {
  tmpFile, err := os.CreateTemp("testdata", "tempdb-*.sqlite")
  if err != nil { t.Fatalf("Failed to Create temporary file: %v", err) }
  defer os.Remove(tmpFile.Name())

  dsn := "file:" + tmpFile.Name() + "?cache=shared&mode=rwc&_foreign_keys=true"
  dbManager, _ := dbdao.NewDBManager("sqlite3", dsn)
  defer dbManager.DB.Close()

  createCodesSqlFilePath := filepath.Join("..", "..", "..", "configs", "sql", "sqlite3", "codes.sql")
  sqlBytes, err := os.ReadFile(createCodesSqlFilePath)
  createCodesSqlText := string(sqlBytes)
  dbManager.DB.Exec(createCodesSqlText)

  testCode := "1234"
  testName := "某田株式会社"
  insertSql := "INSERT INTO codes (code, name) VALUES (?, ?)"
  dbManager.DB.Exec(insertSql, testCode, testName)

  // Assertion
  var actualCode string
  var actualName string
  selectSql := "SELECT code, name FROM codes WHERE code = ?"
  row := dbManager.DB.QueryRow(selectSql, testCode)
  row.Scan(&actualCode, &actualName)

  if testCode != actualCode { t.Errorf("Unexpected code. want: %s, got: %s", testCode, actualCode) }
  if testName != actualName { t.Errorf("Unexpected name. want: %s, got: %s", testName, actualName) }
}

func TestSQLite3InsertDailyStocks_Success(t *testing.T) {
  // Prepare
  tmpFile, err := os.CreateTemp("testdata", "tempdb-*.sqlite")
  if err != nil { t.Fatalf("Failed to Create temporary file: %v", err) }
  defer os.Remove(tmpFile.Name())

  dsn := "file:" + tmpFile.Name() + "?cache=shared&mode=rwc&_foreign_keys=true"
  dbManager, _ := dbdao.NewDBManager("sqlite3", dsn)
  defer dbManager.DB.Close()

  createCodesSqlFilePath := filepath.Join("..", "..", "..", "configs", "sql", "sqlite3", "codes.sql")
  sqlBytes, err := os.ReadFile(createCodesSqlFilePath)
  createCodesSqlText := string(sqlBytes)
  dbManager.DB.Exec(createCodesSqlText)

  testCode := "1234"
  testName := "某田株式会社"
  insertSql := "INSERT INTO codes (code, name) VALUES (?, ?)"
  dbManager.DB.Exec(insertSql, testCode, testName)

  // Test
  testYyyymmdd := 20250329
  testOpen := 100.5
  testHigh := 110
  testLow := 93
  testClose := 100
  testRsi := 30.5
  testMacd := 40
  testSignal := 50
  insertSql = `INSERT INTO daily_codes (
yyyymmdd, code, open, high, low, close,
rsi, macd, signal)
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
  dbManager.DB.Exec(insertSql, testYyyymmdd, testCode, testOpen, testHigh, testLow, testClose, testRsi, testMacd, testSignal)

  // Assertion
  var actualYyyymmdd string
  var actualOpen float64
  // selectSql := `SELECT yyyymmdd, open FROM daily_stocks WHERE code = ?`
  selectSql := `SELECT yyyymmdd, open FROM daily_stocks`
  row := dbManager.DB.QueryRow(selectSql, testCode)
  row.Scan(&actualYyyymmdd, &actualOpen)
  t.Logf(actualYyyymmdd)
}
*/
