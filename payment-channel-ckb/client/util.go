package client

import (
	"math/big"
)

// CKByteToShannon converts a given amount in CKByte to Shannon.
func CKByteToShannon(ckbyteAmount *big.Float) (shannonAmount *big.Int) {
	shannonPerCKByte := new(big.Int).Exp(big.NewInt(10), big.NewInt(8), nil)
	shannonPerCKByteFloat := new(big.Float).SetInt(shannonPerCKByte)
	shannonAmountFloat := new(big.Float).Mul(ckbyteAmount, shannonPerCKByteFloat)
	shannonAmount, _ = shannonAmountFloat.Int(nil)
	return shannonAmount
}

// ShannonToCKByte converts a given amount in Shannon to CKByte.
func ShannonToCKByte(shannonAmount *big.Int) *big.Float {
	shannonPerCKByte := new(big.Int).Exp(big.NewInt(10), big.NewInt(8), nil)
	shannonPerCKByteFloat := new(big.Float).SetInt(shannonPerCKByte)
	shannonAmountFloat := new(big.Float).SetInt(shannonAmount)
	return new(big.Float).Quo(shannonAmountFloat, shannonPerCKByteFloat)
}
