package agents

import (
	"github.com/iwind/TeaGo/maps"
)

// 通知级别类型
type NoticeLevel = uint8

// 通知级别常量
const (
	NoticeLevelNone    = uint8(0)
	NoticeLevelInfo    = uint8(1)
	NoticeLevelWarning = uint8(2)
	NoticeLevelError   = uint8(3)
	NoticeLevelSuccess = uint8(4)
)

// 所有的通知级别
func AllNoticeLevels() []maps.Map {
	return []maps.Map{
		{
			"name":    "信息",
			"code":    NoticeLevelInfo,
			"bgcolor": "#f8ffff",
			"color":   "#276f86",
		},
		{
			"name":    "警告",
			"code":    NoticeLevelWarning,
			"bgcolor": "#fffaf3",
			"color":   "#573a08",
		},
		{
			"name":    "错误",
			"code":    NoticeLevelError,
			"bgcolor": "#fff6f6",
			"color":   "#9f3a38",
		},
		{
			"name":    "成功",
			"code":    NoticeLevelSuccess,
			"bgcolor": "#fcfff5",
			"color":   "#2c662d",
		},
	}
}

// 获取通知级别名称
func FindNoticeLevelName(level NoticeLevel) string {
	for _, l := range AllNoticeLevels() {
		if l["code"] == level {
			return l["name"].(string)
		}
	}
	return "信息"
}

// 获取通知级别信息
func FindNoticeLevel(level NoticeLevel) maps.Map {
	for _, l := range AllNoticeLevels() {
		if l["code"] == level {
			return l
		}
	}
	return FindNoticeLevel(NoticeLevelInfo)
}
