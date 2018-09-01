// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:35<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package depta

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func create_2d_matrix(x, y int) [][]float64 {
	result := make([][]float64, x)
	for i := 0; i < x; i++ {
		result[i] = make([]float64, y)
	}
	return result
}

func DeepCopyRecord(el Record) Record {
	var mod bytes.Buffer
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)

	err := enc.Encode(el)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	var el_result Record
	err = dec.Decode(&el_result)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return el_result
}
