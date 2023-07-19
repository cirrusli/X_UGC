package bloomfilter

type BloomHashFunc func(string) int64 //创建一个函数类型

func NewFunc() []BloomHashFunc { //把所有hash函数放入切片中
	Funcs := make([]BloomHashFunc, 0)
	var f BloomHashFunc
	f = BKDRHash
	Funcs = append(Funcs, f)
	f = SDBMHash
	Funcs = append(Funcs, f)
	f = DJBHash
	Funcs = append(Funcs, f)
	return Funcs
}

func BKDRHash(str string) int64 {
	seed := int64(131) // 31 131 1313 13131 131313 etc..
	hash := int64(0)
	for i := 0; i < len(str); i++ {
		hash = (hash * seed) + int64(str[i])
	}
	return hash & 0x7FFFFFFF
}
func SDBMHash(str string) int64 {
	hash := int64(0)
	for i := 0; i < len(str); i++ {
		hash = int64(str[i]) + (hash << 6) + (hash << 16) - hash
	}
	return hash & 0x7FFFFFFF
}
func DJBHash(str string) int64 {
	hash := int64(0)
	for i := 0; i < len(str); i++ {
		hash = ((hash << 5) + hash) + int64(str[i])
	}
	return hash & 0x7FFFFFFF
}
