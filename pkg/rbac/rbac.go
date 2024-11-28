package rbac

import "fmt"

func IntToBinary(n, length int) []int {
	fmt.Println(n)
	fmt.Println(length)
	binary := make([]int, length)
	for i := 0; i < length; i++ {
		binary[i] = n % 2
		n /= 2
	}
	fmt.Println(binary)
	return binary
}
