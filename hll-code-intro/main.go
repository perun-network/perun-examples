// Copyright (c) 2021, PolyCrypt GmbH, Germany. All rights reserved.
// This file is part of perun-tutorial. Use of this source code is
// governed by the Apache 2.0 license that can be found in the LICENSE file.

package main

import "fmt"

func main() {
	// 1. Set up everything.
	bob := setupNode()

	// 2. Open a channel where Alice and Bob both have 500 ETH.
	initAmount := ethToWei(500)
	channel := bob.openChannel(initAmount)

	// 3. Send some payments over the channel.
	for i := queryNumber(); i > 0; i-- {
		sendBalance(channel, ethToWei(1))
	}

	// Wait for one Enter key press.
	fmt.Scanln()
	// 4. Close the channel.
	closeChannel(channel)
}
