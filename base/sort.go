package base

import "fmt"

func Quick(arr []int) {
	arr = append(arr, 10)
	//arr[len(arr)-1]=9
	fmt.Println("-------", arr, len(arr), cap(arr))
}
