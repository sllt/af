// Package mathx implements some functions for math calculation.
package mathx

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// Exponent calculate x^n.
func Exponent(x, n int64) int64 {
	if n == 0 {
		return 1
	}

	t := Exponent(x, n/2)

	if n%2 == 1 {
		return t * t * x
	}

	return t * t
}

// Fibonacci calculate fibonacci number before n.
func Fibonacci(first, second, n int) int {
	if n <= 0 {
		return 0
	}
	if n < 3 {
		return 1
	} else if n == 3 {
		return first + second
	} else {
		return Fibonacci(second, first+second, n-1)
	}
}

// Factorial calculate x!.
func Factorial(x uint) uint {
	var f uint = 1
	for ; x > 1; x-- {
		f *= x
	}
	return f
}

// Percent calculate the percentage of value to total.
func Percent(val, total float64, n int) float64 {
	if total == 0 {
		return float64(0)
	}
	tmp := val / total * 100
	result := RoundToFloat(tmp, n)

	return result
}

// RoundToString round off to n decimal places.
func RoundToString[T constraints.Float | constraints.Integer](x T, n int) string {
	tmp := math.Pow(10.0, float64(n))
	x *= T(tmp)
	r := math.Round(float64(x))
	result := strconv.FormatFloat(r/tmp, 'f', n, 64)
	return result
}

// RoundToFloat round off to n decimal places.
func RoundToFloat[T constraints.Float | constraints.Integer](x T, n int) float64 {
	tmp := math.Pow(10.0, float64(n))
	x *= T(tmp)
	r := math.Round(float64(x))
	return r / tmp
}

// TruncRound round off n decimal places.
func TruncRound[T constraints.Float | constraints.Integer](x T, n int) T {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n+1)+"f", x)
	temp := strings.Split(floatStr, ".")
	var newFloat string
	if len(temp) < 2 || n >= len(temp[1]) {
		newFloat = floatStr
	} else {
		newFloat = temp[0] + "." + temp[1][:n]
	}
	result, _ := strconv.ParseFloat(newFloat, 64)
	return T(result)
}

// FloorToFloat round down to n decimal places.
func FloorToFloat[T constraints.Float | constraints.Integer](x T, n int) float64 {
	tmp := math.Pow(10.0, float64(n))
	x *= T(tmp)
	r := math.Floor(float64(x))
	return r / tmp
}

// FloorToString round down to n decimal places.
func FloorToString[T constraints.Float | constraints.Integer](x T, n int) string {
	tmp := math.Pow(10.0, float64(n))
	x *= T(tmp)
	r := math.Floor(float64(x))
	result := strconv.FormatFloat(r/tmp, 'f', n, 64)
	return result
}

// CeilToFloat round up to n decimal places.
func CeilToFloat[T constraints.Float | constraints.Integer](x T, n int) float64 {
	tmp := math.Pow(10.0, float64(n))
	x *= T(tmp)
	r := math.Ceil(float64(x))
	return r / tmp
}

// CeilToString round up to n decimal places.
func CeilToString[T constraints.Float | constraints.Integer](x T, n int) string {
	tmp := math.Pow(10.0, float64(n))
	x *= T(tmp)
	r := math.Ceil(float64(x))
	result := strconv.FormatFloat(r/tmp, 'f', n, 64)
	return result
}

// Max return max value of numbers.
func Max[T constraints.Integer | constraints.Float](numbers ...T) T {
	max := numbers[0]

	for _, v := range numbers {
		if max < v {
			max = v
		}
	}

	return max
}

// MaxBy return the maximum value of a slice using the given comparator function.
func MaxBy[T any](slice []T, comparator func(T, T) bool) T {
	var max T

	if len(slice) == 0 {
		return max
	}

	max = slice[0]

	for i := 1; i < len(slice); i++ {
		val := slice[i]

		if comparator(val, max) {
			max = val
		}
	}

	return max
}

// Min return min value of numbers.
func Min[T constraints.Integer | constraints.Float](numbers ...T) T {
	min := numbers[0]

	for _, v := range numbers {
		if min > v {
			min = v
		}
	}

	return min
}

// MinBy return the minimum value of a slice using the given comparator function.
func MinBy[T any](slice []T, comparator func(T, T) bool) T {
	var min T

	if len(slice) == 0 {
		return min
	}

	min = slice[0]

	for i := 1; i < len(slice); i++ {
		val := slice[i]

		if comparator(val, min) {
			min = val
		}
	}

	return min
}

// Sum return sum of passed numbers.
func Sum[T constraints.Integer | constraints.Float](numbers ...T) T {
	var sum T

	for _, v := range numbers {
		sum += v
	}

	return sum
}

// Average return average value of numbers.
func Average[T constraints.Integer | constraints.Float](numbers ...T) T {
	var sum T
	n := T(len(numbers))

	for _, v := range numbers {
		sum += v
	}
	return sum / n
}

// Range creates a slice of numbers from start with specified count, element step is 1.
func Range[T constraints.Integer | constraints.Float](start T, count int) []T {
	size := count
	if count < 0 {
		size = -count
	}

	result := make([]T, size)

	for i, j := 0, start; i < size; i, j = i+1, j+1 {
		result[i] = j
	}

	return result
}

// RangeWithStep creates a slice of numbers from start to end with specified step.
func RangeWithStep[T constraints.Integer | constraints.Float](start, end, step T) []T {
	result := []T{}

	if start >= end || step == 0 {
		return result
	}

	for i := start; i < end; i += step {
		result = append(result, i)
	}

	return result
}

// AngleToRadian converts angle value to radian value.
func AngleToRadian(angle float64) float64 {
	radian := angle * (math.Pi / 180)
	return radian
}

// RadianToAngle converts radian value to angle value.
func RadianToAngle(radian float64) float64 {
	angle := radian * (180 / math.Pi)
	return angle
}

// PointDistance get two points distance.
func PointDistance(x1, y1, x2, y2 float64) float64 {
	a := x1 - x2
	b := y1 - y2
	c := math.Pow(a, 2) + math.Pow(b, 2)

	return math.Sqrt(c)
}

// IsPrime checks if number is prime number.
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}

	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

// GCD return greatest common divisor (GCD) of integers.
func GCD[T constraints.Integer](integers ...T) T {
	result := integers[0]

	for k := range integers {
		result = gcd(integers[k], result)

		if result == 1 {
			return 1
		}
	}

	return result
}

// find greatest common divisor (GCD)
func gcd[T constraints.Integer](a, b T) T {
	if b == 0 {
		return a
	}

	return gcd(b, a%b)
}

// LCM return Least Common Multiple (LCM) of integers.
func LCM[T constraints.Integer](integers ...T) T {
	result := integers[0]

	for k := range integers {
		result = lcm(integers[k], result)
	}

	return result
}

// find Least Common Multiple (LCM) via GCD.
func lcm[T constraints.Integer](a, b T) T {
	if a == 0 || b == 0 {
		panic("lcm function: provide non zero integers only.")
	}
	return a * b / gcd(a, b)
}

// Cos returns the cosine of the radian argument.
func Cos(radian float64, precision ...int) float64 {
	t := 1.0 / (2.0 * math.Pi)
	radian *= t
	radian -= 0.25 + math.Floor(radian+0.25)
	radian *= 16.0 * (math.Abs(radian) - 0.5)
	radian += 0.225 * radian * (math.Abs(radian) - 1.0)

	if len(precision) == 1 {
		return TruncRound(radian, precision[0])
	}

	return TruncRound(radian, 3)
}

// Sin returns the sine of the radian argument.
func Sin(radian float64, precision ...int) float64 {
	return Cos((math.Pi/2)-radian, precision...)
}

// Log returns the logarithm of base n.
func Log(n, base float64) float64 {
	return math.Log(n) / math.Log(base)
}

// Abs returns the absolute value of x.
func Abs[T constraints.Integer | constraints.Float](x T) T {
	if x < 0 {
		return (-x)
	}

	return x
}

// Div returns the result of x divided by y.
func Div[T constraints.Float | constraints.Integer](x T, y T) float64 {
	return float64(x) / float64(y)
}
