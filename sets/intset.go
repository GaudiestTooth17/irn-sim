package sets

type void struct{}
type IntSet map[int]void

func IntSetFromSlice(s []int) IntSet {
	set := make(IntSet)
	for _, v := range s {
		set[v] = void{}
	}
	return set
}

func EmptyIntSet() IntSet {
	return make(IntSet)
}

func (s IntSet) Values() []int {
	vs := make([]int, len(s))
	i := 0
	for v := range s {
		vs[i] = v
		i++
	}
	return vs
}

func (s IntSet) Contains(v int) bool {
	_, ok := s[v]
	return ok
}

func (s IntSet) Add(v int) {
	s[v] = void{}
}

// Return the intersection of the receiver and paramter
func (s0 IntSet) Intersection(s1 IntSet) IntSet {
	newSet := make(map[int]void)
	for v := range s0 {
		if _, ok := s1[v]; ok {
			newSet[v] = void{}
		}
	}
	return newSet
}

// Return the union of the receiver and paramter
func (s0 IntSet) Union(s1 IntSet) IntSet {
	newSet := make(IntSet)
	for v := range s0 {
		newSet.Add(v)
	}
	for v := range s1 {
		newSet.Add(v)
	}
	return newSet
}

func (s IntSet) Filter(keep func(v int) bool) IntSet {
	filteredSet := make(IntSet)
	for v := range s {
		if keep(v) {
			filteredSet.Add(v)
		}
	}
	return filteredSet
}
