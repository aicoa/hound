/*
 * @Author: aicoa
 * @Date: 2024-02-05 23:03:34
 * @Last Modified by:   aicoa
 * @Last Modified time: 2024-02-05 23:03:34
 */
package custom

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const WaitTime = 5

var (
	Logtime int64
	Console *ConsoleLog
)

type ConsoleLog struct {
	widget.Label
}

func init() {
	Console = NewConsoleLog()
}
func NewConsoleLog() *ConsoleLog {
	c := &ConsoleLog{}
	c.Wrapping = fyne.TextWrapBreak
	c.ExtendBaseWidget(c)
	return c
}
func (this *ConsoleLog) Append(text string) {
	this.Text += text
	this.Refresh()
}
func ConsoleWindow() {
	ShowCustomeDiglog(theme.FileTextIcon(), "日志中心(资产收集/爆破状况)", "清空日志", Frame(container.NewVScroll(Console)), func() {
		Console.SetText("")
	}, fyne.NewSize(500, 700))
}
func LogProgress(currentNum int64, countNum int, errfinfo interface{}) {
	if (time.Now().Unix() - Logtime) > WaitTime {
		Console.Append(fmt.Sprintf("已完成 %v/%v %v\n", currentNum, countNum, errfinfo))
		Logtime = time.Now().Unix()
	}
}
