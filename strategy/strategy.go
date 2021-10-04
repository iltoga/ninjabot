package strategy

import (
	"github.com/iltoga/ninjabot/model"
	"github.com/iltoga/ninjabot/service"
)

type Strategy interface {
	Timeframe() string
	WarmupPeriod() int
	Indicators(dataframe *model.Dataframe)
	OnCandle(dataframe *model.Dataframe, broker service.Broker)
}
