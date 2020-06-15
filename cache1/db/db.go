package db

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	// Mysql 数据库驱动
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	m "github.com/zhuge20100104/kaifang/cache1/models"
	"github.com/zhuge20100104/kaifang/cache1/utils"
)

// Db 外部可访问的全局数据库连接对象
var Db *DB

// DB DB实现对象，主要实现对数据库的增删改查，建表等
type DB struct {
	Db *sqlx.DB
}

const (
	// CreateTableSQL 建表语句
	CreateTableSQL = `create table if not exists kfperson(
		id int primary key auto_increment,
		name varchar(100),
		idcard char(18)
	);`
)

func init() {
	// 新建一个DB对象
	Db = &DB{}

	// 初始化DB对象的ChanData对象
	db, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/kaifang")
	utils.PanicErrorHand(err, "sqlx.Open DB")

	// 初始化DB对象的Db字段，代表一个sqlx.DB对象
	Db.Db = db
}

func (imp *DB) insertKFPerson(chanData chan *m.KFPerson, exitChan chan bool) {
	defer func() {
		exitChan <- true
	}()

	for kfPerson := range chanData {
		for {
			result, err := imp.Db.Exec("insert into kfperson(name, idcard) values(?,?)",
				kfPerson.Name, kfPerson.IDCard)
			utils.DefErrorHand(err, "imp.Db.Exec Insert into DB")
			// 插入失败，先处理 error
			if err != nil {
				<-time.After(5 * time.Second)
			} else {
				// 验证插入成功
				if n, e := result.RowsAffected(); e == nil && n > 0 {
					fmt.Printf("插入%s成功\n", kfPerson.Name)
					break
				}
			}
		}
	}
}

// InitDB 初始化数据库
func (imp *DB) InitDB() {
	// 打开文本大数据文件
	file, err := os.Open("D:/result/kaifang_good.txt")
	utils.PanicErrorHand(err, "os.Open Good File")
	defer file.Close()

	// 创建数据库表对象
	_, err = imp.Db.Exec(CreateTableSQL)
	utils.PanicErrorHand(err, "imp.Db.Exec Create Table")

	chanData := make(chan *m.KFPerson, 1000000)
	exitChan := make(chan bool)
	// 开启100个携程，写入数据库
	for i := 0; i < 100; i++ {
		go imp.insertKFPerson(chanData, exitChan)
	}

	reader := bufio.NewReader(file)
	for {
		lineBytes, _, err := reader.ReadLine()

		// 说明读完了，退出
		if err == io.EOF {
			close(chanData)
			break
		}
		utils.DefErrorHand(err, "reader.ReadLine 逐行读取Good File")
		lineStr := string(lineBytes)
		fields := strings.Split(lineStr, ",")
		name, idcard := fields[0], fields[4]
		name = strings.TrimSpace(name)
		if len(strings.Split(name, "")) > 20 {
			fmt.Printf("%s 怪物出现了\n", name)
			continue
		}
		kfPerson := m.KFPerson{
			Name:   name,
			IDCard: idcard,
		}
		chanData <- &kfPerson
	}

	// 阻塞等待插入数据库完成
	<-exitChan
	fmt.Println("插入数据库完成!")
}

// QueryResultByName 通过名字查询开房人信息的函数
func (imp *DB) QueryResultByName(name string) (*m.QueryResult, error) {
	kfs := make([]m.KFPerson, 0)
	err := imp.Db.Select(&kfs, "select id, name, idcard from kfperson where name like ?", name)
	if err != nil {
		utils.DefErrorHand(err, "imp.Db.Select")
		return nil, err
	}
	qRes := &m.QueryResult{
		Value:     kfs,
		CacheTime: time.Now().UnixNano(),
		Count:     1,
	}
	return qRes, nil
}

// CloseDB 关闭数据库操作
func (imp *DB) CloseDB() {
	imp.Db.Close()
}
