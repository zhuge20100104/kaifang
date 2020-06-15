package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/zhuge20100104/kaifang/dataclean/utils"
)

func main() {
	file, err := os.Open("D:/result.txt")
	// 打开文件失败就 panic掉，不再执行
	utils.PanicErrorHand(err, "os.Open")
	defer file.Close()

	// 打开好结果文件对象
	goodFile, err := os.OpenFile("D:/result/kaifang_good.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	utils.PanicErrorHand(err, "os.OpenFile: GoodFile")
	defer goodFile.Close()

	// 打开坏结果文件对象
	badFile, err := os.OpenFile("D:/result/kaifang_bad.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	utils.PanicErrorHand(err, "os.OpenFile: BadFile")
	defer badFile.Close()

	// 构造缓冲读取器
	reader := bufio.NewReaderSize(file, 10*1024)

	for {
		lineBytes, _, err := reader.ReadLine()
		// 读取完毕，退出
		if err == io.EOF {
			break
		}

		// 处理除了io.EOF之外的其他错误
		if err != nil {
			utils.DefErrorHand(err, "reader.ReadLine")
		}

		// 默认就是UTF-8，不用Converting
		utfStr := string(lineBytes)
		fields := strings.Split(utfStr, ",")
		// 身份证在第四列，身份证长度为18
		if len(fields) > 4 && len(fields[4]) == 18 {
			// 写入Good File
			goodFile.WriteString(utfStr + "\n")
		} else {
			// 写入Bad File
			badFile.WriteString(utfStr + "\n")
		}
	}
}
