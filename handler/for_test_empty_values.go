package handler

import "auth/tool/random"

const (
	emptyString = ""
	emptyUint32 = uint32(0)
)

var (
	emptyReplaceValueForString string
	emptyReplaceValueForInt64 int64
	emptyReplaceValueForUint32 uint32
)

func init() {
	emptyReplaceValueForString = random.StringConsistOfIntWithLength(10)
	emptyReplaceValueForInt64 = random.Int64WithLength(10)
	emptyReplaceValueForUint32 = uint32(random.Int64WithLength(10))
}
