package models

import "os"

// Province 省份结构体
type Province struct {
	ID       string      // 省份ID对象
	Name     string      // 省份名称对象
	File     *os.File    // 当前省份写入文件
	ChanData chan string // 当前省份的chan数据
}
