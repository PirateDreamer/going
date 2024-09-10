package excelx

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// 创建文件、样式、表头、内容、支持拦截器自定义切面

// 导入、表头映射，数据过滤转化

const (
	ExcelTagKey = "excel"
	Pattern     = "name:(.*?);|index:(.*?);|width:(.*?);|needMerge:(.*?);|replace:(.*?);"
)

type ExcelCellInfo struct {
	Value     interface{}
	Name      string // 表头标题
	Index     int    // 列下标(从0开始)
	Width     int    // 列宽
	NeedMerge bool   // 是否需要合并
	Replace   string // 替换（需要替换的内容_替换后的内容。比如：1_未开始 ==> 表示1替换为未开始）
	Style     string // 自定义样式
}

func NewDefaultExcelCellInfo() *ExcelCellInfo {
	return &ExcelCellInfo{
		Index: -1,
	}
}

func (e *ExcelCellInfo) ResolveTag(tag string) (err error) {
	re := regexp.MustCompile(Pattern)
	matches := re.FindAllStringSubmatch(tag, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			for i, val := range match {
				if i != 0 && val != "" {
					e.setValue(match, val)
				}
			}
		}
	} else {
		err = errors.New("未匹配到值")
		return
	}
	return
}

func (e *ExcelCellInfo) setValue(tag []string, value string) {
	if strings.Contains(tag[0], "name") {
		e.Name = value
	}
	if strings.Contains(tag[0], "index") {
		v, _ := strconv.ParseInt(value, 10, 8)
		e.Index = int(v)
	}
	if strings.Contains(tag[0], "width") {
		v, _ := strconv.ParseInt(value, 10, 8)
		e.Width = int(v)
	}
	if strings.Contains(tag[0], "replace") {
		e.Replace = value
	}
	if strings.Contains(tag[0], "style") {
		e.Style = value
	}
}
