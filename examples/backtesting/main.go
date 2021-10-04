package main

import (
	"context"
	"fmt"

	"github.com/iltoga/ninjabot"
	"github.com/iltoga/ninjabot/examples/strategies"
	"github.com/iltoga/ninjabot/exchange"
	"github.com/iltoga/ninjabot/plot"
	"github.com/iltoga/ninjabot/plot/indicator"
	"github.com/iltoga/ninjabot/storage"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	settings := ninjabot.Settings{
		Pairs: []string{
			"BTCUSDT",
			"ETHUSDT",
			// "BATUSDT",
			"ADAUSDT",
			"BNBUSDT",
		},
	}

	strategy := new(strategies.CrossEMA2)

	suffix := "-2021-crash-to-today-30m"
	// suffix := "-30m"
	timeFrame := "30m"
	csvFeed, err := exchange.NewCSVFeed(
		strategy.Timeframe(),
		exchange.PairFeed{
			Pair:      "BTCUSDT",
			File:      fmt.Sprintf("testdata/btc%s.csv", suffix),
			Timeframe: timeFrame,
		},
		exchange.PairFeed{
			Pair:      "ETHUSDT",
			File:      fmt.Sprintf("testdata/eth%s.csv", suffix),
			Timeframe: timeFrame,
		},
		exchange.PairFeed{
			Pair:      "BATUSDT",
			File:      fmt.Sprintf("testdata/bat%s.csv", suffix),
			Timeframe: timeFrame,
		},
		exchange.PairFeed{
			Pair:      "ADAUSDT",
			File:      fmt.Sprintf("testdata/ada%s.csv", suffix),
			Timeframe: timeFrame,
		},
		exchange.PairFeed{
			Pair:      "BNBUSDT",
			File:      fmt.Sprintf("testdata/bnb%s.csv", suffix),
			Timeframe: "30m",
		},
		// exchange.PairFeed{
		// 	Pair:      "BTCUSDT",
		// 	File:      "testdata/btc-30m.csv",
		// 	Timeframe: "30m",
		// },
		// exchange.PairFeed{
		// 	Pair:      "ETHUSDT",
		// 	File:      "testdata/eth-30m.csv",
		// 	Timeframe: "30m",
		// },
		// exchange.PairFeed{
		// 	Pair:      "BATUSDT",
		// 	File:      "testdata/bat-30m.csv",
		// 	Timeframe: "30m",
		// },
		// exchange.PairFeed{
		// 	Pair:      "ADAUSDT",
		// 	File:      "testdata/ada-30m.csv",
		// 	Timeframe: "30m",
		// },
		// exchange.PairFeed{
		// 	Pair:      "BNBUSDT",
		// 	File:      "testdata/bnb-30m.csv",
		// 	Timeframe: "30m",
		// },
	)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.FromMemory()
	if err != nil {
		log.Fatal(err)
	}

	wallet := exchange.NewPaperWallet(
		ctx,
		"USDT",
		exchange.WithPaperAsset("USDT", 1000),
		exchange.WithDataFeed(csvFeed),
	)

	chart := plot.NewChart(plot.WithIndicators(
		indicator.EMA(13, "red"),
		indicator.EMA(21, "#000"),
		indicator.RSI(14, "purple"),
		indicator.Stoch(8, 3, "red", "blue"),
	))

	bot, err := ninjabot.NewBot(
		ctx,
		settings,
		wallet,
		strategy,
		ninjabot.WithBacktest(wallet),
		ninjabot.WithStorage(storage),
		ninjabot.WithCandleSubscription(chart),
		ninjabot.WithOrderSubscription(chart),
		ninjabot.WithLogLevel(log.WarnLevel),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = bot.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Print bot results
	fmt.Println(bot.Summary())
	wallet.Summary()
	err = chart.Start()
	if err != nil {
		log.Fatal(err)
	}
}
