# going

Go web development extension package, simple and easy to use.

All functions support configuration using configuration files.

support golang version 1.22.2+

## install

```bash
go get github.com/PirateDreamer/going
```

## content

ginx: Based on gin support for HTTP processing.

gormx: Mysql connect and custom log.

zlog: Provide Zap with support for log splitting.

xerr: Custom err.

stl: Golang container.

config: Support local yml and etcd.

gredis: Redis connect.

## use

Create project

```bash
mkdir demo
cd demo
go mod init demo
```

Create configuration files and go main function files

```bash
touch config.yml
touch main.go
```

The content of the config.yml file

```yml
server:
  addr: 0.0.0.0:8080
```

Main.go Content

```go
package main

import (
   "github.com/PirateDreamer/going"
)

func main() {
   router := going.InitService()
    
   going.GranceRun(router)
}
```

Run the project

```bash
go run main.go
```



## Example

[GitHub - PirateDreamer/going-demo: going demo](https://github.com/PirateDreamer/going-demo)