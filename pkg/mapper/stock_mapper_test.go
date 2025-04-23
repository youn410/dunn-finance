package mapper_test

import (
  "testing"

  "dunn-finance/pkg/dto"
  "dunn-finance/pkg/mapper"
  "dunn-finance/pkg/model"
)

func TestToStockDTO_Success(t *testing.T) {
  code := "1234"
  name := "テスト会社"
  stock := &model.Stock{Code: code, Name: name}

  stockDto := mapper.ToStockDTO(stock)
  if code != stockDto.Code { t.Errorf("got %s, want %s", stockDto.Code, code) }
  if name != stockDto.Name { t.Errorf("got %s, want %s", stockDto.Name, name) }
}

func TestToStockModel_Success(t *testing.T) {
  code := "1234"
  name := "テスト会社"
  stockDto := &dto.StockDTO{Code: code, Name: name}

  stock := mapper.ToStockModel(stockDto)
  if code != stock.Code { t.Errorf("got %s, want %s", stock.Code, code) }
  if name != stock.Name { t.Errorf("got %s, want %s", stock.Name, name) }
}
