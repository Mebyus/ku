package token

import "strconv"

func (k Kind) String() string {
	return strconv.FormatUint(uint64(k), 10)
}
