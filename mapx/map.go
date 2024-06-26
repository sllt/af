// Package mathx includes some functions to manipulate map.
package mathx

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"reflect"
	"sort"

	"github.com/sllt/af/slice"
)

// Keys returns a slice of the map's keys.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, len(m))

	var i int
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

// Values returns a slice of the map's values.
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, len(m))

	var i int
	for _, v := range m {
		values[i] = v
		i++
	}

	return values
}

// KeysBy creates a slice whose element is the result of function mapper invoked by every map's key.
func KeysBy[K comparable, V any, T any](m map[K]V, mapper func(item K) T) []T {
	keys := make([]T, 0, len(m))

	for k := range m {
		keys = append(keys, mapper(k))
	}

	return keys
}

// ValuesBy creates a slice whose element is the result of function mapper invoked by every map's value.
func ValuesBy[K comparable, V any, T any](m map[K]V, mapper func(item V) T) []T {
	keys := make([]T, 0, len(m))

	for _, v := range m {
		keys = append(keys, mapper(v))
	}

	return keys
}

// Merge maps, next key will overwrite previous key.
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V, 0)

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

// ForEach executes iteratee funcation for every key and value pair in map.
func ForEach[K comparable, V any](m map[K]V, iteratee func(key K, value V)) {
	for k, v := range m {
		iteratee(k, v)
	}
}

// Filter iterates over map, return a new map contains all key and value pairs pass the predicate function.
func Filter[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// FilterByKeys iterates over map, return a new map whose keys are all given keys.
func FilterByKeys[K comparable, V any](m map[K]V, keys []K) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if slice.Contain(keys, k) {
			result[k] = v
		}
	}
	return result
}

// FilterByValues iterates over map, return a new map whose values are all given values.
func FilterByValues[K comparable, V comparable](m map[K]V, values []V) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if slice.Contain(values, v) {
			result[k] = v
		}
	}
	return result
}

// OmitBy is the opposite of Filter, removes all the map elements for which the predicate function returns true.
func OmitBy[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if !predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// OmitByKeys the opposite of FilterByKeys, extracts all the map elements which keys are not omitted.
func OmitByKeys[K comparable, V any](m map[K]V, keys []K) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if !slice.Contain(keys, k) {
			result[k] = v
		}
	}
	return result
}

// OmitByValues the opposite of FilterByValues. remov all elements whose value are in the give slice.
func OmitByValues[K comparable, V comparable](m map[K]V, values []V) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if !slice.Contain(values, v) {
			result[k] = v
		}
	}
	return result
}

// Intersect iterates over maps, return a new map of key and value pairs in all given maps.
func Intersect[K comparable, V any](maps ...map[K]V) map[K]V {
	if len(maps) == 0 {
		return map[K]V{}
	}
	if len(maps) == 1 {
		return maps[0]
	}

	var result map[K]V

	reducer := func(m1, m2 map[K]V) map[K]V {
		m := make(map[K]V)
		for k, v1 := range m1 {
			if v2, ok := m2[k]; ok && reflect.DeepEqual(v1, v2) {
				m[k] = v1
			}
		}
		return m
	}

	reduceMaps := make([]map[K]V, 2)
	result = reducer(maps[0], maps[1])

	for i := 2; i < len(maps); i++ {
		reduceMaps[0] = result
		reduceMaps[1] = maps[i]
		result = reducer(reduceMaps[0], reduceMaps[1])
	}

	return result
}

// Minus creates a map of whose key in mapA but not in mapB.
func Minus[K comparable, V any](mapA, mapB map[K]V) map[K]V {
	result := make(map[K]V)

	for k, v := range mapA {
		if _, ok := mapB[k]; !ok {
			result[k] = v
		}
	}
	return result
}

// IsDisjoint two map are disjoint if they have no keys in common.
func IsDisjoint[K comparable, V any](mapA, mapB map[K]V) bool {
	for k := range mapA {
		if _, ok := mapB[k]; ok {
			return false
		}
	}
	return true
}

// Entry is a key/value pairs.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Entries transforms a map into array of key/value pairs.
func Entries[K comparable, V any](m map[K]V) []Entry[K, V] {
	entries := make([]Entry[K, V], 0, len(m))

	for k, v := range m {
		entries = append(entries, Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}

	return entries
}

// FromEntries creates a map based on a slice of key/value pairs
func FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V {
	result := make(map[K]V, len(entries))

	for _, v := range entries {
		result[v.Key] = v.Value
	}

	return result
}

// Transform a map to another type map.
func Transform[K1 comparable, V1 any, K2 comparable, V2 any](m map[K1]V1, iteratee func(key K1, value V1) (K2, V2)) map[K2]V2 {
	result := make(map[K2]V2, len(m))

	for k1, v1 := range m {
		k2, v2 := iteratee(k1, v1)
		result[k2] = v2
	}

	return result
}

// MapKeys transforms a map to other type map by manipulating it's keys.
func MapKeys[K comparable, V any, T comparable](m map[K]V, iteratee func(key K, value V) T) map[T]V {
	result := make(map[T]V, len(m))

	for k, v := range m {
		result[iteratee(k, v)] = v
	}

	return result
}

// MapValues transforms a map to other type map by manipulating it's values.
func MapValues[K comparable, V any, T any](m map[K]V, iteratee func(key K, value V) T) map[K]T {
	result := make(map[K]T, len(m))

	for k, v := range m {
		result[k] = iteratee(k, v)
	}

	return result
}

// HasKey checks if map has key or not.
// This function is used to replace the following boilerplate code:
// _, haskey := amap["baz"];
//
//	if haskey {
//	   fmt.Println("map has key baz")
//	}
func HasKey[K comparable, V any](m map[K]V, key K) bool {
	_, haskey := m[key]
	return haskey
}

// MapToStruct converts map to struct
func MapToStruct(m map[string]any, structObj any) error {
	for k, v := range m {
		err := setStructField(structObj, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func setStructField(structObj any, fieldName string, fieldValue any) error {
	structVal := reflect.ValueOf(structObj).Elem()

	fName := getFieldNameByJsonTag(structObj, fieldName)
	if fName == "" {
		return fmt.Errorf("Struct field json tag don't match map key : %s in obj", fieldName)
	}

	fieldVal := structVal.FieldByName(fName)

	if !fieldVal.IsValid() {
		return fmt.Errorf("No such field: %s in obj", fieldName)
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Cannot set %s field value", fieldName)
	}

	val := reflect.ValueOf(fieldValue)

	if fieldVal.Type() != val.Type() {

		if val.CanConvert(fieldVal.Type()) {
			fieldVal.Set(val.Convert(fieldVal.Type()))
			return nil
		}

		if m, ok := fieldValue.(map[string]any); ok {

			if fieldVal.Kind() == reflect.Struct {
				return MapToStruct(m, fieldVal.Addr().Interface())
			}

			if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
				}

				return MapToStruct(m, fieldVal.Interface())
			}

		}

		return fmt.Errorf("Map value type don't match struct field type")
	}

	fieldVal.Set(val)

	return nil
}

func getFieldNameByJsonTag(structObj any, jsonTag string) string {
	s := reflect.TypeOf(structObj).Elem()

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		tag := field.Tag
		name := tag.Get("json")

		if name == jsonTag {
			return field.Name
		}
	}

	return ""
}

// ToSortedSlicesDefault converts a map to two slices sorted by key: one for the keys and another for the values.
func ToSortedSlicesDefault[K constraints.Ordered, V any](m map[K]V) ([]K, []V) {
	keys := make([]K, 0, len(m))

	// store the map’s keys into a slice
	for k := range m {
		keys = append(keys, k)
	}

	// sort the slice of keys
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// adjust the order of values according to the sorted keys
	sortedValues := make([]V, len(keys))
	for i, k := range keys {
		sortedValues[i] = m[k]
	}

	return keys, sortedValues
}

// ToSortedSlicesWithComparator converts a map to two slices sorted by key and using a custom comparison function:
// one for the keys and another for the values.
func ToSortedSlicesWithComparator[K comparable, V any](m map[K]V, comparator func(a, b K) bool) ([]K, []V) {
	keys := make([]K, 0, len(m))

	// store the map’s keys into a slice
	for k := range m {
		keys = append(keys, k)
	}

	// sort the key slice using the provided comparison function
	sort.Slice(keys, func(i, j int) bool {
		return comparator(keys[i], keys[j])
	})

	// adjust the order of values according to the sorted keys
	sortedValues := make([]V, len(keys))
	for i, k := range keys {
		sortedValues[i] = m[k]
	}

	return keys, sortedValues
}
