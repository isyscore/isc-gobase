package req

type PageRequest[T any] struct {
	Current int `json:"current"` // 当前页面
	Size    int `json:"size"`    // 每页大小
	Param   T   `json:"param"`   // 搜索的参数
}

func (p PageRequest[T]) Start() int {
	if p.Current > 1 {
		return (p.Current - 1) * p.Size
	}
	return 0
}
