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
var engineReady chan bool

//NextID getting ID
func NextID() int64 {
	orderID++
	return orderID
}

//NewContract const
func NewContract(symbol string) ib.Contract {
	return ib.Contract{
		Symbol:       symbol,
		SecurityType: "STK", //STK CFD
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

func engineLoop(ibmanager *ib.Engine) {
	engs := make(chan ib.EngineState)
	rc := make(chan ib.Reply)

	engine.SubscribeState(engs)
	engine.SubscribeAll(rc)

	engine.Send(&ib.RequestIDs{})

	for {
		select {
		case r := <-rc:
			//log.Printf("%s - RECEIVE %v",  reflect.TypeOf(r))
			switch r.(type) {

			case (*ib.ErrorMessage):
				r := r.(*ib.ErrorMessage)
				log.Printf("ID: %v Code:%3d Message:'%v'\n", r.ID(), r.Code, r.Message)

			// case (*ib.ManagedAccounts):
			// 	r := r.(*ib.ManagedAccounts)
			// 	for _, acct := range r.AccountsList {
			// 		log.Printf("%s: Account %v\n", acct)
			// 	}

			// case (*ib.Position):
			// 	r := r.(*ib.Position)
			// 	log.Printf("%s: C:%6v P:%10v AvgC:%10.2f\n", r.Contract.Symbol, r.Position, r.AverageCost)

			// case (*ib.OpenOrder):
			// 	r := r.(*ib.OpenOrder)
			// 	commission := FloatAdjustValue(r.OrderState.Commission)
			// 	maxcommission := FloatAdjustValue(r.OrderState.MaxCommission)
			// 	mincommission := FloatAdjustValue(r.OrderState.MinCommission)
			// 	log.Printf("%s OrderID: %v,%v Status: %-9v Symbol: %-5v Action   : %-4v  Quantity        : %4v %v %v l:%6.2f a:%6.2f c:%4.2f %4.2f/%4.2f\n", r.Order.OrderID, r.Order.ParentID, r.OrderState.Status, r.Contract.Symbol, r.Order.Action, r.Order.TotalQty, r.Order.TIF, r.Order.OrderType, r.Order.LimitPrice, r.Order.AuxPrice, commission, mincommission, maxcommission)

			// case (*ib.OrderStatus):
			// 	r := r.(*ib.OrderStatus)
			// 	log.Printf("%s OrderID: %v,%v Status: %-9v Filled: %5v Remaining: %5v AverageFillPrice: %6.2f - WH:'%s'\n", r.ID(), r.ParentID, r.Status, r.Filled, r.Remaining, r.AverageFillPrice, r.WhyHeld)

			// case (*ib.AccountValue):
			// 	r := r.(*ib.AccountValue)
			// 	if r.Currency == "USD" {
			// 		var show bool
			// 		switch r.Key.Key {
			// 		case "AvailableFunds":
			// 			show = true
			// 		case "BuyingPower":
			// 			show = true
			// 		case "TotalCashValue":
			// 			show = true
			// 		case "GrossPositionValue":
			// 			show = true
			// 		case "NetLiquidation":
			// 			show = true
			// 		case "UnrealizedPnL":
			// 			show = true
			// 		case "RealizedPnL":
			// 			show = true
			// 		case "AccruedCash":
			// 			show = true
			// 		default:
			// 			show = false
			// 		}
			// 		if show || gUpdateOverride {
			// 			log.Printf("%s: K:%-26v V:%20v\n", r.Key.Key, r.Value)
			// 		}
			// 	}

			// case (*ib.PortfolioValue):
			// 	r := r.(*ib.PortfolioValue)
			// 	log.Printf("%s: C:%6v P:%10v AvgC:%10.2f uPNL:%8.2f PNL:%8.2f\n", r.Contract.Symbol, r.Position, r.AverageCost, r.UnrealizedPNL, r.RealizedPNL)

			// case (*ib.AccountSummary):
			// 	r := r.(*ib.AccountSummary)
			// 	log.Printf("%s: K:%-26v V:%20v\n", r.Key.Key, r.Value)

			// case (*ib.ExecutionData):
			// 	r := r.(*ib.ExecutionData)
			// 	item, ok := ibmanager.elog[r.Exec.ExecID]
			// 	if !ok {
			// 		item = new(ExecutionInfo)
			// 		ibmanager.elog[r.Exec.ExecID] = item
			// 	}
			// 	item.ExecutionData = *r

			// case (*ib.CommissionReport):
			// 	r := r.(*ib.CommissionReport)
			// 	item, ok := ibmanager.elog[r.ExecutionID]
			// 	if !ok {
			// 		item = new(ExecutionInfo)
			// 		ibmanager.elog[r.ExecutionID] = item
			// 	}
			// 	item.Commission = *r

			// case (*ib.AccountSummaryEnd):
			// 	r := r.(*ib.AccountSummaryEnd)

			// 	if gCancel {
			// 		req := &ib.CancelAccountSummary{}
			// 		req.SetID(r.ID())
			// 		engine.Send(req)
			// 	}

			// case (*ib.ExecutionDataEnd):
			// 	var keys TimeSlice
			// 	for _, k := range ibmanager.elog {
			// 		keys = append(keys, k)
			// 	}

			// 	sort.Sort(keys)

			// 	for _, x := range keys {
			// 		log.Printf("%s: %v %4d %-7s %s %4d %7.2f %4d %7.2f %6.2f %s\n",

			// 			x.ExecutionData.Exec.Time.Format("15:04:05"),
			// 			x.ExecutionData.Exec.OrderID,
			// 			x.ExecutionData.Contract.Symbol,
			// 			x.ExecutionData.Exec.Side,
			// 			x.ExecutionData.Exec.Shares,
			// 			x.ExecutionData.Exec.Price,
			// 			x.ExecutionData.Exec.CumQty,
			// 			x.ExecutionData.Exec.AveragePrice,
			// 			x.Commission.Commission,
			// 			x.ExecutionData.Exec.Exchange)
			// 	}

			// case (*ib.RealtimeBars):
			// 	r := r.(*ib.RealtimeBars)

			// 	symbol, ok := ibmanager.realtimeMap[r.ID()]
			// 	if !ok {
			// 		symbol = ""
			// 	}

			// 	log.Printf("%10s: %v - Open: %10.2f Close: %10.2f Low %10.2f High %10.2f Volume %10.2f Count %10v WAP %10.2f\n",
			// 		symbol,
			// 		time.Unix(r.Time, 0).Format("15:04:05"),
			// 		r.Open,
			// 		r.Close,
			// 		r.Low,
			// 		r.High,
			// 		r.Volume,
			// 		r.Count,
			// 		r.WAP)

			// case (*ib.PositionEnd):

			// case (*ib.AccountDownloadEnd):
			// 	if gCancel {
			// 		req := &ib.RequestAccountUpdates{}
			// 		req.Subscribe = false
			// 		engine.Send(req)
			// 	}

			// case (*ib.OpenOrderEnd):

			// case (*ib.ContractDataEnd):

			// case (*ib.TickSnapshotEnd):

			// case (*ib.AccountUpdateTime):

			case (*ib.NextValidID):
				r := r.(*ib.NextValidID)
				orderID = r.OrderID
				log.Printf("OrderId=%v", orderID)
				engineReady <- true

			default:
				// log.Printf("%s - RECEIVE %v", reflect.TypeOf(r))
				log.Printf("%#v\n", r)
			}
		case newstate := <-engs:
			log.Printf("ERROR: %v\n", newstate)
			if newstate != ib.EngineExitNormal {
				log.Fatalf("ERROR: %v", engine.FatalError())
			}
			return
		}
	}
}

func main() {
	engineReady = make(chan bool)

	var err error
	engine, err = ib.NewEngine(ib.EngineOptions{Gateway: "localhost:4004"})

	time.Sleep(time.Second)
	go engineLoop(engine)

	if err != nil {
		log.Fatalf("error creating engine %v ", err)
	}

	defer engine.Stop()

	if engine.State() != ib.EngineReady {
		log.Fatalf("Engine is not ready")
	}

	<-engineReady

	order("AAPL", 1, 1, true)

	time.Sleep(3 * time.Second)
}
