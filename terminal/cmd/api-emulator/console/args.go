package console

import "strconv"

// Args represents user input
type Args struct {
	Line string
	Args []string
}

// Len is args length
func (a Args) Len() int {
	return len(a.Args)
}

// String arg
func (a Args) String(idx int) (v string, ok bool) {
	if idx >= a.Len() {
		return
	}
	return a.Args[idx], true
}

// Int arg
func (a Args) Int(idx int) (v int, ok bool) {
	s, ok := a.String(idx)
	if !ok {
		return
	}
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return
	}
	return int(i), true
}

// Uint arg
func (a Args) Uint(idx int) (v uint, ok bool) {
	s, ok := a.String(idx)
	if !ok {
		return
	}
	u, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return
	}
	return uint(u), true
}

// Float64 arg
func (a Args) Float64(idx int) (v float64, ok bool) {
	s, ok := a.String(idx)
	if !ok {
		return
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return
	}
	return f, true
}
