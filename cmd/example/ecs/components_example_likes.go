package ecs

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

	return e
}

func (w *World) SetLikes(c Likes, entities ...Entity) {
	w.likesStore.Set(c, entities...)
}

func (w *World) RemoveLikes(entities ...Entity) {
	w.likesStore.Remove(entities...)
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
func (w *World) SetLikesResource(c ...Entity) Entity {
	w.resourceEntity.SetLikes(c...)

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
