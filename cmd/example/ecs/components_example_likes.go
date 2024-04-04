package ecs

import (
	ecspb "github.com/delaneyj/geck/cmd/example/ecs/pb/gen/ecs/v1"
)

type Likes []Entity

func LikesFromEntity(c []Entity) Likes {
	return Likes(c)
}

func (c Likes) ToEntity() []Entity {
	return []Entity(c)
}

func (c Likes) ToEntities() []Entity {
	entities := make([]Entity, len(c))
	copy(entities, c)
	return entities
}

func (c Likes) ToU32s() []uint32 {
	u32s := make([]uint32, len(c))
	for i, e := range c {
		u32s[i] = e.val
	}
	return u32s
}

func LikesFromEntities(e ...Entity) Likes {
	c := make(Likes, len(e))
	copy(c, e)
	return c
}

//#region Events
//#endregion

func (e Entity) HasLikes() bool {
	return e.w.likesStore.Has(e)
}

func (e Entity) ReadLikes() ([]Entity, bool) {
	entities, ok := e.w.likesStore.Read(e)
	if !ok {
		return nil, false
	}
	return entities, true
}

func (e Entity) LikesContains(other Entity) bool {
	entities, ok := e.w.likesStore.Read(e)
	if !ok {
		return false
	}
	for _, entity := range entities {
		if entity == other {
			return true
		}
	}
	return false
}

func (e Entity) RemoveLikes(toRemove ...Entity) Entity {
	entities, ok := e.w.likesStore.Read(e)
	if !ok {
		return e
	}
	clean := make([]Entity, 0, len(entities))
	for _, tr := range toRemove {
		for _, entity := range entities {
			if entity != tr {
				clean = append(clean, entity)
			}
		}
	}
	e.w.likesStore.Set(clean, e)

	e.w.patch.LikesComponents[e.val] = nil

	return e
}

func (e Entity) RemoveAllLikes() Entity {
	e.w.likesStore.Remove(e)

	return e
}

func (e Entity) WritableLikes() (*Likes, bool) {
	return e.w.likesStore.Writeable(e)
}

func (e Entity) SetLikes(other ...Entity) Entity {
	e.w.likesStore.Set(Likes(other), e)

	e.w.patch.LikesComponents[e.w.resourceEntity.val] = Likes(other).ToPB()
	return e
}

func (w *World) SetLikes(c Likes, entities ...Entity) {
	if len(entities) == 0 {
		panic("no entities provided, are you sure you didn't mean to call SetLikesResource?")
	}
	w.likesStore.Set(c, entities...)
	w.patch.LikesComponents[w.resourceEntity.val] = c.ToPB()
}

func (w *World) RemoveLikes(entities ...Entity) {
	w.likesStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.LikesComponents[entity.val] = nil
	}
}

//#region Resources

// HasLikes checks if the world has a Likes}}
func (w *World) HasLikesResource() bool {
	return w.resourceEntity.HasLikes()
}

// Retrieve the Likes resource from the world
func (w *World) LikesResource() ([]Entity, bool) {
	return w.resourceEntity.ReadLikes()
}

// Set the Likes resource in the world
func (w *World) SetLikesResource(e ...Entity) Entity {
	w.resourceEntity.SetLikes(e...)
	return w.resourceEntity
}

// Remove the Likes resource from the world
func (w *World) RemoveLikesResource() Entity {
	w.resourceEntity.RemoveLikes()

	return w.resourceEntity
}

//#endregion

//#region Iterators

type LikesReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Likes]
}

func (iter *LikesReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *LikesReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *LikesReadIterator) NextLikes() (Entity, Likes) {
	e := iter.store.dense[iter.currIdx]
	c := iter.store.components[iter.currIdx]
	iter.currIdx++
	return e, c
}

func (iter *LikesReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) LikesReadIter() *LikesReadIterator {
	iter := &LikesReadIterator{
		w:     w,
		store: w.likesStore,
	}
	iter.Reset()
	return iter
}

type LikesWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Likes]
}

func (iter *LikesWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *LikesWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *LikesWriteIterator) NextLikes() (Entity, *Likes) {
	e := iter.store.dense[iter.currIdx]
	c := &iter.store.components[iter.currIdx]
	iter.currIdx--

	return e, c
}

func (iter *LikesWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) LikesWriteIter() *LikesWriteIterator {
	iter := &LikesWriteIterator{
		w:     w,
		store: w.likesStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) LikesEntities() []Entity {
	return w.likesStore.entities()
}

func (w *World) SetLikesSortFn(lessThan func(a, b Entity) bool) {
	w.likesStore.LessThan = lessThan
}

func (w *World) SortLikes() {
	w.likesStore.Sort()
}

func (w *World) ApplyLikesPatch(e Entity, patch *ecspb.LikesComponent) Entity {

	c := Likes(w.EntitiesFromU32s(patch.Entity...))
	e.w.likesStore.Set(c, e)
	return e
}

func (c Likes) ToPB() *ecspb.LikesComponent {
	pb := &ecspb.LikesComponent{
		Entity: c.ToU32s(),
	}
	return pb
}
