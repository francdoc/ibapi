package main

import (
	"context"
	"fmt"
	"time"

	. "github.com/hadrianl/ibapi"
	filter "github.com/milosgajdos/go-estimate"
	"github.com/milosgajdos/go-estimate/kalman/ukf"
	"github.com/milosgajdos/go-estimate/noise"
	"github.com/milosgajdos/go-estimate/sim"
	"github.com/shopspring/decimal"

	"gonum.org/v1/gonum/mat"
)

type invalidModel struct {
	filter.Model
}

func (m *invalidModel) SystemDims() (nx, nu, ny, nz int) {
	return -10, 0, 8, 0
}

var (
	okModel  *sim.BaseModel
	badModel *invalidModel
	ic       *sim.InitCond
	q        filter.Noise
	r        filter.Noise
	c        *ukf.Config
	u        *mat.VecDense
	z        *mat.VecDense
)

func setup() {
	u = mat.NewVecDense(1, []float64{-1.0})
	z = mat.NewVecDense(1, []float64{-1.5})

	// initial condition
	initState := mat.NewVecDense(2, []float64{1.0, 3.0})
	initCov := mat.NewSymDense(2, []float64{0.25, 0, 0, 0.25})
	ic = sim.NewInitCond(initState, initCov)

	// state and output noise
	q, _ = noise.NewGaussian([]float64{0, 0}, initCov)
	r, _ = noise.NewGaussian([]float64{0}, mat.NewSymDense(1, []float64{0.25}))

	A := mat.NewDense(2, 2, []float64{1.0, 1.0, 0.0, 1.0})
	B := mat.NewDense(2, 1, []float64{0.5, 1.0})
	C := mat.NewDense(1, 2, []float64{1.0, 0.0})
	D := mat.NewDense(1, 1, []float64{0.0})

	okModel = &sim.BaseModel{A: A, B: B, C: C, D: D}
	badModel = &invalidModel{okModel}

	c = &ukf.Config{
		Alpha: 0.75,
		Beta:  2.0,
		Kappa: 3.0,
	}
}

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
