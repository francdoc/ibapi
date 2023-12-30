package main

import (
	"fmt"
	"time"

	. "github.com/hadrianl/ibapi"
)

func main() {
	var err error
	ibwrapper := &Wrapper{}
	ic := NewIbClient(ibwrapper)
	err = ic.Connect("127.0.0.1", 4002, 0) // 7497 for TWS, 4002 for IB Gateway
	if err != nil {
		fmt.Println("Connect failed:", err)
	}

	err = ic.HandShake()
	if err != nil {
		fmt.Println("HandShake failed:", err)
	}

	ic.Run()

	time.Sleep(time.Second * 20)

	ic.LoopUntilDone()
}
