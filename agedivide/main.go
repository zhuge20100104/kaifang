package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	m "github.com/zhuge20100104/kaifang/agedivide/models"
	"github.com/zhuge20100104/kaifang/agedivide/utils"
)

// 年龄划分
func main() {
	// 打开好的开房记录文件
	file, err := os.Open("D:/result/kaifang_good.txt")
	// 打开文件失败就 panic掉，不再执行
	utils.PanicErrorHand(err, "os.Open Good File")
	defer file.Close()

	aMap := make(map[string]*m.Ager)

	// 构造一堆年代对象
	for i := 190; i < 202; i++ {
		ager := m.Ager{
			Decade: strconv.Itoa(i),
		}
		fileName := fmt.Sprintf("D:/agers/%vx.txt", i)
		pFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		utils.PanicErrorHand(err, "os.OpenFile Agers File")
		ager.File = pFile
		ager.ChanData = make(chan string)
		aMap[ager.Decade] = &ager
	}

	// 同步等待组对象
	var wg sync.WaitGroup

	// 开启多个携程，分别写入ager年代文件
	for _, ager := range aMap {
		wg.Add(1)
		go utils.WriteFile(ager, &wg)
	}

	reader := bufio.NewReaderSize(file, 10*1024)

	// 开始主体部分，读取文件Good File.txt
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			// 读完文件时，关闭所有channel
			for _, province := range aMap {
				close(province.ChanData)
			}
			break
		}

		// 处理除 io.EOF以外的其他错误
		if err != nil {
			utils.DefErrorHand(err, "reader.ReadLine Good File")
		}

		lineStr := string(lineBytes)
		fields := strings.Split(lineStr, ",")
		decade := fields[4][6:9]
		if _, ok := aMap[decade]; ok {
			aMap[decade].ChanData <- (lineStr + "\n")
		}

	}

	wg.Wait()

}
