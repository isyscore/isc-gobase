package req

type PageRequest[T any] struct {
	// 当前页面
	Current int
	// 每页大小
	Size int
	// 搜索的参数
	Param T
}

func (p PageRequest[T]) Start() int {
	if p.Current > 1 {
		return (p.Current - 1) * p.Size
	}
	return 0
}
