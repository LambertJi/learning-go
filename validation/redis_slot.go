package validation

import (
	"learning-go/internals/redisx"
)

func TestRedisSlotStrategy() {
	keyCount := 1000000
	for i := 0; i < keyCount; i++ {
		redisx.Set("1", "1")
	}
}
