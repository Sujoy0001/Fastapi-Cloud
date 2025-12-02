package main

import "fmt"

func BinarySearchIterative(arr []int, target int) int {
    lo, hi := 0, len(arr)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if arr[mid] == target {
            return mid
        }
        if arr[mid] < target {
            lo = mid + 1
        } else {
            hi = mid - 1
        }
    }
    return -1
}

func Screach() {
    a := []int{1, 3, 5, 7, 9, 11, 15}
    fmt.Println("Array:", a)

    fmt.Println("Search 7 -> index:", BinarySearchIterative(a, 7))
    fmt.Println("Search 2 -> index:", BinarySearchIterative(a, 2))
}
