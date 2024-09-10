package excelx

// 创建文件、样式、表头、内容、支持拦截器自定义切面

// 导入、表头映射，数据过滤转化

const (
	ExcelTagKey = "excel"
	Pattern     = "name:(.*?);|index:(.*?);|width:(.*?);|needMerge:(.*?);|replace:(.*?);"
)

type ExcelTag struct {
	Value     interface{}
	Name      string // 表头标题
	Index     int    // 列下标(从0开始)
	Width     int    // 列宽
	NeedMerge bool   // 是否需要合并
	Replace   string // 替换（需要替换的内容_替换后的内容。比如：1_未开始 ==> 表示1替换为未开始）
}
