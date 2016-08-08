package main

import (
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gofinance/ib"
)

var engine *ib.Engine

var orderID int64

//NextID getting ID
func NextID() int64 {
	orderID++
	return orderID
}

//NewContract const
func NewContract(symbol string) ib.Contract {
	return ib.Contract{
		Symbol:       symbol,
		SecurityType: "CFD",
		Exchange:     "SMART",
		Currency:     "USD",
	}
}

//NewOrder const
func NewOrder(gtc bool) (ib.Order, error) {
	order, err := ib.NewOrder()

	order.TIF = "DAY"
	if gtc {
		order.TIF = "GTC"
	}

	return order, err
}

func order(symbol string, qua int64, price float64, gtc bool) {

	request := ib.PlaceOrder{
		Contract: NewContract(symbol),
	}

	request.Order, _ = NewOrder(gtc)

	if qua < 0 {
		request.Order.Action = "SELL"
	} else {
		request.Order.Action = "BUY"
	}

	request.Order.TotalQty = int64(math.Abs(float64(qua)))

	priceLog := ""

	if price == 0 {
		request.Order.OrderType = "MKT"
		priceLog = request.Order.OrderType
	} else {
		request.Order.OrderType = "LMT"
		request.Order.LimitPrice = price
		priceLog = strconv.FormatFloat(price, 'f', 2, 64)
	}

	request.SetID(NextID())

	engine.Send(&request)
	log.Printf("%s %d %s @ %s", request.Order.Action, request.Order.TotalQty, symbol, priceLog)
}

func main() {
	var err error
	engine, err = ib.NewEngine(ib.EngineOptions{Gateway: "localhost:4001"})

	if err != nil {
		log.Fatalf("error creating engine %v ", err)
	}

	time.Sleep(time.Second)

	order("AAPL", 1, 0, true)

	time.Sleep(3 * time.Second)
}
