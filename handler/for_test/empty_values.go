package test

import "auth/tool/random"

const (
	EmptyString = ""
	EmptyUint32 = uint32(0)
)

var (
	EmptyReplaceValueForString string
	EmptyReplaceValueForInt64 int64
	EmptyReplaceValueForUint32 uint32
)

func init() {
	EmptyReplaceValueForString = random.StringConsistOfIntWithLength(10)
	EmptyReplaceValueForInt64 = random.Int64WithLength(10)
	EmptyReplaceValueForUint32 = uint32(random.Int64WithLength(10))
}
