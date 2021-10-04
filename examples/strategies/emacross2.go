package strategies

import (
	"github.com/iltoga/ninjabot"
	"github.com/iltoga/ninjabot/service"

	"github.com/markcheno/go-talib"
	log "github.com/sirupsen/logrus"
)

type CrossEMA2 struct{}

func (e CrossEMA2) Timeframe() string {
	return "30m"
}

func (e CrossEMA2) WarmupPeriod() int {
	return 50
}

func (e CrossEMA2) Indicators(df *ninjabot.Dataframe) {
	df.Metadata["ema20"] = talib.Ema(df.Close, 20)
	df.Metadata["ema50"] = talib.Ema(df.Close, 50)
}

func (e *CrossEMA2) OnCandle(df *ninjabot.Dataframe, broker service.Broker) {
	closePrice := df.Close.Last(0)
	assetPosition, quotePosition, err := broker.Position(df.Pair)
	if err != nil {
		log.Error(err)
	}

	if quotePosition > 10 && df.Metadata["ema20"].Crossover(df.Metadata["ema50"]) {
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
		df.Metadata["ema20"].Crossunder(df.Metadata["ema50"]) {
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
