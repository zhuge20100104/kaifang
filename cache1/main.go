package main

import (
	"fmt"
	"os"

	"github.com/zhuge20100104/kaifang/cache1/db"

	m "github.com/zhuge20100104/kaifang/cache1/models"
	"github.com/zhuge20100104/kaifang/cache1/utils"
)

const (
	// CacheLen 缓存长度
	CacheLen = 2
)

var (
	// 开房结果缓存
	kfMap map[string]*m.QueryResult
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
	// 初始化开房缓存对象
	kfMap = make(map[string]*m.QueryResult, 0)

	for {
		var name string
		fmt.Printf("请输入要查询的开房者姓名:")
		fmt.Scanln(&name)

		// 用户选择退出
		if name == "exit" {
			fmt.Println("您选择了退出程序!")
			break
		}

		// 用户选择查询缓存
		if name == "cache" {
			fmt.Printf("共缓存了 %v 条数据\n", len(kfMap))
			for k := range kfMap {
				fmt.Println(k)
			}
			continue
		}

		// 先看看内存中是否有结果，有结果就直接展示
		if qr, ok := kfMap[name]; ok {
			fmt.Printf("共查询到 %d 条结果\n", len(qr.Value))
			// 查询次数加1
			qr.Count++
			fmt.Println(qr.Value)
			continue
		}

		// 没有结果，到数据库中查询，并缓存
		qRes, err := db.Db.QueryResultByName(name)
		if err != nil {
			continue
		}
		fmt.Println(qRes.Value)
		// 加入缓存
		kfMap[name] = qRes

		// 大于了最大缓存长度，需要更新缓存
		if len(kfMap) > CacheLen {
			delKey := utils.UpdateCache(kfMap)
			fmt.Printf("%v被从缓存中淘汰!\n", delKey)
		}
	}
}
