// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:34<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package uhebom

import (
	"bytes"
	"github.com/satori/go.uuid"
	"hash/fnv"
)

type utils_ struct{}

var utils = utils_{}

func (u *utils_) joinString(data ...string) string {
	var template bytes.Buffer
	for _, v := range data {
		template.WriteString(v)
	}
	return template.String()
}

func (u *utils_) tabularFormToBase(datum [][][]*TabularForm) {

}

func (u *utils_) hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (u *utils_) makeId() []byte {
	id, _ := uuid.NewV1()
	return id.Bytes()
}
func (u *utils_) makeIdString() string {
	uuid, _ := uuid.NewV1()
	return uuid.String()
}

func (u *utils_) maxf(data []float64) float64 {
	max := 0.0
	for _, value := range data {
		if value > max {
			max = value // found another, bigger value, replace previous value in max
		}
	}
	return max
}

func (u *utils_) bindSlice(start, end, size int) (int, int) {
	// will return appropriate starts and ends based on list size
	if start < 0 || end < 0 {
		panic("Should not use negative slices. It is not Python")
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

func (u *utils_) sumi(data []int) int {
	sum := 0
	for _, value := range data {
		sum += value
	}
	return sum
}
