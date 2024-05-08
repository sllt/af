package algorithm

import "github.com/sllt/af/constraints"

// BubbleSort applys the bubble sort algorithm to sort the collection, will change the original collection data.
func BubbleSort[T any](slice []T, comparator constraints.Comparator) {
	for i := 0; i < len(slice); i++ {
		breakTag := false
		for j := 0; j < len(slice)-1-i; j++ {
			isCurrGreatThanNext := comparator.Compare(slice[j], slice[j+1]) == 1
			if isCurrGreatThanNext {
				swap(slice, j, j+1)
				breakTag = true
			}
		}
		if !breakTag {
			break
		}
	}
}

// InsertionSort apply the insertion sort algorithm to sort the collection, will change the original collection data.
func InsertionSort[T any](slice []T, comparator constraints.Comparator) {
	for i := 0; i < len(slice); i++ {
		for j := i; j > 0; j-- {
			isPreLessThanCurrent := comparator.Compare(slice[j], slice[j-1]) == -1
			if isPreLessThanCurrent {
				swap(slice, j, j-1)
			} else {
				break
			}
		}
	}
}

// SelectionSort apply the selection sort algorithm to sort the collection, will change the original collection data.
func SelectionSort[T any](slice []T, comparator constraints.Comparator) {
	for i := 0; i < len(slice); i++ {
		min := i
		for j := i + 1; j < len(slice); j++ {
			if comparator.Compare(slice[j], slice[min]) == -1 {
				min = j
			}
		}
		swap(slice, i, min)
	}
}

// ShellSort applys the shell sort algorithm to sort the collection, will change the original collection data.
func ShellSort[T any](slice []T, comparator constraints.Comparator) {
	size := len(slice)

	gap := 1
	for gap < size/3 {
		gap = 3*gap + 1
	}

	for gap >= 1 {
		for i := gap; i < size; i++ {
			for j := i; j >= gap && comparator.Compare(slice[j], slice[j-gap]) == -1; j -= gap {
				swap(slice, j, j-gap)
			}
		}
		gap = gap / 3
	}
}

// QuickSort quick sorting for slice, lowIndex is 0 and highIndex is len(slice)-1.
func QuickSort[T any](slice []T, comparator constraints.Comparator) {
	quickSort(slice, 0, len(slice)-1, comparator)
}

func quickSort[T any](slice []T, lowIndex, highIndex int, comparator constraints.Comparator) {
	if lowIndex < highIndex {
		p := partition(slice, lowIndex, highIndex, comparator)
		quickSort(slice, lowIndex, p-1, comparator)
		quickSort(slice, p+1, highIndex, comparator)
	}
}

// partition split slice into two parts
func partition[T any](slice []T, lowIndex, highIndex int, comparator constraints.Comparator) int {
	p := slice[highIndex]
	i := lowIndex
	for j := lowIndex; j < highIndex; j++ {
		if comparator.Compare(slice[j], p) == -1 { //slice[j] < p
			swap(slice, i, j)
			i++
		}
	}

	swap(slice, i, highIndex)

	return i
}

// HeapSort applys the heap sort algorithm to sort the collection, will change the original collection data.
func HeapSort[T any](slice []T, comparator constraints.Comparator) {
	size := len(slice)

	for i := size/2 - 1; i >= 0; i-- {
		sift(slice, i, size-1, comparator)
	}
	for j := size - 1; j > 0; j-- {
		swap(slice, 0, j)
		sift(slice, 0, j-1, comparator)
	}
}

func sift[T any](slice []T, lowIndex, highIndex int, comparator constraints.Comparator) {
	i := lowIndex
	j := 2*i + 1

	temp := slice[i]
	for j <= highIndex {
		if j < highIndex && comparator.Compare(slice[j], slice[j+1]) == -1 { //slice[j] < slice[j+1]
			j++
		}
		if comparator.Compare(temp, slice[j]) == -1 { //tmp < slice[j]
			slice[i] = slice[j]
			i = j
			j = 2*i + 1
		} else {
			break
		}
	}
	slice[i] = temp
}

// MergeSort applys the merge sort algorithm to sort the collection, will change the original collection data.
func MergeSort[T any](slice []T, comparator constraints.Comparator) {
	mergeSort(slice, 0, len(slice)-1, comparator)
}

func mergeSort[T any](slice []T, lowIndex, highIndex int, comparator constraints.Comparator) {
	if lowIndex < highIndex {
		mid := (lowIndex + highIndex) / 2
		mergeSort(slice, lowIndex, mid, comparator)
		mergeSort(slice, mid+1, highIndex, comparator)
		merge(slice, lowIndex, mid, highIndex, comparator)
	}
}

func merge[T any](slice []T, lowIndex, midIndex, highIndex int, comparator constraints.Comparator) {
	i := lowIndex
	j := midIndex + 1
	temp := []T{}

	for i <= midIndex && j <= highIndex {
		//slice[i] < slice[j]
		if comparator.Compare(slice[i], slice[j]) == -1 {
			temp = append(temp, slice[i])
			i++
		} else {
			temp = append(temp, slice[j])
			j++
		}
	}

	if i <= midIndex {
		temp = append(temp, slice[i:midIndex+1]...)
	} else {
		temp = append(temp, slice[j:highIndex+1]...)
	}

	for k := 0; k < len(temp); k++ {
		slice[lowIndex+k] = temp[k]
	}
}

// CountSort apply the count sort algorithm to sort the collection, don't change the original collection data.
func CountSort[T any](slice []T, comparator constraints.Comparator) []T {
	size := len(slice)
	out := make([]T, size)

	for i := 0; i < size; i++ {
		count := 0
		for j := 0; j < size; j++ {
			//slice[i] > slice[j]
			if comparator.Compare(slice[i], slice[j]) == 1 {
				count++
			}
		}
		out[count] = slice[i]
	}

	return out
}

// swap two slice value at index i and j
func swap[T any](slice []T, i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
