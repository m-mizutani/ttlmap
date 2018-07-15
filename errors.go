package ttlmap

import (
	"encoding/hex"
	"fmt"
)

// ErrNotFound is used when Map.Delete() fails because no element mathces the key
type ErrNotFound struct {
	key elementKey
}

func (x *ErrNotFound) Error() string {
	return fmt.Sprintf("No such key element: %s", hex.Dump(x.key))
}

// ErrDuplicatedKey is used when the key given as Set() argument already exists in table.
type ErrDuplicatedKey struct {
	key elementKey
}

func (x *ErrDuplicatedKey) Error() string {
	return fmt.Sprintf("Key is duplicated: %s", hex.Dump(x.key))
}

// ErrOverMaxTick is returned when tick as Set() argument is over time frame size.
type ErrOverMaxTick struct {
	max tick
	arg tick
}

func (x *ErrOverMaxTick) Error() string {
	return fmt.Sprintf("%d is larger than time frame size, should < %d", x.arg, x.max)
}
