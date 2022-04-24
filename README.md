# sidergo

```
      _      _
     (_)    | |
 ___  _   __| |  ___  _ __   __ _   ___
/ __|| | / _` | / _ \| '__| / _` | / _ \
\__ \| || (_| ||  __/| |   | (_| || (_) |
|___/|_| \__,_| \___||_|    \__, | \___/
                             __/ |
                            |___/

```

![github action](https://github.com/chenjiayao/sidergo/actions/workflows/master.yml/badge.svg)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-brightgreen.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/chenjiayao/sidergo.svg)](https://pkg.go.dev/github.com/chenjiayao/sidergo)

使用 Go 实现 redis-server 部分功能，**该项目不是一个用于生产环境的产品**，旨在通过该项目学习 Go 开发。sidergo 配有系列教程可以作为参考：[sidergo 系列教程](https://sidergo.jaychen.fun/)


## 🔜 快速开始

1. 执行 `go run cmd/main.go`
![](https://raw.githubusercontent.com/chenjiayao/sidergo-posts/master/docs/images/20220424173207.png)

1. 使用 redis-cli 连接到服务端：`redis-cli -p 3101`
![](https://raw.githubusercontent.com/chenjiayao/sidergo-posts/master/docs/images/20220424173309.png)


## 🧑‍💻 已实现功能

- [x] string、set、list、hash、zset 等数据结构
- [x] multi 事务，支持 watch、discard 等操作
- [x] 实现并发安全的 map 作为 redis db 存储数据
- [x] 实现 list 中 blpush、lpop 等阻塞命令
- [x] AOF 持久化
- [x] 支持 key 自动过期
- [x] 实现 unboundChan 用于 AOF channel
- [x] msetnx、incr 等命令原子操作实现
- [x] 核心逻辑的单元测试
- [x] skipList 数据结构实现，用于 redis zset 数据结构的底层存储
- [x] 集群模式


## 🤯 Benchmark

```
SET: 121951.22 requests per second, p50=0.047 msec
GET: 178571.42 requests per second, p50=0.039 msec
INCR: 169491.53 requests per second, p50=0.039 msec
LPUSH: 169491.53 requests per second, p50=0.039 msec
RPUSH: 169491.53 requests per second, p50=0.039 msec
LPOP: 172413.80 requests per second, p50=0.039 msec
RPOP: 175438.59 requests per second, p50=0.039 msec
SADD: 172413.80 requests per second, p50=0.039 msec
HSET: 175438.59 requests per second, p50=0.039 msec
SPOP: 58823.53 requests per second, p50=0.047 msec
LPUSH (needed to benchmark LRANGE): 169491.53 requests per second, p50=0.039 msec
LRANGE_100 (first 100 elements): 56497.18 requests per second, p50=0.079 msec
LRANGE_300 (first 300 elements): 23094.69 requests per second, p50=0.175 msec
LRANGE_500 (first 500 elements): 14992.50 requests per second, p50=0.295 msec
LRANGE_600 (first 600 elements): 9643.20 requests per second, p50=0.911 msec
```

## License

This project is licensed under the [GPL license](https://github.com/chenjiayao/sidergo/blob/master/LICENSE).