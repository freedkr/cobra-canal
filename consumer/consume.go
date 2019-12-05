package consumer

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/always-waiting/cobra-canal/config"
	cobraErrors "github.com/always-waiting/cobra-canal/errors"
	"github.com/always-waiting/cobra-canal/event"
	"github.com/juju/errors"
	"github.com/siddontang/go-log/log"
	"strings"
)

const (
	HEADER     = ">>>>>>>>%s<<<<<<<<\n"
	DONE       = "<<<<<<<<%s>>>>>>>>\n"
	EVENT_LINE = "################"
)

type Consume struct {
	consumer        Consumer
	eventsChan      chan []event.Event
	closed          bool
	isReady         bool
	isConsumerClose chan bool
	errHr           *cobraErrors.ErrHandler
	Log             *log.Logger
}

func InitConsume(cfg *config.ConsumerConfig) (Consume, error) {
	var err error
	consume := Consume{}
	consume.eventsChan = make(chan []event.Event, cfg.GetBufferNum())
	consume.isConsumerClose = make(chan bool, 1)
	consume.errHr = cobraErrors.MakeErrHandler(cfg.ErrSenderCfg.Parse(), cfg.GetBufferNum())
	consume.Log, err = cfg.LogCfg.GetLogger()
	return consume, err
}

func (this *Consume) SetTransferFunc(f func([]event.Event) (interface{}, error)) {
	this.consumer.SetTransferFunc(f)
}

func (this *Consume) SetConsumer(consumer Consumer) {
	this.consumer = consumer
}

func (this *Consume) GetName() string {
	return this.consumer.GetName()
}

func (this *Consume) Push(input []event.Event) {
	if this.closed {
		this.Log.Errorf("%s消费池已经关闭，无法放入事件包", this.GetName())
		return
	}
	this.eventsChan <- input
}

func (this *Consume) Close() error {
	if this.closed {
		return nil
	}
	close(this.eventsChan)
	this.closed = true
	<-this.isConsumerClose
	err := this.consumer.Close()
	this.Log.Infof("%s消费器关闭", this.GetName())
	this.errHr.Close()
	this.Log.Infof("%s消费错误处理器关闭", this.GetName())
	return err
}

func (this *Consume) Start() {
	this.Log.Infof("%s消费池开启...", this.GetName())
	if this.isReady {
		return
	}
	this.isReady = true
	go this.errHr.Send()
	for {
		input, isOpen := <-this.eventsChan
		if !isOpen {
			break
		}
		this.Log.Debugf(HEADER, "消费开始")
		this.Log.Debug("发现如下事件包:")
		for _, e := range input {
			this.Log.Debug(EVENT_LINE)
			this.Log.Debug(e.String())
		}
		this.Log.Debugf("消费器%s转换事件包", this.GetName())
		data, err := this.consumer.Transfer(input)
		if err != nil {
			go this.errHr.Push(this.modifyErr(err, input))
			continue
		}
		this.Log.Debugf("转换后的信息为:%#v\n", data)
		this.Log.Debugf("消费器%s消费事件包", this.GetName())
		if err = this.consumer.Solve(data); err != nil {
			go this.errHr.Push(this.modifyErr(err, input))
		}
		this.Log.Debugf("消费完毕")
	}
	this.Log.Infof("%s消费池关闭", this.GetName())
	this.isConsumerClose <- true
}

func (this *Consume) modifyErr(err error, input []event.Event) (retErr error) {
	if err == nil {
		return
	}
	docTmp := heredoc.Doc(`<PRE>
	事件包为:<br>
	%s<br>
	消费器错误信息:<br>
	%s<br>
	</PRE>`)
	inputStr := make([]string, 0)
	for _, data := range input {
		inputStr = append(inputStr, data.String())
	}
	retErr = errors.Errorf(docTmp, strings.Join(inputStr, "<br>###########<br>"), err.Error())
	return
}