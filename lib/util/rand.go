package util

import (
	"bufio"
	"bytes"
	rand2 "crypto/rand"
	"errors"
	"math/big"
	"math/rand"
	"os"
	"time"
)

const lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"

/*
letterIdxBits = 6：表示用6位来表示一个字母索引。在二进制中，6位可以表示从0到63的数字，足以覆盖一个典型字母表（例如26个英文字母）的索引。

letterIdxMask = 1<<letterIdxBits - 1：这是一个位掩码，用于获取字母索引。1<<letterIdxBits 创建一个只在第7位（从0开始计算）有一个1的二进制数（即1000000，也就是64），然后减去1，得到的是63，即111111，它用来确保索引值不会超出可表示的范围。

letterIdxMax = 63 / letterIdxBits：这表示一个63位数可以容纳多少个6位的字母索引。因为63位可以被6位数整除，所以这个值表示在63位的范围内，可以完全不重复地表示多少个字母索引。这里的值是63 / 6 = 10.5，通常会向下取整为10。
*/
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index;一个字母索引的位数（6位可以表示一个字母索引）
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits;所有位都为1，位数与letterIdxBits相同
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits;在63位中能容纳的最大字母索引数
)

// RandFromChoices 从choices中随机获取
func RandFromChoices(n int, choices string) string {
	b := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, cache, remain := n-1, r.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(choices) {
			b[i] = choices[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// RandLetter,根据长度来随机选取被小写的字母
func RandLetter(n int) string {
	return RandFromChoices(n, lowercaseLetters)
}

func RandomString(randSource *rand.Rand, letterBytes string, n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, randSource.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSource.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain -= 1
	}
	return string(b)
}

func NewRandomString(n int) string {
	var container string
	var randPool = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := bytes.NewBufferString(randPool)
	bigInt := big.NewInt(int64(b.Len()))
	for i := 0; i < n; i++ {
		randomInt, _ := rand2.Int(rand2.Reader, bigInt)
		container += string(randPool[randomInt.Int64()])
	}
	return container
}

func RandomUserAgent() string {
	file, err := os.Open("../../user-agents.txt")
	var useragent_b []string
	if err != nil {
		useragent_b = []string{
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.137 Safari/4E423F",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
			"Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2919.83 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2866.71 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2762.73 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2656.18 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML like Gecko) Chrome/44.0.2403.155 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.1 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.4; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2225.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2226.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2225.0 Safari/537.36",
			"Mozilla/5.0 (X11; OpenBSD i386) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36",
			"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2224.3 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.93 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.124 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2049.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 4.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2049.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.67 Safari/537.36",
			"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.67 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1944.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.3319.102 Safari/537.36",
			"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.2309.372 Safari/537.36",
			"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.2117.157 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.47 Safari/537.36",
			"Mozilla/5.0 (X11; Ubuntu; Linux i686 on x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2820.59 Safari/537.36",
			"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1866.237 Safari/537.36",
		}
		return useragent_b[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(useragent_b))]
	}
	defer file.Close()
	var userAgents []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		userAgents = append(userAgents, scanner.Text())
	}
	if len(userAgents) == 0 {
		return useragent_b[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(useragent_b))]
	}
	return userAgents[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(userAgents))]
}

// 包含上下限 [min, max]
func GetRandomIntWithAll(min, max int) int {
	rand.NewSource(time.Now().UnixNano())
	return int(rand.Intn(max-min+1) + min)
}

// 不包含上限 [min, max)
func GetRandomIntWithMin(min, max int) int {
	rand.NewSource(time.Now().UnixNano())
	return int(rand.Intn(max-min) + min)
}

// IntN returns a uniform random value in [0, max). It errors if max <= 0.
func IntN(max int) (int, error) {
	if max <= 0 {
		return 0, errors.New("max can't be <= 0")
	}
	nBig, err := rand2.Int(rand2.Reader, big.NewInt(int64(max)))
	if err != nil {
		return rand.Intn(max), nil
	}
	return int(nBig.Int64()), nil
}
