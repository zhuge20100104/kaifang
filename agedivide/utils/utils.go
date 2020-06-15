package utils

import (
	"fmt"
	"sync"

	"github.com/axgle/mahonia"
	m "github.com/zhuge20100104/kaifang/agedivide/models"
)

// handleErrorFunc 错误处理函数闭包
// shouldPanic 是否引发恐慌错误
// err 传入的错误信息
// msg 传入的msg信息
func handleErrorFunc(shouldPanic bool) func(error, string) {
	return func(err error, msg string) {
		if err != nil {
			fmt.Printf("[%s] --- %v\n", msg, err)
		}

		if err != nil && shouldPanic {
			panic(err)
		}
	}
}

// PanicErrorHand 处理错误时同时panic
var PanicErrorHand = handleErrorFunc(true)

// DefErrorHand 处理错误时不panic
var DefErrorHand = handleErrorFunc(false)

// ConvertEncoding 将编码从 srcEncoding 转换到 dstEncoding的函数
func ConvertEncoding(srcStr string, srcEncoding, dstEncoding string) (dstStr string, err error) {
	srcDecoder := mahonia.NewDecoder(srcEncoding)
	dstDecoder := mahonia.NewDecoder(dstEncoding)
	utfStr := srcDecoder.ConvertString(srcStr)
	_, dstBytes, err := dstDecoder.Translate([]byte(utfStr), true)
	if err != nil {
		return
	}

	dstStr = string(dstBytes)
	return
}

// WriteFile 写入单个Province文件的函数
func WriteFile(ager *m.Ager, wg *sync.WaitGroup) {
	// 最终同步等待组计数值减1
	defer wg.Done()
	for lineStr := range ager.ChanData {
		_, err := ager.File.WriteString(lineStr)
		DefErrorHand(err, "ager.File.WriteString")
	}
}
