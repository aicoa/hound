package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ShowCustomeDiglog(icon fyne.Resource, title, confirm string, content fyne.CanvasObject, OnTapped func(), size fyne.Size) {
	ok := &widget.Button{Text: confirm, Icon: theme.ConfirmIcon(), Importance: widget.HighImportance, OnTapped: OnTapped}
	if OnTapped == nil {
		ok.Hide()
	}
	close := widget.NewButtonWithIcon("", theme.WindowCloseIcon(), nil)
	close.Importance = widget.LowImportance
	p := widget.NewModalPopUp(container.NewBorder(container.NewBorder(nil, nil, widget.NewIcon(icon), close, NewCenterLabel(title)),
		container.NewHBox(layout.NewSpacer(), ok, layout.NewSpacer()), nil, nil, c), global.win.Canvas())
	close.OnTapped = func() {
		p.Hide()
	}
	p.Resize(size)
	p.Show()
}
