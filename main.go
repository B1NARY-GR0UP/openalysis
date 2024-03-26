package main

import (
	"github.com/B1NARY-GR0UP/openalysis/cmd"
)

// TODO: repo table 每次更新都需要插入新的 item，但是 NodeID 是一样的，ID 是自增的，结合 CreatedAt 来绘制时间序列图
func main() {
	cmd.Execute()
}
