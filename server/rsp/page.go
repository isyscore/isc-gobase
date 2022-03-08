package pojo

type PageReq struct {
	// 当前页面
	Current int

	// 每页大小
	Size int

	// 搜索的参数
	Param interface{}
}

type PageRsp struct {
	// 总个数
	Total int64

	// 分页数据
	Records []interface{}
}

func (pageReq *PageReq) GetStart() int {
	if pageReq.Current > 1 {
		return (pageReq.Current - 1) * pageReq.Size
	}
	return 0
}
