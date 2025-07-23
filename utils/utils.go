package utils


const base62CharSet string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
func EncodeBase62(n int) string {
	res := ""
	for n > 0 {
		res = string(base62CharSet[n%62]) + res
		n = n / 62
	}
	return res
}
