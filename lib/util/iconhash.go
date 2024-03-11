package util

import (
	"encoding/base64"

	"github.com/spaolacci/murmur3"
)

// 这里借鉴了多家iconhash项目的处理，感谢各位大佬
// Reference: https://github.com/Becivells/iconhash
// Mmh3Hash32 计算 mmh3 hash
// 实际上，MurmurHash3 算法产生的是一个整数类型的哈希值，通常用于快速键值对比较或存储在哈希表中。这个整数哈希值就是我们需要的结果。
func Mmh3Hash32(data []byte) (int32, error) {
	h32 := murmur3.New32()
	_, err := h32.Write(data)
	if err != nil {
		return 0, err
	}
	return int32(h32.Sum32()), nil
}

// Base64EncodeMIME 对数据进行Base64编码并按MIME类型格式化（每76个字符后添加换行）。
func Base64EncodeMIME(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
