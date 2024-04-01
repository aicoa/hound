package global

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AppWidgets struct {
	ProgerssWebscan   *widget.Label
	PortscanProgress  *widget.Label
	PortScanTarget    *widget.Entry
	PortBurstResult   *widget.Entry
	PortBurstTarget   *widget.Entry
	SubdomainTarget   *widget.Entry
	WebScanTarget     *widget.Entry
	Win               fyne.Window // 全局窗口
	ThinkDict         *widget.Entry
	FscanText         *widget.Entry
	SubdomainText     *widget.Entry
	VulnerabilityText *widget.Entry
	UsernameText      *widget.Entry
	PasswordText      *widget.Entry
	DirDictText       *widget.Entry
}

var Widgets AppWidgets

// 初始化全局变量
func SetupWidgets() {
	Widgets.ProgerssWebscan = &widget.Label{}
	Widgets.PortscanProgress = &widget.Label{}
	Widgets.PortScanTarget = widget.NewEntry()
	Widgets.PortBurstResult = widget.NewEntry()
	Widgets.PortBurstTarget = widget.NewEntry()
	Widgets.SubdomainTarget = widget.NewEntry()
	Widgets.WebScanTarget = widget.NewEntry()
	Widgets.ThinkDict = widget.NewEntry()
	Widgets.FscanText = widget.NewEntry()
	Widgets.SubdomainText = widget.NewEntry()
	Widgets.VulnerabilityText = widget.NewEntry()
	Widgets.UsernameText = widget.NewEntry()
	Widgets.PasswordText = widget.NewEntry()
	Widgets.DirDictText = widget.NewEntry()
}

// 初始化拖拽方法
func SetupDragAndDrop(entries []*widget.Entry, win fyne.Window) {
	win.SetOnDropped(func(p fyne.Position, u []fyne.URI) {
		for _, entry := range entries {
			if isInEntryArea(entry, p) {
				entry.SetText(u[0].Path())
			}
		}
	})
}

func isInEntryArea(entry *widget.Entry, pos fyne.Position) bool {
	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(entry)
	size := entry.Size()
	return !entry.Disabled() &&
		entryPos.X <= pos.X && entryPos.X+size.Width >= pos.X &&
		entryPos.Y <= pos.Y && entryPos.Y+size.Height >= pos.Y
}
func RefreshTabs(tableapps ...*container.AppTabs) {
	for _, tabs := range tableapps {
		tabs.OnSelected = func(_ *container.TabItem) {
			fyne.CurrentApp().Settings().SetTheme(fyne.CurrentApp().Settings().Theme())
		}
	}
}
