package models

import "os"

// Ager 按年代划分结构体
type Ager struct {
	Decade   string      //年代， 197x, 198x, 199x, 200x, 201x
	File     *os.File    // 每个年代对应的数据存储文件
	ChanData chan string // 每个年代对应的数据通道
}
