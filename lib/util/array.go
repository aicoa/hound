package util

import "sort"

//去除指定元素
type Pair struct {
	Key   string
	Value int
}

//字符串数组中移除特定的元素
func RemoveStringFromArray(slice []string, stringToRemove string) (newSlice []string) {
	for _, item := range slice {
		if item != stringToRemove {
			newSlice = append(newSlice, item)
		}
	}
	return newSlice
}

//泛型数组去重
func RemoveDuplicateElements[T comparable](arr []T) []T {
	encountered := map[T]bool{}
	result := []T{}
	for _, v := range arr {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

//指定元素的替换
func ReplaceElement(slice []string, old, new string) []string {
	newSlice := []string{}
	for _, v := range slice {
		if v == old {
			newSlice = append(newSlice, new)
		} else {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

//泛型判断数组中是否存在某个元素
//如果数组长度为0也返回成功 ->dirsearch排除状态码
func IsElementInArray[T comparable](t T, arr []T) bool {
	if len(arr) == 0 {
		return true
	}
	for _, v := range arr {
		if t == v {
			return true
		}
	}
	return false
}

//map按照value排序
func SortMapByValue(m map[string]int) []Pair {
	//将map转为切片
	var pairs []Pair
	for k, v := range m {
		pairs = append(pairs, Pair{k, v})
	}
	func(pairs []Pair) {
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].Value > pairs[j].Value })
	}(pairs)
	return pairs
}

func SplitInt(n, slice int) []int {
	var result []int
	for n > slice {
		result = append(result, slice)
		n -= slice
	}
	result = append(result, n)
	return result
}

// Combination 返回两个字符串切片s1和s2中元素的所有组合，组合中的元素由split分隔。
func Combination(s1, s2 []string, split string) []string {
	// 如果其中一个切片为空，返回空结果
	if len(s1) == 0 || len(s2) == 0 {
		return nil
	}

	// 预先分配足够的空间
	temp := make([]string, 0, len(s1)*len(s2))

	for _, v1 := range s1 {
		for _, v2 := range s2 {
			temp = append(temp, v1+split+v2)
		}
	}

	return temp
}
