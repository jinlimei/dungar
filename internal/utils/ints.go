package utils

// JoinInt64s will take all parameters of []int64 and return a map of them in the exact order.
func JoinInt64s(data ...[]int64) []int64 {
	out := make([]int64, 0)

	for _, group := range data {
		for _, val := range group {
			out = append(out, val)
		}
	}

	return out
}

// ReverseInt64s will reverse the array given
func ReverseInt64s(arr []int64) []int64 {
	rev := make([]int64, len(arr))

	k := 0
	for i := len(arr) - 1; i >= 0; i-- {
		rev[k] = arr[i]
		k++
	}

	return rev
}
