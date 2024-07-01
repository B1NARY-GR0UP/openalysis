# 如何使用 API？

## 获取库

```shell
go get -u github.com/B1NARY-GR0UP/openalysis
```

## 示例

```go
package main

import (
	"context"

	"github.com/B1NARY-GR0UP/openalysis/api"
)

func main() {
	err := api.ReadInConfig("config.yaml")
	if err != nil {
		panic("failed to read config file")
	}
	api.Init()
	api.Start(context.Background())
	// api.Restart(context.Background())
}
```