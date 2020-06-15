package main

import (
	"fmt"
	"os"

	"github.com/zhuge20100104/kaifang/intodb/db"

	"github.com/zhuge20100104/kaifang/intodb/utils"
)

// init函数，只在模块加载时执行一次
func init() {
	if utils.CheckFileExists("D:/db/db.mark") {
		fmt.Println("数据库已经初始化!")
		return
	}
	fmt.Println("开始初始化数据库...")
	db.Db.InitDB()
	_, err := os.Create("D:/db/db.mark")
	utils.DefErrorHand(err, "os.Create mark File")
}

// 数据入库
func main() {
	// 最终一定要关闭数据库连接池
	defer db.Db.CloseDB()
}
