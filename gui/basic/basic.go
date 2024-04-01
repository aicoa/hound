/*
 * @Author: aicoa
 * @Date: 2024-03-31 16:44:25
 * @Last Modified by: aicoa
 * @Last Modified time: 2024-03-31 17:11:13
 */

package basic

import (
	"fmt"
	"hound/common"
	"hound/gui/global"
	"hound/gui/mytheme"
	"hound/lib/util"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewCenterLabel(text string) *widget.Label {
	return &widget.Label{
		Text:       text,
		Alignment:  fyne.TextAlignCenter,
		TextStyle:  fyne.TextStyle{},          // 默认文本样式
		Truncation: fyne.TextTruncateEllipsis, // 超出部分用省略号显示
	}
}

// CheckListBox 是带复选框的列表组件。
type CheckListBox struct {
	widget.List
	Options  []string // 列表中的所有选项
	Selected []string // 已选中的选项
}

func NewCheckListBox(options []string) *CheckListBox {
	clb := &CheckListBox{Options: options}
	clb.Length = func() int {
		return len(clb.Options)
	}
	clb.CreateItem = func() fyne.CanvasObject {
		check := widget.NewCheck("", nil)
		checkRef := check // 捕获check的引用
		check.OnChanged = func(b bool) {
			if b {
				clb.Selected = append(clb.Selected, checkRef.Text)
			} else {
				clb.Selected = util.RemoveStringFromArray(clb.Selected, checkRef.Text)
			}
		}
		return check
	}
	clb.UpdateItem = func(id widget.ListItemID, o fyne.CanvasObject) {
		check := o.(*widget.Check)
		spl := strings.Split(clb.Options[id], "\\") // 在这里声明并初始化 spl
		check.Text = spl[len(spl)-1]                // 现在可以直接使用 spl
		check.Refresh()
	}
	clb.ExtendBaseWidget(clb)
	return clb
}

// 全局变量定义
const WaitTime = 5

var (
	LogTime int64
	Console *ConsoleLog // 全局日志记录器
)

type ConsoleLog struct {
	widget.Label
}

func init() {
	Console = NewConsoleLog() // 初始化日志记录器
}

func NewConsoleLog() *ConsoleLog {
	cl := &ConsoleLog{}
	cl.Wrapping = fyne.TextWrapBreak
	cl.ExtendBaseWidget(cl)
	return cl
}

func (cl *ConsoleLog) Append(text string) {
	// 限制日志长度以避免性能下降
	const maxLogLength = 10000
	if len(cl.Text) > maxLogLength {
		cl.Text = "..." + cl.Text[len(cl.Text)-maxLogLength:]
	}
	cl.Text += text
	cl.Refresh()
}

func ConsoleWindow() {
	ShowCustomDialog(theme.FileTextIcon(), "日志中心(资产收集|暴破状况)", "清空日志", Frame(container.NewVScroll(Console)), func() {
		Console.SetText("")
	}, fyne.NewSize(500, 700))
}

func LogProgress(currentNum int64, countNum int, errinfo interface{}) {
	if (time.Now().Unix() - LogTime) > WaitTime {
		Console.Append(fmt.Sprintf("已完成 %v/%v %v\n", currentNum, countNum, errinfo))
		LogTime = time.Now().Unix()
	}
}

func ShowCustomDialog(icon fyne.Resource, title, confirm string, c fyne.CanvasObject, OnTapped func(), size fyne.Size) {
	ok := &widget.Button{Text: confirm, Icon: theme.ConfirmIcon(), Importance: widget.HighImportance, OnTapped: OnTapped}
	if OnTapped == nil {
		ok.Hide()
	}
	close := widget.NewButtonWithIcon("", theme.WindowCloseIcon(), nil)
	close.Importance = widget.LowImportance
	p := widget.NewModalPopUp(container.NewBorder(container.NewBorder(nil, nil, widget.NewIcon(icon), close, NewCenterLabel(title)),
		container.NewHBox(layout.NewSpacer(), ok, layout.NewSpacer()), nil, nil, c), global.Widgets.Win.Canvas())
	close.OnTapped = func() {
		p.Hide()
	}
	p.Resize(size)
	p.Show()
}

func NewFileEntry(placeholder string) *widget.Entry {
	e := widget.NewEntry()
	e.PlaceHolder = placeholder
	e.ActionItem = widget.NewButtonWithIcon("", theme.FileTextIcon(), func() {
		if !e.Disabled() {
			d := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
				if uc != nil {
					e.SetText(uc.URI().Path())
					e.Refresh()
				}
			}, global.Widgets.Win)
			d.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
			d.Show()
		}
	})
	return e
}

// 网站扫描与未授权检测/爆破部分
type ForwordEntry struct {
	*EntryBase
}

func NewForwordEntry() *ForwordEntry {
	fe := &ForwordEntry{EntryBase: NewEntryBase()}
	// ForwordEntry 特定的初始化代码
	fe.MultiLine = true
	fe.Wrapping = fyne.TextWrapBreak
	fe.ExtendBaseWidget(fe)
	return fe
}

func (fe *ForwordEntry) TappedSecondary(ev *fyne.PointEvent) {
	commonItems := fe.createCommonMenuItems() // 获取基类提供的通用菜单项
	customItems := []*fyne.MenuItem{
		fe.createBaseMenuItem("网站扫描", theme.MailSendIcon(), fe.scanWebsites),
		fe.createBaseMenuItem("未授权检测/爆破", theme.MailSendIcon(), fe.checkBurstUnauthorized),
	}
	m := &fyne.Menu{Items: append(commonItems, customItems...)} // 合并菜单项

	rc := widget.NewPopUpMenu(m, global.Widgets.Win.Canvas())
	rc.ShowAtPosition(ev.AbsolutePosition)
}

//

const defaultSpaceDirectory = "./config/space/"

type HistoryEntry struct {
	*EntryBase
	module string
}

func NewHistoryEntry(m string) *HistoryEntry {
	he := &HistoryEntry{
		EntryBase: NewEntryBase(),
		module:    m,
	}
	he.PlaceHolder = "Search..."
	he.ActionItem = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		he.SetText("")
	})
	he.ExtendBaseWidget(he)
	return he
}

type NumberEntry struct {
	*EntryBase
	Number int
}

func NewNumEntry(defaultNum string) *NumberEntry {
	dN, _ := strconv.Atoi(defaultNum)
	ne := &NumberEntry{
		EntryBase: NewEntryBase(),
		Number:    dN,
	}
	ne.Text = defaultNum
	ne.OnChanged = func(s string) {
		if _, err := strconv.Atoi(s); err != nil {
			ne.SetText(defaultNum)
			ne.Number, _ = strconv.Atoi(defaultNum)
		} else {
			ne.Number, _ = strconv.Atoi(ne.Text)
		}
	}
	ne.ExtendBaseWidget(ne)
	return ne
}

type SelectList struct {
	widget.Button
	Options []string
	Checked []string
	PopUp   *widget.PopUp
}

func NewSelectList(text string, options []string) *SelectList {
	sl := &SelectList{Options: options}
	sl.Text = text
	sl.Icon = theme.MenuDropDownIcon()
	sl.IconPlacement = widget.ButtonIconTrailingText // 使图标跟随在文本之后
	sl.Alignment = widget.ButtonAlignCenter
	cg := widget.NewCheckGroup(sl.Options, func(s []string) {
		if len(s) > 0 {
			sl.Icon = nil
			sl.SetText(strings.Join(s, " | "))
			sl.Checked = s
		} else {
			sl.Text = text
			sl.Checked = options
			sl.Icon = theme.MenuDropDownIcon()
			sl.Refresh()
		}
	})
	cg.Horizontal = false
	if sl.PopUp == nil {
		sl.PopUp = widget.NewPopUp(cg, global.Widgets.Win.Canvas())
	}
	sl.ExtendBaseWidget(sl)
	return sl
}

type ClickMode int

const (
	SuperClick ClickMode = iota
	SimpleClick
)

type SuperLabel struct {
	widget.Label
	ClickMode
}

// 给SuperLabel设置格式，比如发送模块只能给有用的地方使用
func NewSuperLabel(text string, mode ClickMode) *SuperLabel {
	label := &SuperLabel{ClickMode: mode}
	label.Text = text
	label.Alignment = fyne.TextAlignCenter
	label.Truncation = fyne.TextTruncateEllipsis
	label.ExtendBaseWidget(label)
	return label
}

const (
	sortOff int = iota
	sortAsc
	sortDesc
)

// 表格内容用superlabel，并且支持排序
func NewTableWithUpdateHeader(data *[][]string, width []float32, mode ClickMode) *widget.Table {
	var sorts = make([]int, len((*data)[0]))
	table := widget.NewTable(
		func() (rows int, cols int) {
			return len((*data)[1:]), len((*data)[0])
		}, func() fyne.CanvasObject {
			return NewSuperLabel("", mode)
		}, func(id widget.TableCellID, o fyne.CanvasObject) {
			if lb, ok := o.(*SuperLabel); ok {
				lb.SetText((*data)[1:][id.Row][id.Col])
			}
		})
	table.ShowHeaderRow = true
	table.CreateHeader = func() fyne.CanvasObject { // 一定得先CreateHeader才能使得表格表头为其他类型控件
		return widget.NewButton("000", func() {})
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		b := o.(*widget.Button)
		if id.Col == -1 {
			b.SetText(strconv.Itoa(id.Row))
			b.Importance = widget.LowImportance
			b.Disable()
		} else {
			b.SetText((*data)[0][id.Col])
			switch sorts[id.Col] {
			case sortAsc:
				b.Icon = theme.MoveUpIcon()
			case sortDesc:
				b.Icon = theme.MoveDownIcon()
			default:
				b.Icon = nil
			}
			b.OnTapped = func() {
				applySort(sorts, data, id.Col, table)
			}
			b.Enable()
			b.Refresh()
		}
	}
	for i, v := range width {
		table.SetColumnWidth(i, v)
	}
	return table
}

func NewVulnerabilityTable(data *[]common.VulnerabilityInfo, width []float32) *widget.Table {
	var sorts = make([]int, 5)
	table := widget.NewTable(
		func() (rows int, cols int) {
			return len((*data)), 5
		}, func() fyne.CanvasObject {
			return container.NewStack(NewSuperLabel("", SimpleClick), canvas.NewText("", color.White), &widget.Button{Icon: theme.ZoomInIcon(), Importance: widget.LowImportance})
		}, func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*fyne.Container).Objects[0].(*SuperLabel)
			r := o.(*fyne.Container).Objects[1].(*canvas.Text) // canvas.Text可以设置任意字体颜色
			b := o.(*fyne.Container).Objects[2].(*widget.Button)
			l.Show()
			b.Hide()
			r.Hide()
			if id.Col == 4 {
				l.Hide()
				b.Show()
				b.OnTapped = func() {
					req := NewMultiLineEntry((*data)[id.Row].Request, false)
					req.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
					resp := NewMultiLineEntry((*data)[id.Row].Response, false)
					resp.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
					hbox := container.NewHSplit(widget.NewCard("Request", "", req), widget.NewCard("Response", "", resp))
					ShowCustomDialog(mytheme.DetailIcon(), "数据包详情", "", container.NewStack(container.NewBorder(
						container.NewGridWithColumns(2, NewFormItem("漏洞ID:", widget.NewLabel((*data)[id.Row].Name)),
							NewFormItem("拓展信息:", widget.NewLabel((*data)[id.Row].TransInfo.ExtInfo))), nil, nil, nil, hbox)), nil, fyne.NewSize(800, 800))
				}
			} else if id.Col == 0 {
				l.SetText((*data)[id.Row].Id)
			} else if id.Col == 1 {
				l.SetText((*data)[id.Row].Name)
			} else if id.Col == 2 {
				l.Hide()
				r.Show()
				switch (*data)[id.Row].RiskLevel {
				case "CRITICAL":
					r.Color = &color.RGBA{75, 0, 130, 255}
				case "HIGH":
					r.Color = &color.RGBA{200, 0, 0, 200}
				case "MEDIUM":
					r.Color = &color.RGBA{255, 140, 0, 255}
				case "LOW":
					r.Color = &color.RGBA{0, 64, 128, 255}
				case "INFO":
					r.Color = &color.RGBA{0, 200, 0, 255}
				}
				r.Alignment = fyne.TextAlignCenter
				r.Text = (*data)[id.Row].RiskLevel
				r.Refresh()
			} else {
				l.SetText((*data)[id.Row].Url)
			}
		})
	table.ShowHeaderRow = true
	table.CreateHeader = func() fyne.CanvasObject { // 一定得先CreateHeader才能使得表格表头为其他类型控件
		return widget.NewButton("000", func() {})
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		b := o.(*widget.Button)
		if id.Col == -1 {
			b.SetText(strconv.Itoa(id.Row))
			b.Importance = widget.LowImportance
			b.Disable()
		} else {
			b.SetText(common.VulHeader[id.Col])
			switch sorts[id.Col] {
			case sortAsc:
				b.Icon = theme.MoveUpIcon()
			case sortDesc:
				b.Icon = theme.MoveDownIcon()
			default:
				b.Icon = nil
			}
			b.Importance = widget.MediumImportance
			b.OnTapped = func() {
				applySort2(sorts, data, id.Col, table)
			}
			b.Enable()
			b.Refresh()
		}
	}
	for i, v := range width {
		table.SetColumnWidth(i, v)
	}
	return table
}

type TappedSelect struct {
	widget.Select
	Parent *fyne.Container
}

func NewTappedSelect(options []string, parent *fyne.Container) *TappedSelect {
	s := &TappedSelect{
		Parent: parent,
	}
	s.Options = options
	s.ExtendBaseWidget(s)
	return s
}
