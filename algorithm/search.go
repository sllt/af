// Package algorithm contain some basic algorithm functions. eg. sort, search, list, linklist, stack, queue, tree, graph.
package algorithm

import "github.com/sllt/af/constraints"

// Search algorithms see https://github.com/TheAlgorithms/Go/tree/master/search

// LinearSearch return the index of target in slice base on equal function.
// If not found return -1
func LinearSearch[T any](slice []T, target T, equal func(a, b T) bool) int {
	for i, v := range slice {
		if equal(v, target) {
			return i
		}
	}
	return -1
}

// BinarySearch return the index of target within a sorted slice, use binary search (recursive call itself).
// If not found return -1.
func BinarySearch[T any](sortedSlice []T, target T, lowIndex, highIndex int, comparator constraints.Comparator) int {
	if highIndex < lowIndex || len(sortedSlice) == 0 {
		return -1
	}

	midIndex := int(lowIndex + (highIndex-lowIndex)/2)
	isMidValGreatTarget := comparator.Compare(sortedSlice[midIndex], target) == 1
	isMidValLessTarget := comparator.Compare(sortedSlice[midIndex], target) == -1

	if isMidValGreatTarget {
		return BinarySearch(sortedSlice, target, lowIndex, midIndex-1, comparator)
	} else if isMidValLessTarget {
		return BinarySearch(sortedSlice, target, midIndex+1, highIndex, comparator)
	}

	return midIndex
}

// BinaryIterativeSearch return the index of target within a sorted slice, use binary search (no recursive).
// If not found return -1.
func BinaryIterativeSearch[T any](sortedSlice []T, target T, lowIndex, highIndex int, comparator constraints.Comparator) int {
	startIndex := lowIndex
	endIndex := highIndex

	var midIndex int
	for startIndex <= endIndex {
		midIndex = int(startIndex + (endIndex-startIndex)/2)
		isMidValGreatTarget := comparator.Compare(sortedSlice[midIndex], target) == 1
		isMidValLessTarget := comparator.Compare(sortedSlice[midIndex], target) == -1

		if isMidValGreatTarget {
			endIndex = midIndex - 1
		} else if isMidValLessTarget {
			startIndex = midIndex + 1
		} else {
			return midIndex
		}
	}
	return -1
}
