// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:34<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package depta

import (
	"github.com/satori/go.uuid"
	"hash/fnv"
	"log"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func make_id() []byte {
	id := uuid.NewV1()
	log.Println(string(id.String()))
	return id.Bytes()
}
func make_id_string() string {
	return uuid.NewV1().String()
}

func maxf(data []float64) float64 {
	max := 0.0
	for _, value := range data {
		if value > max {
			max = value // found another, bigger value, replace previous value in max
		}
	}
	return max
}

func bind_slice(start, end, size int) (int, int) {
	// will return appropriate starts and ends based on list size
	if start < 0 || end < 0 {
		panic("Should not use negative slices. It is not a Python")
	}
	if size == 0 {
		return 0, 0
	}

	if start > size {
		start = size
	}
	if end > size {
		end = size
	}
	if start > end {
		panic("End larger than start")
	}
	return start, end
}

func sumi(data []int) int {
	sum := 0
	for _, value := range data {
		sum += value
	}
	return sum
}
