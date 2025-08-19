package models

import (
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

type CommonModel struct {
	gorm.Model
}

// tbl_seven_fifty
type TblSevenFifty struct {
	ID     uint   `gorm:"primarykey,column:id"`
	Symbol string `gorm:"column:symbol"`
}

func (TblSevenFifty) TableName() string {
	return "tbl_seven_fifty"
}

var _ Tabler = (*TblSevenFifty)(nil)
