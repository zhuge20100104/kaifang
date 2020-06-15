package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	m "github.com/zhuge20100104/kaifang/provicedivide/models"
	"github.com/zhuge20100104/kaifang/provicedivide/utils"
)

// 省份划分
func main() {
	// 打开好的开房记录文件
	file, err := os.Open("D:/result/kaifang_good.txt")
	// 打开文件失败就 panic掉，不再执行
	utils.PanicErrorHand(err, "os.Open Good File")
	defer file.Close()

	// 省份和Code对应数组
	// Code在最后两位
	ps := []string{"北京市11", "天津市12", "河北省13",
		"山西省14", "内蒙古自治区15", "辽宁省21",
		"吉林省22", "黑龙江省23", "上海市31",
		"江苏省32", "浙江省33", "安徽省34",
		"福建省35", "江西省36", "山东省37",
		"河南省41", "湖北省42", "湖南省43",
		"广东省44", "广西壮族自治区45", "海南省46",
		"重庆市50", "四川省51", "贵州省52",
		"云南省53", "西藏自治区54", "陕西省61",
		"甘肃省62", "青海省63", "宁夏回族自治区64",
		"新疆维吾尔自治区65", "台湾省71", "香港特别行政区81",
		"澳门特别行政区82"}

	pm := make(map[string]*m.Province)

	// 初始化Province Map对象
	for _, p := range ps {
		name := p[:len(p)-2]
		id := p[len(p)-2:]
		province := m.Province{ID: id, Name: name}
		pm[id] = &province
	}

	// 打开34个文件，用于写入各个省份的数据
	for _, province := range pm {
		fileName := fmt.Sprintf("D:/province/%v.txt", province.Name)
		pFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		utils.PanicErrorHand(err, "os.OpenFile: Open Province File")
		province.File = pFile
		defer pFile.Close()
		province.ChanData = make(chan string)
	}

	// 同步等待组对象
	var wg sync.WaitGroup

	// 开启34个线程，分别写入34个省份文件
	for _, province := range pm {
		wg.Add(1)
		go utils.WriteFile(province, &wg)
	}

	reader := bufio.NewReaderSize(file, 10*1024)

	// 开始主体部分，读取文件Good File.txt
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			// 读完文件时，关闭所有channel
			for _, province := range pm {
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
		id := fields[4][:2]
		if _, ok := pm[id]; ok {
			pm[id].ChanData <- (lineStr + "\n")
		}

	}

	wg.Wait()

}
