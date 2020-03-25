package avlTree

import (
	"encoding/json"
)

type Set interface {
	json.Marshaler
	Contains(element Comparable) bool
	Remove(element Comparable) Set
	Add(element Comparable) Set
	LeastUpperBound(element Comparable) Set
	Size() int64
	Iter() Iterator
	IterGTE(lub Comparable) Iterator
}

type setIter struct {
	wrapped Iterator
}

type setImpl struct {
	root Node
}

func EmptySet() Set {
	return setImpl{
		root: EmptyTree(),
	}
}

func (this setImpl) Contains(element Comparable) bool {
	return !this.root.Find(element).IsEmpty()
}

func (this setImpl) Remove(element Comparable) Set {
	return setImpl{
		root: this.root.Delete(element),
	}
}

func (this setImpl) Add(element Comparable) Set {
	return setImpl{
		root: this.root.Update(element, nil),
	}
}

func (this setImpl) MarshalJSON() ([]byte, error) {
	var elements []interface{} = nil

	iter := this.Iter()

	for iter.MoveNext() {
		elements = append(elements, iter.Current())
	}

	return json.Marshal(elements)
}

func (this setImpl) LeastUpperBound(element Comparable) Set {
	return setImpl{
		root: this.root.LeastUpperBound(element),
	}
}

func (this setImpl) Size() int64 {
	return this.root.Size()
}

func (this setImpl) Iter() Iterator {
	return setIter{
		wrapped: this.root.Iter(),
	}
}

func (this setImpl) IterGTE(lub Comparable) Iterator {
	return setIter{
		wrapped: this.root.IterGTE(lub),
	}
}

func (this setIter) MoveNext() bool {
	return this.wrapped.MoveNext()
}

func (this setIter) Current() interface{} {
	return this.wrapped.Current().(Node).Key()
}
