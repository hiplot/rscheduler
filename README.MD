# rscheduler：hiplot任务调度优化解决方案

## 简介

  `rscheduler`本质为在`hiplotlib`之上搭建的一个R语言会话层，进行多会话的任务调度，提供会话管理功能，替代原有的`plumber`+`pm2`的调度方案

  工作流由原先的` 前端提交任务->后端处理->plmber执行` 改为`前端提交任务->后端处理->消息队列<-rscheduler拉取任务->hicli执行` ，引入的消息队列解耦了后端逻辑与任务调度与执行逻辑，使`hiplotlib`得以水平拓展，大大提高任务处理效率，并且通过拉模式主动从消息队列中消费任务，可根据自身状态动态调整拉取任务数量，一定程度上防止了CPU负载过高和内存溢出

## 代码实现

在`rscheduler`代码实现中，有三个核心概念，`scheduler`、`processor`、`task`

* `scheduler`为调度器，用于管理多个`processor`，执行创建、销毁、复用会话等操作，管理全局的内存和CPU使用率

* `processor`为处理器，对应一个R会话，实现会话创建、销毁、健康检查等功能，运行时唯一绑定一个`task`执行任务

* `task`为任务，解析消息队列拉取来数据并封装为`task`

除此之外还有`rslog`日志功能及`monitor`进程信息监控功能等

## 配置项

`rscheduler`的配置项从三个维度出发，分别对应上面的三个核心概念

* `scheduler`支持配置`maxCPU`允许创建新任务的最大CPU使用率和`maxMem`允许创建新任务的最大内存使用率，在达到最大内存使用率时，会自动触发`scheduler`的`gc`功能清除空闲的`processor`。另外还支持配置`maxBusyProcessor最大繁忙处理器数量`和`maxIdleProcessor最大空闲处理器数量`。

* `pocessor`支持配置`maxBusyProcessorMem繁忙处理器最大内存使用量`和`maxIdleProcessorMem空闲处理器最大内存使用量`，自动杀死超限任务，防止大型任务占用服务器所有资源，并且避免了会话长期复用累积的内存消耗问题

* `task`支持配置`taskTimeout`任务超时时间，防止任务异常阻塞等情况的发生

## 已实现功能

上述配置项中提到的功能均已实现，除此之外还支持的功能有

* `task`运行时CPU使用率和内存消耗统计
* 不同种类任务的会话分离（目前由于hiplot支持的任务种类较多，迁移困难，所以统一定义为common通用，`commonInit.R`由`plumber`执行函数迁移而来）
* `processor`健康检查（防止R会话异常退出无回调的情况）

## 待实现功能

* 监控看板
* `docker`一键部署
* 任务执行失败的重试机制
* 用户自定义是否保留会话的相关机制
* 多服务之间的负载均衡及状态共享

## 其他信息

详情请查阅 [语雀文档rscheduler](https://www.yuque.com/liang09255/kqaphq/ius7rerk4bw6c4md?singleDoc)，内容包括架构设计及代码设计的一些细节信息

