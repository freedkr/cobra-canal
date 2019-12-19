package rules

import (
	"sync"

	"github.com/always-waiting/cobra-canal/config"
	cobraErrors "github.com/always-waiting/cobra-canal/errors"
	"github.com/always-waiting/cobra-canal/event"
	"github.com/juju/errors"
	"github.com/siddontang/go-log/log"
)

const (
	LOAD_ERR1 = "规则名为空"
	LOAD_ERR2 = "规则(%s)未注册"
	LOAD_ERR3 = "生成%s规则失败"
	HEADER    = ">>>>>>>>开始处理<<<<<<<<"
)

type Rule struct {
	ruler        []Ruler                 `description:"规则,数组是指统一规则的多个worker"`
	eventChannel chan event.Event        `description:"事件队列"`
	errHr        *cobraErrors.ErrHandler `description:"错误处理对象"`
	isRulerClose chan bool
	closed       bool
	isReady      bool
	Log          *log.Logger
	rulerNum     int
	aggregator   config.Aggregatable
	name         string
}

var ruleMakers = map[string]func(config.RuleConfig) (Ruler, error){
	"fake": func(cfg config.RuleConfig) (Ruler, error) {
		r := &BasicRuler{}
		r.SetName("fake")
		r.AddFilterFunc(func(e *event.Event) (bool, error) {
			return true, nil
		})
		return r, nil
	},
}

func RegisterRuleMaker(name string, f func(config.RuleConfig) (Ruler, error)) {
	ruleMakers[name] = f
}

func CreateRule(cfg config.RuleConfig) (rule Rule, err error) {
	if rule, err = InitRule(cfg); err != nil {
		return
	}
	if cfg.Name == "" {
		rule.Log.Info("构建fake规则......")
		cfg.Name = "fake"
	}
	f, ok := ruleMakers[cfg.Name]
	if !ok {
		err = errors.Errorf(LOAD_ERR2, cfg.Name)
		return
	}
	for i := 0; i < rule.rulerNum; i++ {
		var ruler Ruler
		if ruler, err = f(cfg); err != nil {
			return
		}
		if ruler == nil {
			err = errors.Errorf(LOAD_ERR3, cfg.Name)
			return
		}
		ruler.SetNumber(i)
		ruler.SetLogger(rule.Log)
		if err = ruler.LoadConfig(cfg); err != nil {
			return
		}
		if rule.IsAggre() {
			ruler.SetAggregator(rule.aggregator)
		}
		rule.SetRuler(ruler)
	}
	return
}

func InitRule(cfg config.RuleConfig) (rule Rule, err error) {
	rule = Rule{}
	rule.eventChannel = make(chan event.Event, cfg.GetBufferNum())
	rule.errHr = cobraErrors.MakeErrHandler(cfg.ErrSenderCfg.Parse(), cfg.GetBufferNum())
	rule.isRulerClose = make(chan bool, 1)
	rule.Log, err = cfg.LogCfg.GetLogger()
	rule.rulerNum = cfg.Worker()
	rule.aggregator = cfg.InitAggregator()
	rule.name = cfg.Name
	return
}

func (this *Rule) SetName(name string) {
	this.name = name
}

func (this *Rule) GetName() string {
	return this.name
}

func (this *Rule) IsAggre() bool {
	return this.aggregator != nil
}

func (this *Rule) SetRuler(r Ruler) {
	if this.ruler == nil {
		this.ruler = make([]Ruler, 0)
	}
	this.ruler = append(this.ruler, r)
}

func (this *Rule) Push(e event.Event) {
	if this.closed {
		this.Log.Errorf("%s规则事件池已经关闭，不能放入事件", this.name)
		return
	}
	this.eventChannel <- e
}

func (this *Rule) Close() error {
	if this.closed {
		return nil
	}
	close(this.eventChannel)
	this.closed = true
	<-this.isRulerClose
	var err error
	for _, ruler := range this.ruler {
		e := ruler.Close()
		if e != nil {
			if err != nil {
				err = errors.Errorf("%s->%s", err, e)
			} else {
				err = e
			}
		}
	}
	this.errHr.Close()
	this.Log.Infof("%s规则的错误处理器关闭", this.name)
	return err
}

func (this *Rule) Reset() error {
	var err error
	this.eventChannel = make(chan event.Event, cap(this.eventChannel))
	if this.errHr != nil {
		this.errHr.Reset()
	}
	if this.IsAggre() {
		this.aggregator.Reset()
	}
	for _, r := range this.ruler {
		if err = r.Reset(); err != nil {
			break
		}
	}
	if err != nil {
		return err
	}
	this.isReady = false
	this.closed = false
	return nil
}

func (this *Rule) Start() {
	if this.isReady {
		return
	}
	this.Log.Infof("%s规则的事件池开启...", this.name)
	this.isReady = true
	go this.errHr.Send()
	this.Log.Infof("%s规则的错误处理器开启", this.name)
	var wg sync.WaitGroup
	for _, ruler := range this.ruler {
		ruler.Start()
		wg.Add(1)
		go func(r Ruler) {
			for {
				e, isOpen := <-this.eventChannel
				if !isOpen {
					break
				}
				this.Log.Debugf("Rule%d: %s", r.GetNumber(), HEADER)
				this.Log.Debugf("Rule%d: %s规则发现有事件需要处理:%s", r.GetNumber(), this.name, e.String())
				if err := r.HandleEvent(e); err != nil {
					go this.errHr.Push(err)
				}
				this.Log.Debugf("Rule%d: 处理完毕", r.GetNumber())
			}
			wg.Done()
		}(ruler)
	}
	wg.Wait()
	this.Log.Infof("%s规则的事件池关闭", this.name)
	if this.IsAggre() {
		this.aggregator.Stop()
		this.Log.Infof("%s规则关闭聚合器", this.name)
	}
	this.isRulerClose <- true
}

func (this *Rule) IsClosed() bool {
	return this.closed
}

func (this *Rule) RulerNum() int {
	return this.rulerNum
}

func (this *Rule) ActiveRulerNum() int {
	var num int
	for _, r := range this.ruler {
		if !r.IsClosed() {
			num = num + 1
		}
	}
	return num
}
