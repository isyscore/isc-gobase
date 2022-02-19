package isc

func Partition[T any](list []T, partition int) [][]T {
	return PartitionWithCal(list, func(int) int {
		return partition
	})
}

// PartitionWithCal 计算partition数
// f: 入参为数组长度,返回partition数
func PartitionWithCal[T any](list []T, f func(int) int) [][]T {
	var array [][]T

	length := len(list)
	if length == 0 {
		return array
	}

	partiton := f(length)
	if partiton <= 0 {
		array = append(array, list)
		return array
	}
	//list下标
	n := 0
	partitonLen := length / partiton
	for i := 0; i < partiton; i++ {
		var subList []T
		for j := 0; j < partitonLen; j++ {
			subList = append(subList, list[n])
			n++
		}
		//list有多余
		if i == partiton-1 {
			for k := n; k < length; k++ {
				subList = append(subList, list[k])
			}
		}
		array = append(array, subList)
	}
	return array
}
