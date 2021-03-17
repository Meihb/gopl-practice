// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 244.

// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"os"
	"time"
)

//!+

func main() {
	// ...create abort channel...

	//!-


	//!+abort
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()
	//!-abort

	//!+
	fmt.Println("Commencing countdown.  Press return to abort.")

	select {
	case <-time.After(10 * time.Second): //意思是10s后 从chan 中返回 一值 time.After函数会立即返回一个channel，并起一个新的goroutine在经过特定的时间后向该channel发送一个独立的值。
		// Do nothing.
	case <-abort:
		fmt.Println("Launch aborted!")
		return

	default: //不再有阻塞功能

	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
