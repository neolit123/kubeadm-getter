// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/neolit123/tokenized-getter/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
