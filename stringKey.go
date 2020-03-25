package avlTree

import "strings"

type StringKey string

func (lhs StringKey) Compare(rhs interface{}) int {
	rhsKey, ok := rhs.(StringKey)

	if !ok {
		rhsKey = ""
	}

	return strings.Compare(string(lhs), string(rhsKey))
}
