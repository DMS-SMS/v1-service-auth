package handler

import "auth/tool/random"

const (
	emptyString = ""
	emptyUint32 = uint32(0)
)

var (
	emptyReplaceValueForString string
	emptyReplaceValueForInt int64
)

func init() {
	emptyReplaceValueForString = random.StringConsistOfIntWithLength(10)
	emptyReplaceValueForInt = random.Int64WithLength(10)
}