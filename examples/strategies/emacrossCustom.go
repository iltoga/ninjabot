package strategies

import (
	"github.com/rodrigo-brito/ninjabot"
	"github.com/rodrigo-brito/ninjabot/service"

	"github.com/markcheno/go-talib"
	log "github.com/sirupsen/logrus"
)

type CrossEMACustom struct {
	buys  int
	sells int
}

func (e CrossEMACustom) Timeframe() string {
	return "4h"
}

func (e CrossEMACustom) WarmupPeriod() int {
	return 21
}

func (e CrossEMACustom) Indicators(df *ninjabot.Dataframe) {
	df.Metadata["ema8"] = talib.Ema(df.Close, 13)
	df.Metadata["ema21"] = talib.Sma(df.Close, 21)
	df.Metadata["RSI14"] = talib.Rsi(df.Close, 14)
}

func (e *CrossEMACustom) OnCandle(df *ninjabot.Dataframe, broker service.Broker) {
	closePrice := df.Close.Last(0)
	assetPosition, quotePosition, err := broker.Position(df.Pair)
	if err != nil {
		log.Error(err)
	}

	if quotePosition > 10 && df.Metadata["ema8"].Crossover(df.Metadata["ema21"]) {
		_, err := broker.CreateOrderMarketQuote(ninjabot.SideTypeBuy, df.Pair, quotePosition/2)
		if err != nil {
			log.WithFields(map[string]interface{}{
				"pair":  df.Pair,
				"side":  ninjabot.SideTypeBuy,
				"close": closePrice,
				"asset": assetPosition,
				"quote": quotePosition,
			}).Error(err)
		}
	}

	if assetPosition > 0 &&
		df.Metadata["ema8"].Crossunder(df.Metadata["ema21"]) {
		_, err := broker.CreateOrderMarket(ninjabot.SideTypeSell, df.Pair, assetPosition)
		if err != nil {
			log.WithFields(map[string]interface{}{
				"pair":  df.Pair,
				"side":  ninjabot.SideTypeSell,
				"close": closePrice,
				"asset": assetPosition,
				"quote": quotePosition,
				"size":  assetPosition,
			}).Error(err)
		}
	}
}
