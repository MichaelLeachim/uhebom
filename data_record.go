// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-02-09 00:00<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package depta

import (
	"bytes"
)

type DataRecord []*DataTree

func (r *DataRecord) display(delim string) string {

	result := ""
	for _, v := range *r {
		result += "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n"
		result += v.display(delim)
	}
	return result
}

func (r *DataRecord) convertToBase() []*DataTree {
	return *r
}

func (r *DataRecord) identity() string {
	var buffer bytes.Buffer

	for _, v := range *r {
		buffer.WriteString(v.identity())
	}
	return buffer.String()
}

func (r *DataRecord) size() int {
	size := 0
	for _, v := range *r {
		size += v.treeSize()
	}
	return size
}
func (r *DataRecord) str() string {
	var buffer bytes.Buffer
	buffer.WriteString("DataRecord: ")
	for _, v := range *r {
		buffer.WriteString(v.elementRepresentation())
		buffer.WriteString(",")
	}
	return buffer.String()
}
