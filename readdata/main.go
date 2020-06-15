package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/zhuge20100104/kaifang/readdata/utils"
)

func main() {
	file, err := os.Open("D:/result.txt")
	// 打开文件失败就 panic掉，不再执行
	utils.PanicErrorHand(err, "os.Open")
	defer file.Close()

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
		fmt.Println(utfStr)
	}
}
