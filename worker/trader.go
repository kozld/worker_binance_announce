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
	cp, _, err := client.SpotApi.GetCurrencyPair(ctx, currencyPair)
	if err != nil {
		return err
	}

	tickers, _, err := client.SpotApi.ListTickers(ctx, &gateapi.ListTickersOpts{CurrencyPair: optional.NewString(cp.Id)})
	if err != nil {
		return err
	}

	// TODO: Check tickers list length > 0

	lastPrice, err := strconv.ParseFloat(tickers[0].Last, 64)
	if err != nil {
		return err
	}
	log.Printf("(%s) Last price: %f", currency, lastPrice)

	// Increase last price +30%
	lastPrice = lastPrice + (lastPrice * 0.3)
	log.Printf("(%s) Target price (+30 percent): %f", currency, lastPrice)

	// Figure out amount
	amount := (t.conf.QuantityUSDT - 1) / lastPrice

	newOrder := gateapi.Order{
		CurrencyPair: cp.Id,
		Type:         "limit",
		Account:      "spot", // create spot order. set to "margin" if creating margin orders
		Side:         "buy",
		Amount:       fmt.Sprintf("%f", amount),
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

	return nil
}
