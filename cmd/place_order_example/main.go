package main

import (
	"context"
	"fmt"
	"time"

	. "github.com/hadrianl/ibapi"
	"github.com/shopspring/decimal"
)

func main() {
	var err error
	ibwrapper := &Wrapper{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	ic := NewIbClient(ibwrapper)
	ic.SetContext(ctx)
	err = ic.Connect("127.0.0.1", 4002, 0) // 7497 for TWS, 4002 for IB Gateway
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

	ic.ReqTickByTickData(1, &contract, "BidAsk", 0, true)
	ic.ReqHistoricalData(ic.GetReqID(), &contract, "BidAsk", "4800 S", "1 min", "TRADES", false, 1, true, nil)
	ic.ReqAutoOpenOrders(true)

	quantity, err := decimal.NewFromString("1")
	fmt.Println("quantity:", quantity, "err:", err)

	mktOrder_buy := NewMarketOrder("BUY", quantity)
	lmtOrder_buy := NewLimitOrder("BUY", 144, quantity)

	ic.PlaceOrder(ibwrapper.GetNextOrderID(), &contract, mktOrder_buy)
	ic.PlaceOrder(ibwrapper.GetNextOrderID(), &contract, lmtOrder_buy)

	mktOrder_sell := NewMarketOrder("SELL", quantity)
	lmtOrder_sell := NewLimitOrder("SELL", 10, quantity)

	ic.PlaceOrder(ibwrapper.GetNextOrderID(), &contract, mktOrder_sell)
	ic.PlaceOrder(ibwrapper.GetNextOrderID(), &contract, lmtOrder_sell)

	time.Sleep(time.Second * 5)

	fmt.Println(" ================= Requesting positions ================= ")
	ic.ReqPositions()
	fmt.Println(" ================= Requesting account summary ================= ")
	ic.ReqOpenOrders()

	ic.LoopUntilDone()
}
