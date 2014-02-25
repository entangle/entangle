package utils

// String set.
type StringSet map[string]struct{}

// Add value to set.
func (s StringSet) Add(val string) {
	s[val] = struct{}{}
}

// Check if a value is contained in set.
func (s StringSet) Contains(val string) bool {
	_, contained := s[val]
	return contained
}

// Remove a value from set.
func (s StringSet) Remove(val string) {
	delete(s, val)
}

// Int set.
type IntSet map[int]struct{}

// Add value to set.
func (s IntSet) Add(val int) {
	s[val] = struct{}{}
}

// Check if a value is contained in set.
func (s IntSet) Contains(val int) bool {
	_, contained := s[val]
	return contained
}

// Remove a value from set.
func (s IntSet) Remove(val int) {
	delete(s, val)
}

// Uint set.
type UintSet map[uint]struct{}

// Add value to set.
func (s UintSet) Add(val uint) {
	s[val] = struct{}{}
}

// Check if a value is contained in set.
func (s UintSet) Contains(val uint) bool {
	_, contained := s[val]
	return contained
}

// Remove a value from set.
func (s UintSet) Remove(val uint) {
	delete(s, val)
}
