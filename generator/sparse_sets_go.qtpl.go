// Code generated by qtc from "sparse_sets_go.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// package generator
//

//line generator/sparse_sets_go.qtpl:3
package generator

//line generator/sparse_sets_go.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line generator/sparse_sets_go.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line generator/sparse_sets_go.qtpl:3
func streamsparseSetTemplate(qw422016 *qt422016.Writer, data *ecsTmplData) {
//line generator/sparse_sets_go.qtpl:3
	qw422016.N().S(`
package `)
//line generator/sparse_sets_go.qtpl:4
	qw422016.E().S(data.PackageName)
//line generator/sparse_sets_go.qtpl:4
	qw422016.N().S(`

const ssTombstoneIndex = -1

type SparseSet[T any] struct {
	sparse []int
	dense  []Entity
	data      []T
}

func NewSparseSet[T any]() *SparseSet[T] {
	return &SparseSet[T]{
	}
}

func (s *SparseSet[T]) search(idx int) int {
	sl, dl := len(s.sparse), len(s.dense)
	if idx >= sl || dl == 0 {
		return -1
	}

	denseIdx := s.sparse[idx]
	if denseIdx < 0 || denseIdx >= dl {
		return -1
	}

	dense := s.dense[denseIdx]
	if denseIdx < len(s.dense) && dense.Index() == idx {
		return denseIdx
	}

	return -1
}

func (s *SparseSet[T]) grow(idx int) {
	if idx >= len(s.sparse) {
		toGrow := idx - len(s.sparse) + 1
		arr := make([]int, toGrow)
		for i := range arr {
			arr[i] = ssTombstoneIndex
		}
		s.sparse = append(s.sparse, arr...)
	}
}

func (s *SparseSet[T]) Upsert(e Entity, c T) (old T, wasAdded bool) {
	idx := e.Index()
	searchIdx := s.search(idx)
	if searchIdx != -1 {
		old = s.data[searchIdx]
		s.data[searchIdx] = c
		return old, false
	}

	s.grow(idx)
	s.sparse[idx] = len(s.dense)
	s.dense = append(s.dense, e)
	s.data = append(s.data, c)
	return old, true
}

func (s *SparseSet[T]) Remove(e Entity) (wasRemoved bool) {
	idx := e.Index()
	sIdx := s.search(idx)
	if sIdx == -1 {
		return false
	}

	lastIdx := len(s.dense) - 1
	lastEntity := s.dense[lastIdx]
	lastEntityIdx := lastEntity.Index()
	s.dense[sIdx] = lastEntity
	s.data[sIdx] = s.data[lastIdx]
	s.sparse[lastEntityIdx] = sIdx
	s.sparse[idx] = ssTombstoneIndex
	s.dense = s.dense[:lastIdx]
	s.data = s.data[:lastIdx]
	return true
}

func (s *SparseSet[T]) Contains(e Entity) bool {
	return s.search(e.Index()) != -1
}

func (s *SparseSet[T]) Data(e Entity) (T, bool) {
	idx := s.search(e.Index())
	if idx == -1 {
		var zero T
		return zero, false
	}
	return s.data[idx], true
}

func (s *SparseSet[T]) DataMutable(e Entity) (*T,bool) {
	idx := s.search(e.Index())
	if idx == -1 {
		return nil, false
	}
	return &s.data[idx], true
}

func (s *SparseSet[T]) All(yield func(e Entity, c T) bool) {
	for i, e := range s.dense {
		data := s.data[i]
		if !yield(e, data) {
			break
		}
	}
}

func (s *SparseSet[T]) AllMutable(yield func(e Entity, c *T) bool) {
	for i, e := range s.dense {
		data := &s.data[i]
		if !yield(e, data) {
			break
		}
	}
}

func (s *SparseSet[T]) AllEntities(yield func(e Entity) bool) {
	for _, e := range s.dense {
		if !yield(e) {
			break
		}
	}
}

func (s *SparseSet[T]) Clear() {
	s.sparse = s.sparse[:0]
	s.dense = s.dense[:0]
	s.data = s.data[:0]
}

func (s *SparseSet[T]) Len() int {
	return len(s.dense)
}

func (s *SparseSet[T]) Cap() int {
	return cap(s.dense)
}

`)
//line generator/sparse_sets_go.qtpl:145
}

//line generator/sparse_sets_go.qtpl:145
func writesparseSetTemplate(qq422016 qtio422016.Writer, data *ecsTmplData) {
//line generator/sparse_sets_go.qtpl:145
	qw422016 := qt422016.AcquireWriter(qq422016)
//line generator/sparse_sets_go.qtpl:145
	streamsparseSetTemplate(qw422016, data)
//line generator/sparse_sets_go.qtpl:145
	qt422016.ReleaseWriter(qw422016)
//line generator/sparse_sets_go.qtpl:145
}

//line generator/sparse_sets_go.qtpl:145
func sparseSetTemplate(data *ecsTmplData) string {
//line generator/sparse_sets_go.qtpl:145
	qb422016 := qt422016.AcquireByteBuffer()
//line generator/sparse_sets_go.qtpl:145
	writesparseSetTemplate(qb422016, data)
//line generator/sparse_sets_go.qtpl:145
	qs422016 := string(qb422016.B)
//line generator/sparse_sets_go.qtpl:145
	qt422016.ReleaseByteBuffer(qb422016)
//line generator/sparse_sets_go.qtpl:145
	return qs422016
//line generator/sparse_sets_go.qtpl:145
}
