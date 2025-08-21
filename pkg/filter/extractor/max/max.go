package max

import (
	"sync"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
)

func Find(input <-chan filter.Artifact) uint64 {
	resultFunction := sync.OnceValue(func() uint64 {
		var maxSize uint64 = 0

		for value := range input {
			if value.Size > maxSize {
				maxSize = value.Size
			}
		}
		return maxSize
	})
	return resultFunction()
}
