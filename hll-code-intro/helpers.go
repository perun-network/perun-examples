// Copyright (c) 2021, PolyCrypt GmbH, Germany. All rights reserved.
// This file is part of perun-tutorial. Use of this source code is
// governed by the Apache 2.0 license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/big"
	"strconv"
)

func noError(err error) {
	if err != nil {
		panic(err)
	}
}

func ethToWei(eth float64) *big.Int {
	//1 Ether = 10^18 Wei
	var ethPerWei = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	wei, _ := new(big.Float).Mul(big.NewFloat(eth), new(big.Float).SetInt(ethPerWei)).Int(new(big.Int))
	return wei
}

func queryNumber() int {
	for {
		fmt.Println("Enter an Integer")
		var input string
		fmt.Scanln(&input)

		value, err := strconv.Atoi(input)
		if err == nil {
			return value
		}
	}
}
