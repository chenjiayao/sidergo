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

使用 golang 实现 redis 

1. go run cmd/main.go，默认监听 3101 端口
2. redis-cli -p 3101


已经实现的功能：
1. TCP 层解析 redis 通信协议。
2. 数据结构 skipList 实现，用作 redis zset 数据结构的底层存储。
3. multi 事务支持，支持 watch 和 discard 等操作。
4. 使用 sync.Map 实现自旋锁保证 msetnx, incr 等命令的原子操作。
5. 实现并发安全 map 提高并发量。
6. string，set，hash，list，zset 等命令实现，兼容 redis server。
7. 实现 unboundChan 用于 AOF 写入。
8. 实现 list 的 blpush，blpop 等阻塞命令。

更多文档正在完善中。。。

