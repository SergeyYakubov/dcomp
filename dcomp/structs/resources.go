package structs

import (
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"sort"
)

// Keeps resources and their priorities
type ResourcePrio map[string]int

type Resource struct {
	Server      server.Server
	DataManager server.Server
}

type pair struct {
	Key   string
	Value int
}

type pairList []pair

func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m ResourcePrio, reverse bool) pairList {
	p := make(pairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = pair{k, v}
		i++

	}
	if reverse {
		sort.Sort(sort.Reverse(p))
	} else {
		sort.Sort(p)
	}

	return p
}

func (prio ResourcePrio) Sort() (sorted []string) {
	sorted = make([]string, len(prio))
	p := sortMapByValue(prio, true)
	for i := range p {
		sorted[i] = p[i].Key
	}
	return
}
