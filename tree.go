package avlTree

type Comparable interface {
	Compare(rhs interface{}) int
}
type Iterator interface {
	MoveNext() bool
	Current() interface{}
}

type nodeImpl struct {
	left   *nodeImpl
	right  *nodeImpl
	key    Comparable
	data   interface{}
	n      int64
	height int64
}

func EmptyTree() Node {
	var ret *nodeImpl = nil
	return ret
}

type Node interface {
	Find(key Comparable) Node
	Update(key Comparable, data interface{}) Node
	Left() Node
	Right() Node
	Delete(key Comparable) Node
	Size() int64
	IsEmpty() bool
	Height() int64
	Key() Comparable
	Data() interface{}
	LeastUpperBound(key Comparable) Node
	GreatestLowerBound(key Comparable) Node
	Iter() Iterator
	IterGTE(lub Comparable) Iterator
	Least() Node
	Most() Node
}

func (this *nodeImpl) Find(key Comparable) Node {
	if this.IsEmpty() {
		return this
	}

	comp := this.key.Compare(key)

	if comp < 0 {
		return this.Right().Find(key)
	}

	if comp > 0 {
		return this.Left().Find(key)
	}

	return this
}

func createNode(left *nodeImpl, right *nodeImpl, key Comparable, data interface{}) *nodeImpl {
	return &nodeImpl{
		left:   left,
		right:  right,
		key:    key,
		data:   data,
		n:      left.Size() + right.Size() + 1,
		height: max(left.Height(), right.Height()) + 1,
	}
}

func (this *nodeImpl) Update(key Comparable, data interface{}) Node {
	return this.doUpdate(key, data)
}

func (this *nodeImpl) doUpdate(key Comparable, data interface{}) *nodeImpl {
	if this == nil {
		return createNode(nil, nil, key, data)
	}

	comp := this.key.Compare(key)

	if comp == 0 {
		return createNode(this.left, this.right, key, data)
	}

	if comp < 0 {
		return createNode(this.left, this.right.doUpdate(key, data), this.key, this.data).rebalance()
	}

	return createNode(this.left.doUpdate(key, data), this.right, this.key, this.data).rebalance()
}

func (this *nodeImpl) Left() Node {
	if this == nil {
		return this
	}
	return this.left
}

func (this *nodeImpl) Right() Node {
	if this == nil {
		return this
	}
	return this.right
}

func (this *nodeImpl) Delete(key Comparable) Node {
	return this.doDelete(key)
}

func (this *nodeImpl) doDelete(key Comparable) *nodeImpl {
	if this == nil {
		return this
	}

	comp := this.key.Compare(key)

	if comp == 0 {
		return this.doDeleteCurrent().rebalance()
	}

	if comp < 0 {
		return createNode(
			this.left,
			this.right.doDelete(key),
			this.key,
			this.data,
		).rebalance()
	}

	return createNode(
		this.left.doDelete(key),
		this.right,
		this.key,
		this.data,
	).rebalance()
}

func (this *nodeImpl) doDeleteCurrent() *nodeImpl {
	if this.left.IsEmpty() {
		return this.right
	}

	if this.right.IsEmpty() {
		return this.left
	}

	replacement := this.left.rightMost()

	return createNode(
		this.left.doDelete(replacement.key),
		this.right,
		replacement.key,
		replacement.data,
	).rebalance()
}

func (this *nodeImpl) rightMost() *nodeImpl {
	current := this

	for current.right != nil {
		current = current.right
	}
	return current
}

func (this *nodeImpl) Size() int64 {
	if this == nil {
		return 0
	}
	return this.n
}

func (this *nodeImpl) IsEmpty() bool {
	return this == nil
}

func (this *nodeImpl) Height() int64 {
	if this == nil {
		return 0
	}
	return this.height
}

func (this *nodeImpl) balanceFactor() int64 {
	return this.Right().Height() - this.Left().Height()
}

func (this *nodeImpl) rebalance() *nodeImpl {
	balance := this.balanceFactor()
	if abs(balance) <= 1 {
		return this
	}

	if balance > 0 {
		rightBalance := this.right.balanceFactor()

		if rightBalance > 0 {
			return this.rotateLeft()
		}

		return this.rotateRightLeft()
	} else {
		leftBalance := this.left.balanceFactor()
		if leftBalance < 0 {
			return this.rotateRight()
		}

		return this.rotateLeftRight()
	}
}

func max(i, j int64) int64 {
	if i > j {
		return i
	}
	return j
}

func abs(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}

func (this *nodeImpl) rotateLeft() *nodeImpl {
	return createNode(
		createNode(this.left, this.right.left, this.key, this.data),
		this.right.right,
		this.right.key,
		this.right.data,
	)
}

func (this *nodeImpl) rotateRight() *nodeImpl {
	return createNode(
		this.left.left,
		createNode(this.left.right, this.right, this.key, this.data),
		this.left.key,
		this.left.data,
	)
}

func (this *nodeImpl) rotateRightLeft() *nodeImpl {
	return createNode(
		this.left,
		this.right.rotateRight(),
		this.key,
		this.data,
	).rotateLeft()
}

func (this *nodeImpl) rotateLeftRight() *nodeImpl {
	return createNode(
		this.left.rotateLeft(),
		this.right,
		this.key,
		this.data,
	).rotateRight()
}

func (this *nodeImpl) Key() Comparable {
	return this.key
}

func (this *nodeImpl) Data() interface{} {
	return this.data
}

func (this *nodeImpl) LeastUpperBound(key Comparable) Node {
	if this == nil {
		return this
	}

	comp := this.key.Compare(key)

	if comp < 0 {
		return this.right.LeastUpperBound(key)
	}

	if comp > 0 {
		ret := this.left.LeastUpperBound(key)

		if ret.IsEmpty() {
			return this
		}
		return ret
	}

	return this
}

func (this *nodeImpl) GreatestLowerBound(key Comparable) Node {
	if this == nil {
		return this
	}

	comp := this.key.Compare(key)

	if comp < 0 {
		ret := this.right.GreatestLowerBound(key)

		if ret.IsEmpty() {
			return this
		}
		return ret
	}

	if comp > 0 {
		return this.left.GreatestLowerBound(key)
	}

	return this
}

func (this *nodeImpl) Least() Node {
	if this.IsEmpty() {
		return this
	}

	var ret = this

	for !ret.left.IsEmpty() {
		ret = ret.left
	}

	return ret
}

func (this *nodeImpl) Most() Node {
	if this.IsEmpty() {
		return this
	}

	var ret = this

	for !ret.right.IsEmpty() {
		ret = ret.right
	}

	return ret
}

type treeIterator struct {
	stack   []*nodeImpl
	current *nodeImpl
}

func (this *nodeImpl) Iter() Iterator {

	ret := treeIterator{
		stack:   nil,
		current: this,
	}

	for ret.current != nil {
		ret.stack = append(ret.stack, ret.current)
		ret.current = ret.current.left
	}
	return &ret
}

func (this *nodeImpl) IterGTE(lub Comparable) Iterator {
	ret := treeIterator{
		stack:   nil,
		current: this,
	}

	for ret.current != nil {
		comp := ret.current.key.Compare(lub)

		if comp == 0 {
			ret.stack = append(ret.stack, ret.current)
			ret.current = nil
		} else if comp < 0 {
			ret.current = ret.current.right
		} else {
			ret.stack = append(ret.stack, ret.current)
			ret.current = ret.current.left
		}
	}

	return &ret

}

func (this *treeIterator) Current() interface{} {
	return this.current
}

func (this *treeIterator) MoveNext() bool {
	if this.current != nil {
		this.current = this.current.right
		for this.current != nil {
			this.stack = append(this.stack, this.current)
			this.current = this.current.left
		}
	}

	if len(this.stack) != 0 {
		this.current = this.stack[len(this.stack)-1]
		this.stack = this.stack[:len(this.stack)-1]
		return true
	}

	return false
}
