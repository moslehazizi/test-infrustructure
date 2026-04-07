package entity

import (
	"math/big"

	"gorm.io/gorm"
)

type Factorial struct {
	gorm.Model
	Input  string
	Output string
}

type FactorialEvent struct {
	Input  *big.Int `json:"input"`
	Output *big.Int `json:"output"`
}
