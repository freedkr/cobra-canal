package rules

import (
	"fmt"
	"sync"

	"github.com/always-waiting/cobra-canal/config"
	"github.com/always-waiting/cobra-canal/consumer"
	"github.com/always-waiting/cobra-canal/event"

	"github.com/juju/errors"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
)

const (
	AGGREGATOR_HEADER = "Rule%d: -------->缓存键值为:%s<--------"
)

type BasicRuler struct {
	name             string
	number           int
	aggregator       config.Aggregatable
	consumers        map[string]*consumer.Consume
	DBClient         *client.Conn
	dbLock           sync.Mutex
	filter           FilterHandler
	isReady          bool
	closed           bool
	closeAggregation chan bool
	Log              *log.Logger
	transferFunc     map[string]func([]event.Event) (interface{}, error)
	hasLoadConfig    bool
}

func (this *BasicRuler) Debugf(tmp string, i ...interface{}) {
	nTmp := fmt.Sprintf("Rule%d: %s", this.number, tmp)
	this.Log.Debugf(nTmp, i...)
}

func (this *BasicRuler) Infof(tmp string, i ...interface{}) {
	nTmp := fmt.Sprintf("Rule%d: %s", this.number, tmp)
	this.Log.Infof(nTmp, i...)
}

func (this *BasicRuler) Errorf(tmp string, i ...interface{}) {
	nTmp := fmt.Sprintf("Rule%d: %s", this.number, tmp)
	this.Log.Errorf(nTmp, i...)
}

func (this *BasicRuler) Info(i string) {
	this.Log.Infof("Rule%d: %s", this.number, i)
}

func (this *BasicRuler) Debug(i string) {
	this.Log.Debugf("Rule%d: %s", this.number, i)
}

func (this *BasicRuler) Error(i string) {
	this.Log.Errorf("Rule%d: %s", this.number, i)
}

func (this *BasicRuler) SetLogger(l *log.Logger) {
	this.Log = l
}

func (this *BasicRuler) GetNumber() int {
	return this.number
}

func (this *BasicRuler) SetNumber(i int) {
	this.number = i
}

func (this *BasicRuler) SetAggregator(ag config.Aggregatable) {
	this.aggregator = ag
}

func (this *BasicRuler) LoadConfig(ruleCfg config.RuleConfig) (err error) {
	if this.hasLoadConfig {
		return
	}
	this.hasLoadConfig = true
	this.aggregator = ruleCfg.InitAggregator()
	this.closeAggregation = make(chan bool, 1)
	if len(ruleCfg.ReplySync) != 0 {
		this.filter.LoadReplySyncFilter(ruleCfg.ReplySync)
	}
	if ruleCfg.HasTableFilter() {
		this.filter.LoadTableFilter(ruleCfg.TableFilterCfg)
	}
	if len(ruleCfg.ConsumerCfg) == 0 {
		this.Log.Info("构建fake消费器......")
		ruleCfg.ConsumerCfg = append(ruleCfg.ConsumerCfg, &config.ConsumerConfig{})
	}
	this.consumers = make(map[string]*consumer.Consume)
	for _, consumerCfg := range ruleCfg.ConsumerCfg {
		if consume, err := consumer.CreateConsume(consumerCfg); err != nil {
			return err
		} else {
			if f, ok := this.transferFunc[consume.GetName()]; ok {
				consume.SetTransferFunc(f)
			}
			consume.SetRulerNum(this.number)
			this.consumers[consume.GetName()] = &consume
		}
	}
	if ruleCfg.MasterDBCfg != nil {
		if this.DBClient, err = client.Connect(
			ruleCfg.MasterDBCfg.Addr,
			ruleCfg.MasterDBCfg.User,
			ruleCfg.MasterDBCfg.Passwd,
			ruleCfg.MasterDBCfg.Db,
		); err != nil {
			return
		}
	}
	return
}

func (this *BasicRuler) DBExecute(cmd string, args ...interface{}) (*mysql.Result, error) {
	defer this.dbLock.Unlock()
	this.dbLock.Lock()
	return this.DBClient.Execute(cmd, args...)
}

func (this *BasicRuler) AddTransferFunc(name string, f func([]event.Event) (interface{}, error)) {
	if this.transferFunc == nil {
		this.transferFunc = make(map[string]func([]event.Event) (interface{}, error))
	}
	this.transferFunc[name] = f
}

/*
func (this *BasicRuler) SetConsumerTransferFunc(name string, f func([]event.Event) (interface{}, error)) (err error) {
	if _, ok := this.consumers[name]; !ok {
		err = errors.Errorf("未初始化的消费器%s", name)
		return
	}
	this.consumers[name].SetTransferFunc(f)
	return err
}
*/

func (this *BasicRuler) Start() {
	this.StartConsume()
	if this.IsAggre() {
		go this.StartAggregation()
	}
}

func (this *BasicRuler) Close() (err error) {
	if this.IsAggre() {
		err = this.CloseAggregation()
	}
	this.CloseConsume()
	if this.DBClient != nil {
		this.DBClient.Close()
	}
	return
}

func (this *BasicRuler) HandleEvent(e event.Event) (err error) {
	// 应用过滤规则
	flag, err := this.Filter(&e)
	if err != nil || !flag {
		err = this.ModifyErr(err)
		this.Log.Debugf("Rule%d: 事件跳过", this.number)
		return
	}
	// 消费事件
	if this.IsAggre() {
		err = this.Aggregate(e)
	} else {
		this.Push([]event.Event{e})
	}
	return
}

func (this *BasicRuler) AddFilterFunc(f func(*event.Event) (bool, error)) {
	this.filter.AddFilterFunc(f)
}

func (this *BasicRuler) SetName(name string) {
	this.name = name
}

func (this *BasicRuler) GetName() string {
	return this.name
}

func (this *BasicRuler) Filter(e *event.Event) (bool, error) {
	this.Log.Debugf("Rule%d: 根据规则进行过滤", this.number)
	return this.filter.Filter(e)
}

func (this *BasicRuler) IsAggre() bool {
	return this.aggregator != nil
}

func (this *BasicRuler) GetAggreRule(name string) *config.IdxRuleConfig {
	return this.aggregator.GetRule(name)
}

func (this *BasicRuler) DiffData(r *config.IdxRuleConfig, d1, d2 map[string]interface{}) (map[string]interface{}, error) {
	return this.aggregator.DiffData(r, d1, d2)
}

func (this *BasicRuler) GetAggreKey(e event.Event) (string, error) {
	return this.aggregator.GetIdxValue(e)
}

func (this *BasicRuler) StartConsume() {
	if this.isReady {
		return
	}
	this.isReady = true
	for _, consumer := range this.consumers {
		go consumer.Start()
	}
}

func (this *BasicRuler) CloseConsume() {
	if this.closed {
		return
	}
	this.closed = true
	var wg sync.WaitGroup
	for _, consume := range this.consumers {
		wg.Add(1)
		go func(c *consumer.Consume) {
			c.Close()
			wg.Done()
		}(consume)
	}
	wg.Wait()
}

func (this *BasicRuler) Push(events []event.Event) {
	var wg sync.WaitGroup
	for _, consume := range this.consumers {
		wg.Add(1)
		go func(c *consumer.Consume) {
			this.Log.Debugf("Rule%d: %s规则向%s消费池推送事件包", this.number, this.name, c.GetName())
			c.Push(events)
			wg.Done()
		}(consume)
	}
	wg.Wait()
}

func (this *BasicRuler) StartAggregation() {
	sendKey := this.aggregator.GetSendChan()
	for {
		key, isOpen := <-sendKey
		if !isOpen {
			break
		}
		events, err := this.aggregator.MoveEvents(key)
		if err != nil {
			continue
		}
		this.Log.Debugf("Rule%d: 聚合消费%s键的事件包", this.number, key)
		this.Push(events)
	}
	this.closeAggregation <- true
}

func (this *BasicRuler) CloseAggregation() error {
	//this.aggregator.Stop()
	<-this.closeAggregation
	return nil
}

func (this *BasicRuler) Aggregate(e event.Event) (err error) {
	key, err := this.aggregator.GetIdxValue(e)
	if err != nil {
		return
	}
	if key != "" {
		this.Log.Debugf(AGGREGATOR_HEADER, this.number, key)
		if this.aggregator.HasIdx(key) {
			err = this.aggregator.AppendEvent(key, e)
		} else {
			err = this.aggregator.CreateEvent(key, e)
		}
	}
	return
}

func (this *BasicRuler) ModifyErr(err error) (ret error) {
	if err != nil {
		str := fmt.Sprintf("Rule%d: 规则%s报错: %s", this.number, this.GetName(), err.Error())
		ret = errors.New(str)
	}
	return
}
