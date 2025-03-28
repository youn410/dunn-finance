package csvreader

import (
  "bytes"
  "encoding/csv"
  "fmt"
  "io"
  "log"
  "math"
  "os"
  "reflect"
  "strings"
  "strconv"
  "time"
)

// Input
//   - csvHeaders: For example, ["フィールド1", "フィールド2", ... ]
//   - structField2CsvHeaderMapping: {"Field1": "フィールド1", "Field2": "フィールド2", ... }
// Return
//   - {"Field1": 0, "Field2": 1, ... }
func GetStructField2CSVHeaderIndexMapping[T any](csvHeaders []string, structField2CsvHeaderMapping map[string]string) (map[string]int, error) {
  structField2CSVHeaderIndexMapping := make(map[string]int)

  var data T
  structType := reflect.TypeOf(data)
  for i := 0; i < structType.NumField(); i++ {
    field := structType.Field(i)
    fieldName := field.Name

    mapped := false
    if targetCSVHeader, exists := structField2CsvHeaderMapping[fieldName]; exists {
      for idx, csvHeader := range csvHeaders {
        if targetCSVHeader == csvHeader {
          structField2CSVHeaderIndexMapping[fieldName] = idx
          mapped = true
          break
        }
      }
    } else {
      return structField2CSVHeaderIndexMapping, fmt.Errorf("[ERROR] Field '%s' not in field2headerMap '%v'", fieldName, structField2CsvHeaderMapping)
    }

    if !mapped {
      return structField2CSVHeaderIndexMapping, fmt.Errorf("[ERROR] Field '%s' not in CSV Headers: '%v'", fieldName, csvHeaders)
    }
  }

  return structField2CSVHeaderIndexMapping, nil
}

// Input
//   T: parsed csv one row data struct
//   structField2CSVHeaderIndexMapping: For example, {"Field1": 0, "Filed2": 1, ... }
//   row: CSV row. For example, []string{"1", "ora", ... }
func ParseCSVRow[T any](structField2CSVHeaderIndexMapping map[string]int, row []string) (T, error) {
  var data T

  maxIndexInMapping := math.MinInt
  for _, v := range structField2CSVHeaderIndexMapping {
    if v > maxIndexInMapping {
      maxIndexInMapping = v
    }
  }
  if maxIndexInMapping + 1 > len(row) {
    return data, fmt.Errorf("[ERROR] Failed to parse csv row: unexpected row length")
  }

  structValue := reflect.ValueOf(&data).Elem()
  structType := structValue.Type()
  for structField, index := range structField2CSVHeaderIndexMapping {
    log.Printf("[INFO] struct field: %s, index: %d\n", structField, index)

    field := structValue.FieldByName(structField)
    if !field.IsValid() || !field.CanSet() {
      log.Printf("[ERROR] No field '%s' in struct '%s'\n", structField, structType.Name())
      continue
    }

    switch field.Kind() {
      case reflect.Int, reflect.Int64:
        fieldValue, err := strconv.ParseInt(strings.ReplaceAll(row[index], ",", ""), 10, 64)
        if err != nil {
          return data, fmt.Errorf("[ERROR] Failed to ParseInt: err=%s, row=%q, field=%s", err, row, field)
        } else {
          field.SetInt(fieldValue)
        }
      case reflect.Float64:
        fieldValue, err := strconv.ParseFloat(strings.ReplaceAll(row[index], ",", ""), 64)
        if err != nil {
          return data, fmt.Errorf("[ERROR] Failed to ParseFloat: err=%s, row=%q, field=%s", err, row, field)
        } else {
          field.SetFloat(fieldValue)
        }
      case reflect.String:
        field.SetString(row[index])
      case reflect.Struct:
        detailedStructField, _ := structType.FieldByName(structField)
        if detailedStructField.Type == reflect.TypeOf(time.Time{}) {
          date, err := time.Parse("2006/01/02", row[index])
          if err != nil { return data, fmt.Errorf("[ERROR] Failed to Parse date: err=%s, row=%q, field=%s", err, row, field) }
          field.Set(reflect.ValueOf(date))
        }
    }
  }

  return data, nil
}

type SourceType interface {
  ~string | ~[]byte
}

// Input
//   - source: CSV data. file path or bytes.
// Return
//  - CSV header array
//  - CSV rows 2D array
//  - error
func ReadCSV[S SourceType](source S) ([]string, [][]string, error) {
  var reader *csv.Reader

  switch v := any(source).(type) {
    case string: // File path
      file, err := os.Open(v)
      if err != nil { return nil, nil, fmt.Errorf("Failed to open file: %w", err) }
      defer file.Close()

      reader = csv.NewReader(file)
    case []byte:
      reader = csv.NewReader(bytes.NewReader(v))
  }

  header, err := reader.Read()
  if err != nil { return nil, nil, fmt.Errorf("Failed to read header: %w", err) }

  var rows [][]string

  for {
    row, err := reader.Read()
    if err == io.EOF { break }
    if err != nil { return nil, nil, fmt.Errorf("Failed to read row: %w", err) }

    rows = append(rows, row)
  }

  return header, rows, nil
}
