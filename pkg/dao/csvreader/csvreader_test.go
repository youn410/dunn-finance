package  csvreader

import (
  "bytes"
  "fmt"
  "testing"
  "strings"
  "time"
  "log"
  "io/ioutil"
  "os"
  "reflect"
)

var _ = fmt.Println

type ParsedDataStruct struct {
  IntField int
  Float64Field float64
  DateField time.Time
  StringField string
}

var structField2CSVHeaderIndexMapping = map[string]int {
  "IntField": 0,
  "Float64Field": 1,
  "DateField": 2,
  "StringField": 3,
}

func createCSVHeaders() []string {
  return []string {
    "整数",
    "小数",
    "日付",
    "文字列",
  }
}

func createStructField2CsvHeaderMapping() map[string]string {
  return map[string]string {
    "IntField": "整数",
    "Float64Field": "小数",
    "DateField": "日付",
    "StringField": "文字列",
  }
}

func TestMain(m *testing.M) {
  log.SetOutput(ioutil.Discard)

  exitCode := m.Run()

  os.Exit(exitCode)
}

func TestGetStructField2CSVHeaderIndexMappingSuccess(t *testing.T) {
  structField2CsvHeaderMapping := createStructField2CsvHeaderMapping()
  csvHeaders := createCSVHeaders()
  var expected = map[string]int {
    "IntField": 0,
    "DateField": 2,
    "Float64Field": 1,
    "StringField": 3,
  }

  actual, _ := GetStructField2CSVHeaderIndexMapping[ParsedDataStruct](csvHeaders, structField2CsvHeaderMapping)
  if !reflect.DeepEqual(actual, expected) {
    t.Errorf("got: %v, want: %v", actual, expected)
  }
}

func TestGetStructField2CSVHeaderIndexMappingFailure(t *testing.T) {
  t.Run("structField2CsvHeaderMapping has missing field", func(t *testing.T){
    structField2CsvHeaderMapping := createStructField2CsvHeaderMapping()
    csvHeaders := createCSVHeaders()
    delete(structField2CsvHeaderMapping, "StringField")

    _, err := GetStructField2CSVHeaderIndexMapping[ParsedDataStruct](csvHeaders, structField2CsvHeaderMapping)
    if err == nil { t.Errorf("No error occured.") }
  })

  t.Run("csvHeaders has missing field", func(t *testing.T){
    structField2CsvHeaderMapping := createStructField2CsvHeaderMapping()
    csvHeaders := createCSVHeaders()
    csvHeaders = csvHeaders[:len(csvHeaders) - 1]

    _, err := GetStructField2CSVHeaderIndexMapping[ParsedDataStruct](csvHeaders, structField2CsvHeaderMapping)
    if err == nil { t.Errorf("No error occured.") }
  })
}

func TestParseCSVRowSuccessWithErrorLog(t *testing.T) {
  var buf bytes.Buffer
  // Keep original log output
  originalOutput := log.Writer()
  defer log.SetOutput(originalOutput)
  log.SetOutput(&buf)

  t.Run("Failed to parse the row with invalid ParsedDataStruct", func(t *testing.T) {
    defer func() {
      buf.Reset()
    }()

    row := []string{"1", "1.1", "2024/03/03", "hoge"}
    type InvalidParsedDataStruct struct {
      IntField int
      // Float64Field float64 // missing field
      DateField time.Time
      StringField string
    }

    _, err := ParseCSVRow[InvalidParsedDataStruct](structField2CSVHeaderIndexMapping, row)
    if err != nil { t.Errorf("Error occured: %s.", err) }
    if !strings.Contains(buf.String(), "[ERROR] No field 'Float64Field' in struct 'InvalidParsedDataStruct'") {
      t.Errorf("No error log is outputted.")
    }
  })
}

func TestParseCSVRowSuccess(t *testing.T) {
  tests := map[string]struct {
    row []string
    expected ParsedDataStruct
  }{
    "Succeeded to parse the row": {row: []string{"1", "1.1", "2024/03/03", "hoge"}, expected: ParsedDataStruct{1, 1.1, time.Date(2024, time.March, 3, 0, 0, 0, 0, time.UTC), "hoge"}},
    "Succeeded to parse float64 field which value is int": {row: []string{"1", "1", "2024/03/03", "hoge"}, expected: ParsedDataStruct{1, 1, time.Date(2024, 3, 3, 0, 0, 0, 0, time.UTC), "hoge"}},
    "Succeeded to parse float64 field with a seperator": {row: []string{"1", "1,000.1", "2024/03/03", "hoge"}, expected: ParsedDataStruct{1, 1000.1, time.Date(2024, 3, 3, 0, 0, 0, 0, time.UTC), "hoge"}},
    "Succeeded to parse float64 field with seperators": {row: []string{"1", "1,001,000.1", "2024/03/03", "hoge"}, expected: ParsedDataStruct{1, 1001000.1, time.Date(2024, 3, 3, 0, 0, 0, 0, time.UTC), "hoge"}},
    "Succeeded to parse the row with more field": {row: []string{"1", "1.1", "2024/03/03", "hoge", "unexpectedValid"}, expected: ParsedDataStruct{1, 1.1, time.Date(2024, time.March, 3, 0, 0, 0, 0, time.UTC), "hoge"}},
  }

  for name, test := range tests {
    t.Run(name, func(t *testing.T){
      actual, _ := ParseCSVRow[ParsedDataStruct](structField2CSVHeaderIndexMapping, test.row)
      if actual != test.expected {
        t.Errorf("got %v, want: %v", actual, test.expected)
      }
    })
  }
}

func TestParseCSVRowFailure(t *testing.T) {
  tests := map[string]struct {
    row []string
  }{
    "Failed to parse the row with float value in int field": {row: []string{"1.1", "1.1", "2024/03/03", "hoge"}},
    "Failed to parse the row with invalid float value": {row: []string{"1", "invalid float", "2024/03/03", "hoge"}},
    "Failed to parse the row with invalid date value": {row: []string{"1", "1.1", "invalid date", "hoge"}},
    "Failed to parse the row with fewer field": {row: []string{"1.1", "1.1", "2024/03/03"}},
    // "Failed to parse the row with invalid ParsedDataStruct": {row: []string{"1", "1.1", "2024/03/03", "hoge"}, expected: InvalidParsedDataStruct},
  }

  for name, test := range tests {
    t.Run(name, func(t *testing.T){
      _, err := ParseCSVRow[ParsedDataStruct](structField2CSVHeaderIndexMapping, test.row)
      if err == nil { t.Errorf("No error occured.") }
    })
  }
}

func TestReadCSVFromBytesSuccess(t *testing.T) {
  t.Run("Read CSV from bytes", func(t *testing.T){
    rowString := "整数,小数,文字列,日付" + "\n" +
      `"1","1.1",オラオラ,2025/03/13` + "\n" +
      `"2,222","2",無駄無駄,2025/03/12` + "\n" +
      `"3,333,333","3,333.3",ホゲホゲ,2025/03/11`
    rowBytes := []byte(rowString)

    expectedHeader := []string{"整数", "小数", "文字列", "日付"}
    expectedRows := [][]string{
      {"1", "1.1", "オラオラ", "2025/03/13"},
      {"2,222", "2", "無駄無駄", "2025/03/12"},
      {"3,333,333", "3,333.3", "ホゲホゲ", "2025/03/11"},
    }

    actualHeader, actualRows, err := ReadCSV(rowBytes)
    if err != nil { t.Errorf("Error occured: %s.", err) }
    if !reflect.DeepEqual(actualHeader, expectedHeader) { t.Errorf("Failed to read CSV header. got: %v, want: %v", actualHeader, expectedHeader) }
    if !reflect.DeepEqual(actualRows, expectedRows) { t.Errorf("Failed to read CSV rows. got: %v, want: %v", actualRows, expectedRows)  }
  })
}

func TestReadCSVFromBytesFailure(t *testing.T) {
  t.Run("Read CSV from bytes with invalid header", func(t *testing.T){
    rowString := `"整数","小数,文字列,日付` + "\n" + // missing '"'
      `"1","1.1",オラオラ,2025/03/13`
    rowBytes := []byte(rowString)

    _, _, err := ReadCSV(rowBytes)
    if err == nil { t.Errorf("No error occured.") }
  })

  t.Run("Read CSV from bytes with invalid row", func(t *testing.T){
    rowString := "整数,小数,文字列,日付" + "\n" +
      `"1,"1.1",オラオラ,2025/03/13` + "\n" + // missing '"'
      `"2,222","2",無駄無駄,2025/03/12`
    rowBytes := []byte(rowString)

    _, _, err := ReadCSV(rowBytes)
    if err == nil { t.Errorf("No error occured.") }
  })
}

func TestReadCSVFromFileSuccess(t *testing.T) {
  t.Run("Read CSV from file", func(t *testing.T){
    csvPath := "testdata/test_data_valid.csv"

    expectedHeader := []string{"整数", "小数", "文字列", "日付"}
    expectedRows := [][]string{
      {"1", "1.1", "オラオラ", "2025/03/13"},
      {"2,222", "2", "無駄無駄", "2025/03/12"},
      {"3,333,333", "3,333.3", "ホゲホゲ", "2025/03/11"},
    }

    actualHeader, actualRows, err := ReadCSV(csvPath)
    if err != nil { t.Errorf("Error occured: %s.", err) }
    if !reflect.DeepEqual(actualHeader, expectedHeader) { t.Errorf("Failed to read CSV header. got: %v, want: %v", actualHeader, expectedHeader) }
    if !reflect.DeepEqual(actualRows, expectedRows) { t.Errorf("Failed to read CSV rows. got: %v, want: %v", actualRows, expectedRows)  }
  })
}

func TestReadCSVFromFileFailure(t *testing.T) {
  t.Run("Read CSV from not existing file", func(t *testing.T){
    csvPath := "testdata/not_existing.csv"

    _, _, err := ReadCSV(csvPath)
    if err == nil { t.Errorf("No error occured.") }
  })

  t.Run("Read CSV from file with invalid header", func(t *testing.T){
    csvPath := "testdata/test_data_with_invalid_header.csv"

    _, _, err := ReadCSV(csvPath)
    if err == nil { t.Errorf("No error occured.") }
  })

  t.Run("Read CSV from file with invalid row", func(t *testing.T){
    csvPath := "testdata/test_data_with_invalid_row.csv"

    _, _, err := ReadCSV(csvPath)
    if err == nil { t.Errorf("No error occured.") }
  })
}
