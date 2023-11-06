package sliceutil

// RemoveDuplicates removes duplicate elements from a given slice and returns a new slice with unique elements.
// The function preserves the order of the elements in the original slice.
//
// Parameters:
// - slice: The slice of type T containing potentially duplicate elements.
//
// Returns:
// - []T: The new slice containing unique elements from the input slice.
func RemoveDuplicates[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}

	for _, item := range slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}

	return list
}

// Intersect returns a new slice containing elements that are present in both input slices.
// The function preserves the order of the elements in the first slice.
//
// Parameters:
// - sliceA: The first slice containing elements of type T.
// - sliceB: The second slice containing elements of type T.
//
// Returns:
// - []T: The new slice containing elements that are present in both input slices.
func Intersect[T comparable](sliceA []T, sliceB []T) []T {
	result := make([]T, 0)

	for _, a := range sliceA {
		if Contains(sliceB, a) {
			result = append(result, a)
		}
	}

	return result
}

// ContainsAll returns whether the given first slice contains all elements of the second slice.
//
// Parameters:
// - sliceAll: The first slice of type T containing all elements.
// - slicePart: The second slice of type T containing elements that must be present in the first slice.
//
// Returns:
// - bool: Whether the given first slice contains all elements of the second slice.

func ContainsAll[T comparable](sliceAll []T, slicePart []T) bool {
	for _, e := range slicePart {
		if !Contains(sliceAll, e) {
			return false
		}
	}

	return true
}

// Contains returns whether the given string slice contains the given element.
//
// Parameters:
// - slice: The input slice of type T containing all elements.
// - element: The element of type T to check for.
//
// Returns:
// - bool: Whether the given slice contains the given element.
func Contains[T comparable](slice []T, element T) bool {
	for _, sliceElement := range slice {
		if sliceElement == element {
			return true
		}
	}

	return false
}
