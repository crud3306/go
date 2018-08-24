package main

import (
	_ "fmt"
)

func main() {

}

// 冒泡排序 (方式1)
// ==================
func BubbleSort1(data []int) {
	for i := 1; i < len(data); i++ {
		isChange := false // 标志本轮有没有交换，没有交换就退出本轮
		for j := 0; j < len(data)-i; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
				isChange = true
			}
		}
		if !isChange {
			break
		}
	}
}


// 冒泡排序 (方式2)
// ==================
func BubbleSort2(items []int) {
    var (
        n       = len(items)
        swapped = true
    )
    for swapped {
        swapped = false
        for i := 0; i < n-1; i++ {
            if items[i] > items[i+1] {
                items[i+1], items[i] = items[i], items[i+1]
                swapped = true
            }
        }
        n = n - 1
    }
}


// 快速排序
// ==================
func QuickSort(data []int) {
    if len(data) <= 1 {
        return
    }
    pos := 0
    tail := len(data) - 1
    for pos < tail {
        if data[0] < data[pos+1] {
            data[pos+1], data[tail] = data[tail], data[pos+1]
            tail--
        } else {
            pos++
        }
    }

    data[0], data[pos] = data[pos], data[0]
    QuickSort(data[:pos])
    QuickSort(data[pos+1:])
}



// 选择排序
// ==================
func selectionSort(items []int) {
    var n = len(items)
    for i := 0; i < n; i++ {
        var minIdx = i
        for j := i; j < n; j++ {
            if items[j] < items[minIdx] {
                minIdx = j
            }
        }
        items[i], items[minIdx] = items[minIdx], items[i]
    }
}


// 插入排序
// ==================
func insertionSort(items []int) {
    var n = len(items)
    for i := 1; i < n; i++ {
        j := i
        for j > 0 {
            if items[j-1] > items[j] {
                items[j-1], items[j] = items[j], items[j-1]
            }
            j = j - 1
        }
    }
}


// xxx排序
// ==================




// http://blog.studygolang.com/2017/07/sorting-algorithms-primer/


