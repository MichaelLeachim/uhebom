// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-02-09 00:06<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package depta

import (
	"bytes"
	"strconv"
)

const (
	ALMOST_SIMILAR = 0.8
)

type DataRegion struct {
	Parent  *DataTree
	Start   int
	K       int
	Covered int
	Score   float64
	Items   []*DataRecord
}

func newDataRegion(parent *DataTree, start, k, covered int, score float64) *DataRegion {
	d := DataRegion{}
	d.Parent = parent
	d.Start = start
	d.K = k
	d.Score = score
	d.Covered = covered
	return &d
}

func (d *DataRegion) iter(k int) [][]*DataTree {
	result := make([][]*DataTree, 0)

	for i := d.Start; i < d.Start+d.Covered; i += k {
		result = append(result, d.Parent.Children[i:i+k])
	}
	return result
}

func (d *DataRegion) asHTMLTable(show_id bool) string {
	var buffer bytes.Buffer
	buffer.WriteString("<table>")
	for i, item := range d.Items {
		buffer.WriteString("<tr>")
		if show_id {
			buffer.WriteString("<td>")
			buffer.WriteString(strconv.Itoa(i + 1))
			buffer.WriteString("</td>")
		}
		for _, field := range *item {
			buffer.WriteString("<td>")
			buffer.WriteString(field.Data)
			buffer.WriteString("</td>")
		}
		buffer.WriteString("</tr>")
	}
	buffer.WriteString("</table>")
	return buffer.String()
}

func (d *DataRegion) asPlainTexts() [][]string {
	result := make([][]string, 0)
	for _, v := range d.Items {
		sub_result := make([]string, 0)
		for _, v2 := range *v {
			sub_result = append(sub_result, v2.Data)
		}
		result = append(result, sub_result)
	}
	return result
}
