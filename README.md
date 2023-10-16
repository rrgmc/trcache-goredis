# trcache go-redis [![GoDoc](https://godoc.org/github.com/rrgmc/trcache-goredis?status.png)](https://godoc.org/github.com/rrgmc/trcache-goredis)

This is a [trcache](https://github.com/rrgmc/trcache) wrapper for [go-redis](https://github.com/redis/go-redis).

## Info

### go-redis library

| info        |          |
|-------------|----------|
| Generics    | No       |
| Key types   | `string` |
| Value types | `any`    |
| TTL         | Yes      |

### wrapper

| info              |                  |
|-------------------|------------------|
| Default codec     | `GOBCodec`       |
| Default key codec | `StringKeyCodec` |

## Installation

```shell
go get github.com/rrgmc/trcache-goredis
```
