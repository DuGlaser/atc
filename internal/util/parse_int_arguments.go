package util

import (
	"fmt"
	"strconv"
)

func ParseIntArguments(s []string) ([]int, error) {
	is := []int{}

	visited := map[int]bool{}

	for _, si := range s {
		i, err := strconv.ParseInt(si, 10, 8)
		if err != nil {
			return nil, err
		}

		ii := int(i)
		if visited[ii] {
			return nil, fmt.Errorf("Duplicate value %d", ii)
		}

		visited[ii] = true
		is = append(is, ii)
	}

	return is, nil
}
