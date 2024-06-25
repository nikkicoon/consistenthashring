package pkg

import (
	"bytes"
	"cmp"
	"crypto/sha1"
	"fmt"
	"math/rand/v2"
	"slices"
)

func binarySearchBytes(arr [][]byte, target []byte, left int, right int) (int, bool) {
	if right == -1 {
		return 0, false
	}
	if len(arr) == 0 || right == -1 {
		// no elements in slice
		fmt.Println("no elements in slice")
		return 0, false
	}
	if left > right {
		// not found
		return -1, false
	}
	mid := left + (right-left)/2
	if bytes.Equal(arr[mid], target) {
		// found, return
		return mid, true
	} else if bytes.Compare(arr[mid], target) < 0 {
		// a less than b
		return binarySearchBytes(arr, target, mid+1, right)
	} else {
		return binarySearchBytes(arr, target, left, mid-1)
	}
}

// SortedInsertByte inserts t (byte slice) into ts (slice of byte slices).
func SortedInsertByte(ts [][]byte, t []byte) [][]byte {
	// find slot
	i, ok := binarySearchBytes(ts, t, 0, len(ts)-1)
	if !ok {
		// value not found in slice
		i = len(ts)
	}
	// if value is not found, assume index=0
	// fmt.Printf("inserting at position %d, length of ts %d\n", i, len(ts))
	return slices.Insert(ts, i, t)
}

// SortedInsert inserts t into ts, where t and ts are of cmp.Ordered type
func SortedInsert[T cmp.Ordered](ts []T, t T) []T {
	// find slot
	i, _ := slices.BinarySearch(ts, t)
	return slices.Insert(ts, i, t)
}

// CalculateHash returns the sha1 binary value for a given string `key`.
func CalculateHash(key string) []byte {
	// SHA1 hash, length of result: 20
	hash := sha1.New()
	hash.Write([]byte(key))
	return hash.Sum(nil)
	/*
		data := []byte(key)
		return sha1.Sum(data)

	*/
}

// ShuffleSlice randomizes the order of slice `s` via rand.Shuffle.
func ShuffleSlice[T any](s []T) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}
