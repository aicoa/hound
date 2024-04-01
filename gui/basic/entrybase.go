package basic

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

type EntryBase struct {
	widget.Entry
}

func NewEntryBase() *EntryBase {
	eb := &EntryBase{}
	eb.ExtendBaseWidget(eb)
	return eb
}

// 重命名为createBaseMenuItem，以避免与子类方法冲突
func (eb *EntryBase) createBaseMenuItem(label string, icon fyne.Resource, action func()) *fyne.MenuItem {
	return &fyne.MenuItem{Label: label, Icon: icon, Action: action}
}

// 提供一个创建通用菜单项的方法
func (eb *EntryBase) createCommonMenuItems() []*fyne.MenuItem {
	return []*fyne.MenuItem{
		eb.createBaseMenuItem("复制", theme.ContentCopyIcon(), func() {
			clipboard.WriteAll(eb.Text)
		}),
		eb.createBaseMenuItem("剪切", theme.ContentCutIcon(), func() {
			clipboard.WriteAll(eb.Text)
			eb.SetText("")
		}),
		eb.createBaseMenuItem("清空", theme.ContentClearIcon(), func() {
			eb.SetText("")
		}),
	}
}
