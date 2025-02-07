package rules

import (
	"net/http"

	"github.com/freedkr/cobra-canal/config"
	"github.com/freedkr/cobra-canal/event"
	"github.com/siddontang/go-log/log"
)

type httpAction interface {
	ServeHTTPStop(http.ResponseWriter, *http.Response)
}

type Ruler interface {
	httpAction
	GetDesc() string
	SetDesc(string)
	Start()
	Close() error
	HandleEvent(event.Event) error
	GetName() string
	SetLogger(*log.Logger)
	LoadConfig(config.RuleConfig) error
	SetNumber(int)
	GetNumber() int
	SetAggregator(config.Aggregatable)
	Reset() error
	IsClosed() bool
	RulerInfo() (RulerInfo, error)
}
