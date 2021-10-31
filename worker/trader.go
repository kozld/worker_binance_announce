package worker

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"

	"github.com/stdi0/worker_binance_announce/config"
)

type Trader struct {
	conf *config.TraderConfig
}

func NewTrader(conf *config.TraderConfig) *Trader {
	return &Trader{conf}
}

func (t *Trader) CreateOrder(currency string, test bool) error {
	client := gateapi.NewAPIClient(gateapi.NewConfiguration())
	if test {
		client.ChangeBasePath("https://fx-api-testnet.gateio.ws/api/v4")
	}
	ctx := context.WithValue(context.Background(), gateapi.ContextGateAPIV4, gateapi.GateAPIV4{
		Key:    t.conf.GateIOApiKey,
		Secret: t.conf.GateIOApiSecret,
	})

	currencyPair := fmt.Sprintf("%s_USDT", currency)
	//currency := "USDT"
	cp, _, err := client.SpotApi.GetCurrencyPair(ctx, currencyPair)
	if err != nil {
		return err
	}
	//log.Printf("testing against currency pair: %s\n", cp.Id)
	//minAmount := cp.MinQuoteAmount

	tickers, _, err := client.SpotApi.ListTickers(ctx, &gateapi.ListTickersOpts{CurrencyPair: optional.NewString(cp.Id)})
	if err != nil {
		return err
	}

	lastPrice, err := strconv.ParseFloat(tickers[0].Last, 64)
	if err != nil {
		return err
	}

	fmt.Println("[PRICE]", lastPrice)
	lastPrice = lastPrice + (lastPrice * 0.3)
	fmt.Println("[PRICE + 30%]", lastPrice)

	// better avoid using float, take the following decimal library for example
	// `go get github.com/shopspring/decimal`
	//orderAmount := decimal.RequireFromString(minAmount).Mul(decimal.NewFromInt32(2))
	//
	//balance, _, err := client.SpotApi.ListSpotAccounts(ctx, &gateapi.ListSpotAccountsOpts{Currency: optional.NewString(currency)})
	//if err != nil {
	//	panicGateError(err)
	//}
	//if decimal.RequireFromString(balance[0].Available).Cmp(orderAmount) < 0 {
	//	log.Fatal("balance not enough")
	//}

	newOrder := gateapi.Order{
		CurrencyPair: cp.Id,
		Type:         "limit",
		Account:      "spot", // create spot order. set to "margin" if creating margin orders
		Side:         "buy",
		Amount:       "100",
		Price:        fmt.Sprintf("%f", lastPrice), // use last price
		TimeInForce:  "gtc",
		AutoBorrow:   false,
	}
	log.Printf("place a spot %s order in %s with amount %s and price %s\n", newOrder.Side, newOrder.CurrencyPair, newOrder.Amount, newOrder.Price)

	createdOrder, _, err := client.SpotApi.CreateOrder(ctx, newOrder)
	if err != nil {
		return err
	}
	log.Printf("order created with ID: %s, status: %s\n", createdOrder.Id, createdOrder.Status)

	//if createdOrder.Status == "open" {
	//	order, _, err := client.SpotApi.GetOrder(ctx, createdOrder.Id, createdOrder.CurrencyPair, nil)
	//	if err != nil {
	//		return err
	//	}
	//	log.Printf("order %s filled: %s, left: %s\n", order.Id, order.FilledTotal, order.Left)
	//	//result, _, err := client.SpotApi.CancelOrder(ctx, createdOrder.Id, createdOrder.CurrencyPair)
	//	//if err != nil {
	//	//	panicGateError(err)
	//	//}
	//	//if result.Status == "cancelled" {
	//	//	log.Printf("order %s cancelled\n", createdOrder.Id)
	//	//}
	//}

	return nil

	//else {
	//	// order finished
	//	trades, _, err := client.SpotApi.ListMyTrades(ctx, createdOrder.CurrencyPair,
	//		&gateapi.ListMyTradesOpts{OrderId: optional.NewString(createdOrder.Id)})
	//	if err != nil {
	//		panicGateError(err)
	//	}
	//	for _, t := range trades {
	//		log.Printf("order %s filled %s with price: %s\n", t.OrderId, t.Amount, t.Price)
	//	}
	//}
}
