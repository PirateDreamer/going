# going

Go web development extension package, simple and easy to use

support golang version 1.22.2+

## install

```
go get github.com/PirateDreamer/going
```

## content

ginx: Based on gin support for HTTP processing

gormx: Mysql connect and custom log

zlog: Provide Zap with support for log splitting

xerr: Custom err

stl: Golang container

config: Support local yml and etcd

gredis: Redis connect

## use

```go
package main

import (
   "going-demo/api"

   "github.com/PirateDreamer/going"
)

func main() {
   router := going.InitService()
   // api service
   api.InitApi()
   going.GranceRun(router)
}
```

## Example

[GitHub - PirateDreamer/going-demo: going demo](https://github.com/PirateDreamer/going-demo)