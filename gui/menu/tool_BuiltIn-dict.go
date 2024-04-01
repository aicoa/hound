package menu

import (
	"hound/common"
	"hound/gui/basic"
	"hound/gui/mytheme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BuiltIn_Dict() {
	showDict := widget.NewMultiLineEntry()
	pro := &widget.Select{Options: []string{"ftp", "ssh", "telnet", "smb", "mssql", "oracle", "mysql", "postgresql", "vnc", "redis", "weblogic", "致远OA"}, Selected: "ftp"}
	file := &widget.Select{Options: []string{"username", "password"}}
	// 选择协议
	pro.OnChanged = func(s string) {
		showDict.SetText("")
		if file.Selected == "username" {
			for _, v := range common.Userdict[s] {
				showDict.Text += v + "\n"
				showDict.Refresh()
			}
		} else {
			for _, v := range common.Passwords {
				showDict.Text += v + "\n"
				showDict.Refresh()
			}
		}
	}
	file.OnChanged = func(s string) {
		showDict.SetText("")
		if pro.Selected == "username" {
			for _, v := range common.Userdict[s] {
				showDict.Text += v + "\n"
				showDict.Refresh()
			}
		} else {
			for _, v := range common.Passwords {
				showDict.Text += v + "\n"
				showDict.Refresh()
			}
		}
	}
	file.SetSelectedIndex(0)
	c := container.NewBorder(container.NewGridWithColumns(2, pro, file), nil, nil, nil, showDict)
	basic.ShowCustomDialog(mytheme.DictIcon(), "内置字典查看器", "内置字典不支持拓展", c, nil, fyne.NewSize(500, 500))
}
