package main

import (
	"context"
	"fmt"
	"time"

	. "github.com/hadrianl/ibapi"
)

func main() {
	var err error
	ibwrapper := &Wrapper{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	ic := NewIbClient(ibwrapper)
	ic.SetContext(ctx)
	err = ic.Connect("127.0.0.1", 7497, 0) // 7497 for TWS, 4002 for IB Gateway
	if err != nil {
		fmt.Println("Connect failed:", err)
	}

	err = ic.HandShake()
	if err != nil {
		fmt.Println("HandShake failed:", err)
	}

	ic.Run()

	time.Sleep(time.Second * 1)

	contract := Contract{Symbol: "EUR", SecurityType: "CASH", Currency: "GBP", Exchange: "IDEALPRO"}
	fmt.Println("contract:", contract)

	ic.ReqHistoricalData(ic.GetReqID(), &contract, "BidAsk", "4800 S", "1 min", "TRADES", false, 1, true, nil)

	ic.LoopUntilDone()
}
