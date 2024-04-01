package basic

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// true为PlaceHolder,False为Text
func NewMultiLineEntry(initialValue string, isPlaceHolder bool) *widget.Entry {
	e := widget.NewMultiLineEntry()
	e.Wrapping = fyne.TextWrapBreak

	if isPlaceHolder {
		e.PlaceHolder = initialValue
	} else {
		e.SetText(initialValue)
	}

	return e
}

// 对齐元素
func NewFormItem(text string, object fyne.CanvasObject) *fyne.Container {
	return container.NewBorder(nil, nil, widget.NewLabel(text), nil, object)
}

// 基础框架
func Frame(content fyne.CanvasObject) fyne.CanvasObject {
	top := canvas.NewLine(theme.FocusColor())
	buttom := canvas.NewLine(theme.FocusColor())
	top.Resize(fyne.NewSize(1, 0))
	buttom.Resize(fyne.NewSize(1, 0))
	return container.NewBorder(top, buttom, nil, nil, content)
}

func (fe *ForwordEntry) createMenuItem(label string, icon fyne.Resource, action func()) *fyne.MenuItem {
	return &fyne.MenuItem{Label: label, Icon: icon, Action: action}
}
