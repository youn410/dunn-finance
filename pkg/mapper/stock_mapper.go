package mapper

import (
  "dunn-finance/pkg/dto"
  "dunn-finance/pkg/model"
)

func ToStockDTO(m *model.Stock) *dto.StockDTO {
  return &dto.StockDTO{
    Code: m.Code,
    Name: m.Name,
  }
}

func ToStockModel(d *dto.StockDTO) *model.Stock {
  return &model.Stock{
    Code: d.Code,
    Name: d.Name,
  }
}
