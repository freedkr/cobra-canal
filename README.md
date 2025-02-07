# Cobra-Canal
binlog监控框架

Cobra-Canal根据用户自定的规则，对事件进行过滤，同时根据用户自定的消费器进行消费。框架支持多规则，对于一个事件，可以
应用于多个规则，同时拥有聚合功能，可以对同一id的多次事件进行聚合处理。在运行中出现的错误，可以自定义错误处理机制。

# Table of Contents

- [Overview](#overview)
- [Concepts](#concepts)

# Overview
Cobra-Canal是一个binlog监控框架，基于[canal](https://github.com/go-mysql-org/go-mysql)项目进行日志获取与解析。

Cobra-Canal可以使日志监控更加定制化，满足各种同步或服务发现场景。Cobra-Canal支持如下功能:
* 自定义过滤规则
* 自定义消费规则
* 自定义错误处理
* 事件聚合
* 多过滤多消费运行
* 丰富的命令

# Concepts
由于mysql可以批量操作，因此Cobra-Canal会把binlog日志变为事件组，一个事件代表一个变更。
事件会被放入每个过滤规则(ruler)对应的事件池中，ruler从事件池中获取事件进行判断，对于不符合条件的事件，直接丢弃。
如果ruler定义了聚合逻辑，会把事件根据定义的逻辑进行聚合，放入聚合器，在聚合时间到达后，聚合器把事件组放入消费池，
否则ruler会直接把事件放入消费池。消费器会从消费池获取事件(组)进行消费。

过滤规则(ruler)以及消费器(consumer)需要用户自定义，框架自带了一个假ruler和假consumer。假ruler没有任何过滤逻辑，假
consumer会把事件打印到日志，可以参见[examples/00-fake](https://github.com/freedkr/cobra-canal/tree/examples/examples/00-fake)

ruler和consumer的编写，可以参见[examples/01-print](https://github.com/freedkr/cobra-canal/tree/examples/examples/01-print)
