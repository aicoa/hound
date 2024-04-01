/*
 * @Author: aicoa
 * @Date: 2024-04-01 22:23:57
 * @Last Modified by: aicoa
 * @Last Modified time: 2024-04-01 22:27:14
 */

package menu

import (
	"encoding/json"
	"fmt"
	"hound/common"
	"hound/gui/basic"
	"hound/gui/global"
	"hound/lib/util"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type webscan struct {
	thread *basic.NumberEntry
}

type subdomain struct {
	dns1 *widget.Entry
	dns2 *widget.Entry
}

type portscan struct {
	thread  *basic.NumberEntry
	timeout *basic.NumberEntry
}

type proxys struct {
	enable   *widget.Check
	mode     *widget.Select
	address  *widget.Entry
	port     *basic.NumberEntry
	username *widget.Entry
	password *widget.Entry
}

type hunter struct {
	key *widget.Entry
}

type fofa struct {
	email *widget.Entry
	key   *widget.Entry
}

type quake struct {
	key *widget.Entry
}

func ConfigCenter() {
	// WEBSCAN
	var ws webscan
	ws.thread = basic.NewNumEntry(fmt.Sprint(common.Profile.WebScan.Thread))
	// SUBDOMAIN
	var sb subdomain
	sb.dns1 = &widget.Entry{Text: common.Profile.Subdomain.DNS1, OnChanged: func(s string) {
		common.Profile.Subdomain.DNS1 = util.RemoveIllegalChar(s)
	}}
	sb.dns2 = &widget.Entry{Text: common.Profile.Subdomain.DNS2, OnChanged: func(s string) {
		common.Profile.Subdomain.DNS2 = util.RemoveIllegalChar(s)
	}}
	// PORTSCAN
	var ps portscan
	ps.thread = basic.NewNumEntry(fmt.Sprint(common.Profile.PortScan.Thread))
	ps.timeout = basic.NewNumEntry(fmt.Sprint(common.Profile.PortScan.Timeout))
	// PROXY
	var pr proxys
	pr.enable = &widget.Check{Text: "Enable", Checked: common.Profile.Proxy.Enable, OnChanged: func(b bool) {
		if !b {
			common.Profile.Proxy.Enable = false
		} else {
			common.Profile.Proxy.Enable = true
		}
	}}
	pr.mode = &widget.Select{Options: []string{"HTTP", "SOCK5"}, Selected: common.Profile.Proxy.Mode, OnChanged: func(s string) {
		if s == "HTTP" {
			common.Profile.Proxy.Mode = "HTTP"
		} else {
			common.Profile.Proxy.Mode = "SOCK5"
		}
	}}
	pr.address = &widget.Entry{Text: common.Profile.Proxy.Address, OnChanged: func(s string) {
		common.Profile.Proxy.Address = s
	}}
	pr.port = basic.NewNumEntry(fmt.Sprint(common.Profile.Proxy.Port))
	pr.username = &widget.Entry{Text: common.Profile.Proxy.Username, OnChanged: func(s string) {
		common.Profile.Proxy.Username = s
	}}
	pr.password = &widget.Entry{Text: common.Profile.Proxy.Password, OnChanged: func(s string) {
		common.Profile.Proxy.Password = s
	}}
	// HUNTER
	var hu hunter
	hu.key = &widget.Entry{Text: common.Profile.Hunter.Api, OnChanged: func(s string) {
		common.Profile.Hunter.Api = util.RemoveIllegalChar(s)
	}}
	// FOFA
	var fo fofa
	fo.email = &widget.Entry{Text: common.Profile.Fofa.Email, OnChanged: func(s string) {
		common.Profile.Fofa.Email = util.RemoveIllegalChar(s)
	}}
	fo.key = &widget.Entry{Text: common.Profile.Fofa.Api, OnChanged: func(s string) {
		common.Profile.Fofa.Api = util.RemoveIllegalChar(s)
	}}
	// QUAKE
	var qk quake
	qk.key = &widget.Entry{Text: common.Profile.Quake.Api, OnChanged: func(s string) {
		common.Profile.Quake.Api = util.RemoveIllegalChar(s)
	}}
	accord := widget.NewAccordion(
		widget.NewAccordionItem("Webscan", widget.NewForm(
			widget.NewFormItem("Thread:", ws.thread),
		)),
		widget.NewAccordionItem("Subdomain", widget.NewForm(
			widget.NewFormItem("DNS1", sb.dns1),
			widget.NewFormItem("DNS2", sb.dns2),
		)),
		widget.NewAccordionItem("Portscan", widget.NewForm(
			widget.NewFormItem("Thread:", ps.thread),
			widget.NewFormItem("Timeout:", ps.timeout),
		)),
		widget.NewAccordionItem("Proxy", widget.NewForm(
			widget.NewFormItem("ON/OFF:", pr.enable),
			widget.NewFormItem("Mode:", pr.mode),
			widget.NewFormItem("Address:", pr.address),
			widget.NewFormItem("Port:", pr.port),
			widget.NewFormItem("Username:", pr.username),
			widget.NewFormItem("Password:", pr.password),
		)),
		widget.NewAccordionItem("Hunter", widget.NewForm(
			widget.NewFormItem("Key:", hu.key),
		)),
		widget.NewAccordionItem("Fofa", widget.NewForm(
			widget.NewFormItem("Email:", fo.email),
			widget.NewFormItem("Key:", fo.key),
		)),
		widget.NewAccordionItem("Quake", widget.NewForm(
			widget.NewFormItem("Key:", qk.key),
		)),
	)
	basic.ShowCustomDialog(theme.SettingsIcon(), "配置中心", "保存", accord, func() {
		pr.mode.SetSelected(pr.mode.Selected)
		b, err := json.MarshalIndent(common.Profile, "", "  ")
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed: %v", err), global.Widgets.Win)
			return
		}
		err = os.WriteFile("./config/config.json", b, 0777) // 将内容重新写入JSON文件
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed: %v", err), global.Widgets.Win)
			return
		}
		dialog.ShowInformation("", "Success", global.Widgets.Win)
	}, fyne.NewSize(500, 500))
}
