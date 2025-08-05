package max

import (
	"sync"
)

func Find(input <-chan uint64) uint64 {
	resultFunction := sync.OnceValue(func() uint64 {
		var m uint64 = 0

		for value := range input {
			if value > m {
				m = value
			}
		}
		return m
	})
	return resultFunction()
}
