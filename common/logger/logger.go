package logger

import (
	"fmt"
	"os"
	"time"
)

const defaultLogFilePath = "./log/"

// IFO相关的记录不影响结果的日志
func Info(message interface{}) {
	log("INF", fmt.Sprintf("%v", message))
}

// ERROR会记录退出的日志
func Error(message interface{}) {
	log("[ERRO]", fmt.Sprintf("%v", message))
}

// Debug记录会影响结果的日志
func Debug(message interface{}) {
	log("[DEBUG]", fmt.Sprintf("%v", message))
}

// log写日志
func log(level string, message string) {
	if _, err := os.Stat(defaultLogFilePath); os.IsNotExist(err) {
		os.Mkdir(defaultLogFilePath, os.ModePerm)
	}
	f, _ := os.OpenFile(defaultLogFilePath+time.Now().Format("2006-01-02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	f.WriteString(fmt.Sprintf("%v %v %v\n", time.Now().Format("2006-01-02 15:04:05"), level, message))
}
