package csvreader_test

import (
  "testing"

  "dunn-finance/pkg/csvreader"
)

var fieldMap = map[int]string{
  0: "Yyyymmdd",
  1: "OpenPrice",
  2: "HighPrice",
  3: "LowPrice",
  4: "ClosePrice",
  5: "DMAPrice5",
  6: "DMAPrice25",
  7: "DMAPrice75",
  8: "VMAP",
  9: "Volume",
  10: "VMA5",
  11: "VMA25",
}

func TestLoadAdjustedDailyOHLCVsFromCSV_Success(t *testing.T) {
  records, err := csvreader.LoadAdjustedDailyOHLCVsFromCSV("5253", "testdata/sbi_timechart_5253_20250720.csv", fieldMap, true, 0, 0)
  if err != nil { t.Fatal(err) }

  // validate record of 2023/03/30
  record := records[565]
  if record.Yyyymmdd != "20230330" { t.Errorf("Expected: 20230330, but got: %s", record.Yyyymmdd) }
  if record.Code != "5253" { t.Errorf("Expected: 5253, but got: %s", record.Code) }
  if *record.OpenPrice != 1465 { t.Errorf("Expected: 1465, but got: %f", *record.OpenPrice) }
  if *record.HighPrice != 1486 { t.Errorf("Expected: 1486, but got: %f", *record.HighPrice) }
  if *record.LowPrice != 1303 { t.Errorf("Expected: 1303, but got: %f", *record.LowPrice) }
  if *record.ClosePrice != 1326 { t.Errorf("Expected: 1326, but got: %f", *record.ClosePrice) }
  if record.DMAPrice5 != nil { t.Errorf("Expected: nil, but got: %f", *record.DMAPrice5) }
  if record.DMAPrice25 != nil { t.Errorf("Expected: nil, but got: %f", *record.DMAPrice25) }
  if record.DMAPrice75 != nil { t.Errorf("Expected: nil, but got: %f", *record.DMAPrice75) }
  if *record.VMAP != 1380.4325 { t.Errorf("Expected: nil, but got: %f", *record.VMAP) }
  if *record.Volume != 9858700 { t.Errorf("Expected: nil, but got: %f", *record.Volume) }
  if record.VMA5 != nil { t.Errorf("Expected: nil, but got: %f", *record.VMA5) }
  if record.VMA25 != nil { t.Errorf("Expected: nil, but got: %f", *record.VMA25) }
}

func TestLoadAdjustedDailyOHLCVsFromCSV_with_offset_limit_Success(t *testing.T) {
  records, err := csvreader.LoadAdjustedDailyOHLCVsFromCSV("5253", "testdata/sbi_timechart_5253_20250720.csv", fieldMap, true, 9, 5)
  if err != nil { t.Fatal(err) }
  if len(records) != 5 { t.Errorf("Expected record length: 5, but is %d", len(records)) }
  if records[0].Yyyymmdd != "20250707" { t.Errorf("Expected: 20250707, but got %s", records[0].Yyyymmdd) }
  if records[4].Yyyymmdd != "20250701" { t.Errorf("Expected: 20250701, but got %s", records[1].Yyyymmdd) }

  records, err = csvreader.LoadAdjustedDailyOHLCVsFromCSV("5253", "testdata/sbi_timechart_5253_20250720.csv", fieldMap, true, 14, 5)
  if err != nil { t.Fatal(err) }
  if len(records) != 5 { t.Errorf("Expected record length: 5, but is %d", len(records)) }
  if records[0].Yyyymmdd != "20250630" { t.Errorf("Expected: 20250630, but got %s", records[0].Yyyymmdd) }
  if records[4].Yyyymmdd != "20250624" { t.Errorf("Expected: 20250624, but got %s", records[1].Yyyymmdd) }
}
