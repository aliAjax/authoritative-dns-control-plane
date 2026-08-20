package operations

import (
	"sort"
	"strconv"
	"strings"
)

type Page struct {
	Limit  int
	Offset int
	Cursor string
}

func ParsePage(limit, offset string) Page {
	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)
	if l < 1 || l > 1000 {
		l = 100
	}
	if o < 0 {
		o = 0
	}
	return Page{Limit: l, Offset: o}
}
func Slice[T any](in []T, p Page) []T {
	if p.Offset >= len(in) {
		return []T{}
	}
	end := p.Offset + p.Limit
	if end > len(in) {
		end = len(in)
	}
	return append([]T(nil), in[p.Offset:end]...)
}
func SortStrings(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}
func UniqueStrings(in []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, x := range in {
		if strings.TrimSpace(x) != "" && !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}
func NextCursor(offset, limit, total int) string {
	if offset+limit >= total {
		return ""
	}
	return EncodeCursor(strconv.Itoa(offset + limit))
}
