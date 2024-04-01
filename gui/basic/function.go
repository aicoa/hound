package basic

import (
	"bufio"
	"fmt"
	"hound/common"
	"hound/common/logger"
	"hound/gui/global"
	"hound/lib/util"
	"net"
	"net/url"
	"os"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

func (fe *ForwordEntry) scanWebsites() {
	// 用于网站扫描的代码逻辑
	go func() {
		lines := common.ParseTargets(fe.Text, common.Mode_Other)
		if len(lines) > 1 {
			global.Widgets.WebScanTarget.SetText("")
			for _, line := range lines {
				if strings.HasPrefix(line, "http") {
					global.Widgets.WebScanTarget.Text += line + "\n"
					global.Widgets.WebScanTarget.Refresh()
				}
			}
		}
	}()
}

func (fe *ForwordEntry) checkBurstUnauthorized() {
	// 用于暴破与未授权检测的代码逻辑
	go func() {
		lines := common.ParseTargets(fe.Text, common.Mode_Other)
		if len(lines) <= 0 {
			return
		} else {
			global.Widgets.PortBurstTarget.SetText("")
			for _, line := range lines {
				for protocol := range common.Userdict {
					if strings.HasPrefix(line, protocol) {
						global.Widgets.PortBurstTarget.Text += line + "\n"
						global.Widgets.PortBurstTarget.Refresh()
					}
				}
			}
		}
	}()
}

func (he *HistoryEntry) TappedSecondary(ev *fyne.PointEvent) {
	// 实现 HistoryEntry 特有的右键菜单逻辑
	commonItems := he.createCommonMenuItems()
	// 创建 HistoryEntry 特有的菜单项（例如历史记录）
	childmenu := fyne.NewMenu("")
	if _, err := os.Stat(defaultSpaceDirectory); err != nil {
		os.Mkdir(defaultSpaceDirectory, 0777)
	}
	historyFile := fmt.Sprintf("%v%v.txt", defaultSpaceDirectory, he.module)
	if _, err := os.Stat(historyFile); err == nil {
		childmenu.Items = append(childmenu.Items, &fyne.MenuItem{
			Label: "清空查询记录",
			Icon:  theme.DeleteIcon(),
			Action: func() {
				if err := os.Remove(historyFile); err != nil {
					dialog.ShowInformation("", "失败!", global.Widgets.Win)
				} else {
					dialog.ShowInformation("", "成功!", global.Widgets.Win)
				}
			},
		})

		childmenu.Items = append(childmenu.Items, fyne.NewMenuItemSeparator())

		for _, line := range common.ParseFile(historyFile) {
			lineCopy := line // 创建临时变量以避免闭包陷阱
			mi := fyne.NewMenuItem(lineCopy, func() {
				he.SetText(lineCopy)
				he.Refresh()
			})
			childmenu.Items = append(childmenu.Items, mi)
		}
	} else {
		childmenu.Items = append(childmenu.Items, &fyne.MenuItem{
			Label:  "未存在查询记录",
			Action: nil,
		})
	}
	historyItems := []*fyne.MenuItem{{Label: "历史记录", Icon: theme.HistoryIcon(), ChildMenu: childmenu}}
	// 合并通用菜单项和历史记录菜单项
	allItems := append(commonItems, historyItems...)

	// 创建并显示菜单
	m := &fyne.Menu{Items: allItems}
	rc := widget.NewPopUpMenu(m, global.Widgets.Win.Canvas())
	rc.ShowAtPosition(ev.AbsolutePosition)
}

func (he *HistoryEntry) WriteHistory(filename string) {
	// 实现 HistoryEntry 特有的历史记录写入逻辑
	file, err := os.OpenFile(defaultSpaceDirectory+filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logger.Info(err)
	}
	b, err := os.ReadFile(defaultSpaceDirectory + filename)
	if err != nil {
		logger.Info(err)
	}
	lines := strings.Split(string(b), "\n")
	write := bufio.NewWriter(file)
	if !util.IsElementInArray[string](he.Text, lines) {
		if len(b) > 0 {
			write.WriteString("\n" + he.Text)
		} else {
			write.WriteString(he.Text)
		}
	}
	write.Flush()
	file.Close()
}

func (sl *SelectList) Tapped(ev *fyne.PointEvent) {
	// 计算 PopUp 应该出现的位置
	h := fyne.CurrentApp().Driver().AbsolutePositionForObject(sl).Y + sl.Size().Height
	w := fyne.CurrentApp().Driver().AbsolutePositionForObject(sl).X

	// 设置 PopUp 的大小和位置，然后显示
	sl.PopUp.Resize(fyne.NewSize(sl.Size().Width, 0)) // 可以根据需要调整 PopUp 的高度
	sl.PopUp.ShowAtPosition(fyne.NewPos(w, h))
}

func (s *SuperLabel) TappedSecondary(ev *fyne.PointEvent) {
	memu := fyne.NewMenu("",
		&fyne.MenuItem{Label: "复制", Icon: theme.ContentCopyIcon(), Action: func() {
			clipboard.WriteAll(s.Text)
		}},
		&fyne.MenuItem{Label: "打开链接", Icon: theme.MailAttachmentIcon(), Action: func() {
			openlink(s.Text)
		}},
	)
	if s.ClickMode == SuperClick {
		memu.Items = append(memu.Items, fyne.NewMenuItemSeparator(),
			&fyne.MenuItem{Label: "将URL发送到Web扫描功能部分", Icon: theme.MailSendIcon(), Action: func() {
				sendwebscan(s.Text)
			}},
			&fyne.MenuItem{Label: "将IP发送到端口扫描功能部分", Icon: theme.MailSendIcon(), Action: func() {
				sendportscan(s.Text)
			}},
			&fyne.MenuItem{Label: "将域名发送到子域名暴破功能部分", Icon: theme.MailSendIcon(), Action: func() {
				sendsubdomain(s.Text)
			}})
	}
	memu.Refresh()
	rc := widget.NewPopUpMenu(memu, global.Widgets.Win.Canvas())
	rc.ShowAtPosition(ev.AbsolutePosition) // 面板出现在鼠标点击位置
}

func (s *SuperLabel) DoubleTapped(*fyne.PointEvent) {
	openlink(s.Text)
}

func openlink(text string) {
	if u, err := url.ParseRequestURI(text); err == nil {
		if err := fyne.CurrentApp().OpenURL(u); err != nil {
			dialog.ShowInformation("提示", "浏览器打开失败", global.Widgets.Win)
		}
	}
}

func sendwebscan(text string) {

	if _, err := url.ParseRequestURI(text); err == nil {
		global.Widgets.WebScanTarget.Text += text + "\n"
		global.Widgets.WebScanTarget.Refresh()
	}
}

// 该部分有问题
func sendportscan(text string) {
	if ip := net.ParseIP(text); ip != nil {
		global.Widgets.PortScanTarget.Text += text + "\n"
		global.Widgets.PortScanTarget.Refresh()
	}
}

func sendsubdomain(text string) {
	domains := util.RegDomain.FindAllString(text, -1)
	if len(domains) > 0 {
		for _, d := range domains {
			global.Widgets.SubdomainTarget.Text += d + "\n"
			global.Widgets.SubdomainTarget.Refresh()
		}
	}
}

// applySort 用于排序字符串类型的表格
func applySort(sorts []int, data *[][]string, col int, t *widget.Table) {
	toggleSortOrder(sorts, col)
	sort.Slice((*data)[1:], func(i, j int) bool {
		return compareValues((*data)[1:][i][col], (*data)[1:][j][col], sorts[col])
	})
	t.Refresh()
}

// applySort2 用于排序VulnerabilityInfo类型的表格
func applySort2(sorts []int, data *[]common.VulnerabilityInfo, col int, t *widget.Table) {
	toggleSortOrder(sorts, col)
	sort.Slice(*data, func(i, j int) bool {
		return compareVulnerabilityInfo((*data)[i], (*data)[j], col, sorts[col])
	})
	t.Refresh()
}

// toggleSortOrder 切换排序顺序
func toggleSortOrder(sorts []int, col int) {
	order := sorts[col] + 1
	if order > sortDesc {
		order = sortOff
	}
	for i := range sorts {
		sorts[i] = sortOff
	}
	sorts[col] = order
}

// compareValues 比较字符串值
func compareValues(a, b string, order int) bool {
	switch order {
	case sortAsc:
		return a < b
	case sortDesc:
		return a > b
	default:
		return false
	}
}

// compareVulnerabilityInfo 根据列比较VulnerabilityInfo
func compareVulnerabilityInfo(a, b common.VulnerabilityInfo, col, order int) bool {
	var comparisonResult bool
	switch col {
	case 0:
		comparisonResult = a.Id < b.Id
	case 1:
		comparisonResult = a.Name < b.Name
	case 2: // Assuming column 2 corresponds to RiskLevel
		comparisonResult = a.RiskLevel < b.RiskLevel
	case 3: // Assuming column 3 corresponds to Url
		comparisonResult = a.Url < b.Url
	default: // Add more cases if there are more fields
		comparisonResult = false
	}
	if order == sortDesc {
		return !comparisonResult
	}
	return comparisonResult
}

func (s *TappedSelect) TappedSecondary(ev *fyne.PointEvent) {
	s.Parent.Remove(s)
	s.Parent.Refresh()
}
