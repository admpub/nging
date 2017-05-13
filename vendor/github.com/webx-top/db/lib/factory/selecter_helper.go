package factory

//MathGCD 最大公约数
func MathGCD(x, y int32) int32 {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

//MathNGCD n个数的最大公约数
func MathNGCD(a []int32, n int) int32 {
	if n == 2 {
		return MathGCD(a[0], a[1])
	} else if n < 2 {
		if len(a) > 1 {
			return a[0]
		}
		return 0
	}
	return MathGCD(a[n-1], MathNGCD(a, n-1))
}
