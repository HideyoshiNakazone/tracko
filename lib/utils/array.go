package utils

import "iter"


func PartitionChannel[T any](items <-chan T, size int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		batch := make([]T, 0, size)

		for item := range items {
			batch = append(batch, item)
			if len(batch) == size {
				if !yield(batch) {
					return
				}
				batch = make([]T, 0, size)
			}
		}

		if len(batch) > 0 {
			yield(batch)
		}
	}
}

func Map[T any, R any](input []T, mapper func(T) R) []R {
    output := make([]R, len(input))
    for i, v := range input {
        output[i] = mapper(v)
    }
    return output
}
