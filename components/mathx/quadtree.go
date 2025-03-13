package mathx

import (
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

// Package quadtree implements a quadtree using rectangular partitions.
// Each point exists in a unique node in the tree or as leaf nodes.
// This implementation is based off of the d3 implementation:
// https://github.com/mbostock/d3/wiki/Quadtree-Geom
var (
	// ErrPointOutsideOfBounds is returned when trying to add a point
	// to a quadtree and the point is outside the bounds used to create the tree.
	ErrPointOutsideOfBounds = fmt.Errorf("quadtree: point outside of bounds")
)

// Quadtree implements a two-dimensional recursive spatial subdivision
// of Vector2er[T]s. This implementation uses rectangular partitions.
type Quadtree[T constraints.Float] struct {
	bound Box2[T]
	root  *node[T]
}

type Vector2er[T constraints.Float] interface {
	Vector2() Vector2[T]
}

// A FilterFunc is a function that filters the points to search for.
type FilterFunc[T constraints.Float] func(p Vector2er[T]) bool

// node represents a node of the quad tree. Each node stores a Value
// and has links to its 4 children
type node[T constraints.Float] struct {
	Value    Vector2er[T]
	Children [4]*node[T]
}

// New creates a new quadtree for the given bound. Added points
// must be within this bound.
func New[T constraints.Float](bound Box2[T]) *Quadtree[T] {
	return &Quadtree[T]{bound: bound}
}

// Bound returns the bounds used for the quad tree.
func (q *Quadtree[T]) Bound() Box2[T] {
	return q.bound
}

// Add puts an object into the quad tree, must be within the quadtree bounds.
// This function is not thread-safe, ie. multiple goroutines cannot insert into
// a single quadtree.
func (q *Quadtree[T]) Add(p Vector2er[T]) error {
	if p == nil {
		return nil
	}

	v2 := p.Vector2()
	if !q.bound.ContainsPoint(v2) {
		return ErrPointOutsideOfBounds
	}

	if q.root == nil {
		q.root = &node[T]{
			Value: p,
		}
		return nil
	} else if q.root.Value == nil {
		q.root.Value = p
		return nil
	}

	q.add(q.root, p, p.Vector2(),
		// q.bound.Left(), q.bound.Right(),
		// q.bound.Bottom(), q.bound.Top(),
		q.bound.Min.X, q.bound.Max.X,
		q.bound.Min.Y, q.bound.Max.Y,
	)

	return nil
}

// add is the recursive search to find a place to add the point
func (q *Quadtree[T]) add(n *node[T], p Vector2er[T], point Vector2[T], left, right, bottom, top T) {
	i := 0

	// figure which child of this internal node the point is in.
	if cy := (bottom + top) / 2.0; point.Y <= cy {
		top = cy
		i = 2
	} else {
		bottom = cy
	}

	if cx := (left + right) / 2.0; point.X >= cx {
		left = cx
		i++
	} else {
		right = cx
	}

	if n.Children[i] == nil {
		n.Children[i] = &node[T]{Value: p}
		return
	} else if n.Children[i].Value == nil {
		n.Children[i].Value = p
		return
	}

	// proceed down to the child to see if it's a leaf yet and we can add the pointer there.
	q.add(n.Children[i], p, point, left, right, bottom, top)
}

// Remove will remove the pointer from the quadtree. By default it'll match
// using the points, but a FilterFunc can be provided for a more specific test
// if there are elements with the same point value in the tree. For example:
//
//	func(pointer Vector2er[T]) {
//		return pointer.(*MyType).ID == lookingFor.ID
//	}
func (q *Quadtree[T]) Remove(p Vector2er[T], eq FilterFunc[T]) bool {
	if eq == nil {
		point := p.Vector2()
		eq = func(pointer Vector2er[T]) bool {
			return point.Equals(pointer.Vector2())
		}
	}

	b := q.bound
	minDistSquared := math.MaxFloat64
	v := &findVisitor[T]{
		point:          p.Vector2(),
		filter:         eq,
		closestBound:   &b,
		minDistSquared: T(minDistSquared),
	}

	newVisit(v).Visit(q.root,
		// q.bound.Left(), q.bound.Right(),
		// q.bound.Bottom(), q.bound.Top(),
		q.bound.Min.X, q.bound.Max.X,
		q.bound.Min.Y, q.bound.Max.Y,
	)

	if v.closest == nil {
		return false
	}

	v.closest.Value = nil

	// if v.closest is NOT a leaf node, values will be shuffled up into this node.
	// if v.closest IS a leaf node, the call is a no-op but we can't delete
	// the now empty node because we don't know the parent here.
	//
	// Future adds will reuse this node if applicable.
	// Removing v.closest parent will cause this node to be removed,
	// but the parent will be a leaf with a nil value.
	removeNode(v.closest)
	return true
}

// removeNode is the recursive fixing up of the tree when we remove a node.
// It will pull up a child value into it's place. It will try to remove leaf nodes
// that are now empty, since their values got pulled up.
func removeNode[T constraints.Float](n *node[T]) bool {
	i := -1
	if n.Children[0] != nil {
		i = 0
	} else if n.Children[1] != nil {
		i = 1
	} else if n.Children[2] != nil {
		i = 2
	} else if n.Children[3] != nil {
		i = 3
	}

	if i == -1 {
		// all children are nil, can remove.
		// n.value ==  nil because it "pulled up" (or removed) by the caller.
		return true
	}

	n.Value = n.Children[i].Value
	n.Children[i].Value = nil

	removeThisChild := removeNode(n.Children[i])
	if removeThisChild {
		n.Children[i] = nil
	}

	return false
}

// Find returns the closest Value/Pointer in the quadtree.
// This function is thread safe. Multiple goroutines can read from
// a pre-created tree.
func (q *Quadtree[T]) Find(p Vector2[T]) Vector2er[T] {
	return q.Matching(p, nil)
}

// Matching returns the closest Value/Pointer in the quadtree for which
// the given filter function returns true. This function is thread safe.
// Multiple goroutines can read from a pre-created tree.
func (q *Quadtree[T]) Matching(p Vector2[T], f FilterFunc[T]) Vector2er[T] {
	if q.root == nil {
		return nil
	}

	b := q.bound
	minDistSquared := math.MaxFloat64
	v := &findVisitor[T]{
		point:          p,
		filter:         f,
		closestBound:   &b,
		minDistSquared: T(minDistSquared),
	}

	newVisit(v).Visit(q.root,
		// q.bound.Left(), q.bound.Right(),
		// q.bound.Bottom(), q.bound.Top(),
		q.bound.Min.X, q.bound.Max.X,
		q.bound.Min.Y, q.bound.Max.Y,
	)

	if v.closest == nil {
		return nil
	}
	return v.closest.Value
}

// KNearest returns k closest Value/Pointer in the quadtree.
// This function is thread safe. Multiple goroutines can read from a pre-created tree.
// An optional buffer parameter is provided to allow for the reuse of result slice memory.
// The points are returned in a sorted order, nearest first.
// This function allows defining a maximum distance in order to reduce search iterations.
func (q *Quadtree[T]) KNearest(buf []Vector2er[T], p Vector2[T], k int, maxDistance ...T) []Vector2er[T] {
	return q.KNearestMatching(buf, p, k, nil, maxDistance...)
}

// KNearestMatching returns k closest Value/Pointer in the quadtree for which
// the given filter function returns true. This function is thread safe.
// Multiple goroutines can read from a pre-created tree. An optional buffer
// parameter is provided to allow for the reuse of result slice memory.
// The points are returned in a sorted order, nearest first.
// This function allows defining a maximum distance in order to reduce search iterations.
func (q *Quadtree[T]) KNearestMatching(buf []Vector2er[T], p Vector2[T], k int, f FilterFunc[T], maxDistance ...T) []Vector2er[T] {
	if q.root == nil {
		return nil
	}

	b := q.bound
	maxDistSquared := math.MaxFloat64
	v := &nearestVisitor[T]{
		point:          p,
		filter:         f,
		k:              k,
		maxHeap:        make(maxHeap[T], 0, k+1),
		closestBound:   &b,
		maxDistSquared: T(maxDistSquared),
	}

	if len(maxDistance) > 0 {
		v.maxDistSquared = maxDistance[0] * maxDistance[0]
	}

	newVisit(v).Visit(q.root,
		// q.bound.Left(), q.bound.Right(),
		// q.bound.Bottom(), q.bound.Top(),
		q.bound.Min.X, q.bound.Max.X,
		q.bound.Min.Y, q.bound.Max.Ceil().Y,
	)

	//repack result
	if cap(buf) < len(v.maxHeap) {
		buf = make([]Vector2er[T], len(v.maxHeap))
	} else {
		buf = buf[:len(v.maxHeap)]
	}

	for i := len(v.maxHeap) - 1; i >= 0; i-- {
		buf[i] = v.maxHeap[0].point
		v.maxHeap.Pop()
	}

	return buf
}

// InBound returns a slice with all the pointers in the quadtree that are
// within the given bound. An optional buffer parameter is provided to allow
// for the reuse of result slice memory. This function is thread safe.
// Multiple goroutines can read from a pre-created tree.
func (q *Quadtree[T]) InBound(buf []Vector2er[T], b Box2[T]) []Vector2er[T] {
	return q.InBoundMatching(buf, b, nil)
}

// InBoundMatching returns a slice with all the pointers in the quadtree that are
// within the given bound and matching the give filter function. An optional buffer
// parameter is provided to allow for the reuse of result slice memory. This function
// is thread safe.  Multiple goroutines can read from a pre-created tree.
func (q *Quadtree[T]) InBoundMatching(buf []Vector2er[T], b Box2[T], f FilterFunc[T]) []Vector2er[T] {
	if q.root == nil {
		return nil
	}

	var p []Vector2er[T]
	if buf != nil {
		p = buf[:0]
	}
	v := &inBoundVisitor[T]{
		bound:    &b,
		pointers: p,
		filter:   f,
	}

	newVisit(v).Visit(q.root,
		// q.bound.Left(), q.bound.Right(),
		// q.bound.Bottom(), q.bound.Top(),
		q.bound.Min.X, q.bound.Max.X,
		q.bound.Min.Y, q.bound.Max.Y,
	)

	return v.pointers
}

// The visit stuff is a more go like (hopefully) implementation of the
// d3.quadtree.visit function. It is not exported, but if there is a
// good use case, it could be.

type visitor[T constraints.Float] interface {
	// Bound returns the current relevant bound so we can prune irrelevant nodes
	// from the search. Using a pointer was benchmarked to be 5% faster than
	// having to copy the bound on return. go1.9
	Bound() *Box2[T]
	Visit(n *node[T])

	// Point should return the specific point being search for, or null if there
	// isn't one (ie. searching by bound). This helps guide the search to the
	// best child node first.
	Vector2() Vector2[T]
}

// visit provides a framework for walking the quad tree.
// Currently used by the `Find` and `InBound` functions.
type visit[T constraints.Float] struct {
	visitor visitor[T]
}

func newVisit[T constraints.Float](v visitor[T]) *visit[T] {
	return &visit[T]{
		visitor: v,
	}
}

func (v *visit[T]) Visit(n *node[T], left, right, bottom, top T) {
	b := v.visitor.Bound()
	// if left > b.Right() || right < b.Left() ||
	// 	bottom > b.Top() || top < b.Bottom() {
	// 	return
	// }
	if left > b.Max.X || right < b.Min.X ||
		bottom > b.Max.Y || top < b.Min.Y {
		return
	}

	if n.Value != nil {
		v.visitor.Visit(n)
	}

	if n.Children[0] == nil && n.Children[1] == nil &&
		n.Children[2] == nil && n.Children[3] == nil {
		// no children check
		return
	}

	cx := (left + right) / 2.0
	cy := (bottom + top) / 2.0

	i := childIndex(cx, cy, v.visitor.Vector2())
	for j := i; j < i+4; j++ {
		if n.Children[j%4] == nil {
			continue
		}

		if k := j % 4; k == 0 {
			v.Visit(n.Children[0], left, cx, cy, top)
		} else if k == 1 {
			v.Visit(n.Children[1], cx, right, cy, top)
		} else if k == 2 {
			v.Visit(n.Children[2], left, cx, bottom, cy)
		} else if k == 3 {
			v.Visit(n.Children[3], cx, right, bottom, cy)
		}
	}
}

type findVisitor[T constraints.Float] struct {
	point          Vector2[T]
	filter         FilterFunc[T]
	closest        *node[T]
	closestBound   *Box2[T]
	minDistSquared T
}

func (v *findVisitor[T]) Bound() *Box2[T] {
	return v.closestBound
}

func (v *findVisitor[T]) Vector2() Vector2[T] {
	return v.point
}

func (v *findVisitor[T]) Visit(n *node[T]) {
	// skip this pointer if we have a filter and it doesn't match
	if v.filter != nil && !v.filter(n.Value) {
		return
	}

	point := n.Value.Vector2()
	if d := point.DistanceToSquared(v.point); d < v.minDistSquared {
		v.minDistSquared = d
		v.closest = n

		d = T(math.Sqrt(float64(d)))
		v.closestBound.Min.X = v.point.X - d
		v.closestBound.Max.X = v.point.X + d
		v.closestBound.Min.Y = v.point.Y - d
		v.closestBound.Max.Y = v.point.Y + d
	}
}

// type pointsQueueItem struct {
// 	point    Vector2er[T]
// 	distance float64 // distance to point and priority inside the queue
// 	index    int     // point index in queue
// }

// type pointsQueue []pointsQueueItem

// func newPointsQueue(capacity int) pointsQueue {
// 	// We make capacity+1 because we need additional place for the greatest element
// 	return make([]pointsQueueItem, 0, capacity+1)
// }

// func (pq pointsQueue) Len() int { return len(pq) }

// func (pq pointsQueue) Less(i, j int) bool {
// 	// We want pop longest distances so Less was inverted
// 	return pq[i].distance > pq[j].distance
// }

// func (pq pointsQueue) Swap(i, j int) {
// 	pq[i], pq[j] = pq[j], pq[i]
// 	pq[i].index = i
// 	pq[j].index = j
// }

// func (pq *pointsQueue) Push(x interface{}) {
// 	n := len(*pq)
// 	item := x.(pointsQueueItem)
// 	item.index = n
// 	*pq = append(*pq, item)
// }

// func (pq *pointsQueue) Pop() interface{} {
// 	old := *pq
// 	n := len(old)
// 	item := old[n-1]
// 	item.index = -1
// 	*pq = old[0 : n-1]
// 	return item
// }

type nearestVisitor[T constraints.Float] struct {
	point          Vector2[T]
	filter         FilterFunc[T]
	k              int
	maxHeap        maxHeap[T]
	closestBound   *Box2[T]
	maxDistSquared T
}

func (v *nearestVisitor[T]) Bound() *Box2[T] {
	return v.closestBound
}

func (v *nearestVisitor[T]) Vector2() Vector2[T] {
	return v.point
}

func (v *nearestVisitor[T]) Visit(n *node[T]) {
	// skip this pointer if we have a filter and it doesn't match
	if v.filter != nil && !v.filter(n.Value) {
		return
	}

	point := n.Value.Vector2()
	if d := point.DistanceToSquared(v.point); d < v.maxDistSquared {
		v.maxHeap.Push(n.Value, d)
		if len(v.maxHeap) > v.k {

			v.maxHeap.Pop()

			// Actually this is a hack. We know how heap works and obtain
			// top element without function call
			top := v.maxHeap[0]

			v.maxDistSquared = top.distance

			// We have filled queue, so we start to restrict searching range
			d = T(math.Sqrt(float64(top.distance)))
			v.closestBound.Min.X = v.point.X - d
			v.closestBound.Max.X = v.point.X + d
			v.closestBound.Min.Y = v.point.Y - d
			v.closestBound.Max.Y = v.point.Y + d
		}
	}
}

type inBoundVisitor[T constraints.Float] struct {
	bound    *Box2[T]
	pointers []Vector2er[T]
	filter   FilterFunc[T]
}

func (v *inBoundVisitor[T]) Bound() *Box2[T] {
	return v.bound
}

func (v *inBoundVisitor[T]) Vector2() (p Vector2[T]) {
	return
}

func (v *inBoundVisitor[T]) Visit(n *node[T]) {
	if v.filter != nil && !v.filter(n.Value) {
		return
	}

	p := n.Value.Vector2()
	if v.bound.Min.X > p.X || v.bound.Max.X < p.X ||
		v.bound.Min.Y > p.Y || v.bound.Max.Y < p.Y {
		return

	}
	v.pointers = append(v.pointers, n.Value)
}

func childIndex[T constraints.Float](cx, cy T, point Vector2[T]) int {
	i := 0
	if point.Y <= cy {
		i = 2
	}

	if point.X >= cx {
		i++
	}

	return i
}

// maxHeap is used for the knearest list. We need a way to maintain
// the furthest point from the query point in the list, hence maxHeap.
// When we find a point closer than the furthest away, we remove
// furthest and add the new point to the heap.
type maxHeap[T constraints.Float] []heapItem[T]

type heapItem[T constraints.Float] struct {
	point    Vector2er[T]
	distance T
}

func (h *maxHeap[T]) Push(point Vector2er[T], distance T) {
	prevLen := len(*h)
	*h = (*h)[:prevLen+1]
	(*h)[prevLen].point = point
	(*h)[prevLen].distance = distance

	i := len(*h) - 1
	for i > 0 {
		up := ((i + 1) >> 1) - 1
		parent := (*h)[up]

		if distance < parent.distance {
			// parent is further so we're done fixing up the heap.
			break
		}

		// swap nodes
		// (*h)[i] = parent
		(*h)[i].point = parent.point
		(*h)[i].distance = parent.distance

		// (*h)[up] = item
		(*h)[up].point = point
		(*h)[up].distance = distance

		i = up
	}
}

// Pop returns the "greatest" item in the list.
// The returned item should not be saved across push/pop operations.
func (h *maxHeap[T]) Pop() {
	lastItem := (*h)[len(*h)-1]
	(*h) = (*h)[:len(*h)-1]

	mh := (*h)
	if len(mh) == 0 {
		return
	}

	// move the last item to the top and reset the heap
	mh[0].point = lastItem.point
	mh[0].distance = lastItem.distance

	i := 0
	for {
		right := (i + 1) << 1
		left := right - 1

		childIndex := i
		child := mh[childIndex]

		// swap with biggest child
		if left < len(mh) && child.distance < mh[left].distance {
			childIndex = left
			child = mh[left]
		}

		if right < len(mh) && child.distance < mh[right].distance {
			childIndex = right
			child = mh[right]
		}

		// non bigger, so quit
		if childIndex == i {
			break
		}

		// swap the nodes
		mh[i].point = child.point
		mh[i].distance = child.distance

		mh[childIndex].point = lastItem.point
		mh[childIndex].distance = lastItem.distance

		i = childIndex
	}
}
