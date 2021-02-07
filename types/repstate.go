package types

import "math"

// RepState - A list that holds the state of the replica server
type RepState []rune

// Equals - Checks if repstates are equal
func (a RepState) Equals(b RepState) bool {
	if len(a) != len(b) {
		return false
	}
	for e := range a {
		if a[e] != b[e] {
			return false
		}
	}
	return true
}

// PrefixRelation - Checks if the two rep states have the same prefix
func (a RepState) PrefixRelation(b RepState) bool {
	var size = int(math.Min(float64(len(a)), float64(len(b))))
	for i := 0; i < size; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// GetCommonPrefix - Returns the common prefix of the two repstates
func (a RepState) GetCommonPrefix(b RepState) RepState {
	var prefix = make(RepState, 0)
	var size = int(math.Min(float64(len(a)), float64(len(b))))
	for i := 0; i < size; i++ {
		if a[i] == b[i] {
			prefix = append(prefix, a[i])
		}
	}
	return prefix
}
