package ecs

import "google.golang.org/protobuf/types/known/emptypb"

type Enemy struct{}

//#region Events
//#endregion

func (e Entity) HasEnemyTag() bool {
	return e.w.enemyStore.Has(e)
}

func (e Entity) TagWithEnemy() Entity {
	e.w.enemyStore.Set(e.w.enemyStore.zero, e)
	e.w.patch.EnemyTags[e.val] = empty
	return e
}

func (e Entity) RemoveEnemyTag() Entity {
	e.w.enemyStore.Remove(e)
	e.w.patch.EnemyTags[e.val] = nil
	return e
}

func (w *World) RemoveEnemyTags(entities ...Entity) {
	w.enemyStore.Remove(entities...)
	for _, entity := range entities {
		w.patch.EnemyTags[entity.val] = nil
	}
}

//#region Iterators

type EnemyReadIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Enemy]
}

func (iter *EnemyReadIterator) HasNext() bool {
	return iter.currIdx < iter.store.Len()
}

func (iter *EnemyReadIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx++
	return e
}

func (iter *EnemyReadIterator) Reset() {
	iter.currIdx = 0
}

func (w *World) EnemyReadIter() *EnemyReadIterator {
	iter := &EnemyReadIterator{
		w:     w,
		store: w.enemyStore,
	}
	iter.Reset()
	return iter
}

type EnemyWriteIterator struct {
	w       *World
	currIdx int
	store   *SparseSet[Enemy]
}

func (iter *EnemyWriteIterator) HasNext() bool {
	return iter.currIdx >= 0
}

func (iter *EnemyWriteIterator) NextEntity() Entity {
	e := iter.store.dense[iter.currIdx]
	iter.currIdx--

	return e
}

func (iter *EnemyWriteIterator) Reset() {
	iter.currIdx = iter.store.Len() - 1
}

func (w *World) EnemyWriteIter() *EnemyWriteIterator {
	iter := &EnemyWriteIterator{
		w:     w,
		store: w.enemyStore,
	}
	iter.Reset()
	return iter
}

//#endregion

func (w *World) EnemyEntities() []Entity {
	return w.enemyStore.entities()
}

func (w *World) ApplyEnemyPatch(e Entity, pb *emptypb.Empty) Entity {
	if pb == nil {
		e.RemoveEnemyTag()
	}
	e.TagWithEnemy()
	return e
}
