package avlTree

type NumericKey uint64

func (lhs NumericKey) Compare(rhs interface{}) int {
	rhsKey, ok := rhs.(NumericKey)

	if !ok {
		rhsKey = 0
	}

	if lhs > rhsKey {
		return 1
	}

	if lhs < rhsKey {
		return -1
	}

	return 0
}
