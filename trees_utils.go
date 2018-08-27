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
	log "github.com/sirupsen/logrus"
)

func as_html_tables(data []*DataRegion, show_id bool) []byte {
	var template bytes.Buffer
	template.WriteString("<style>table {border-collapse: collapse;}table, th, td {border: 1px solid black;}</style>")

	for i, item := range data {
		template.WriteString(fmt.Sprintf("<h1>Table number: %d </h1>", i))
		template.WriteString(item.as_html_table(show_id))
	}
	return template.Bytes()
}

func create_2d_matrix(x, y int) [][]float64 {
	result := make([][]float64, x)
	for i := 0; i < x; i++ {
		result[i] = make([]float64, y)
	}
	return result
}
func create_2d_matrix_TreeAlignment(x, y int) [][]*TreeAlignment {
	result := make([][]*TreeAlignment, x)
	for i := 0; i < x; i++ {
		result[i] = make([]*TreeAlignment, y)
	}
	return result
}

func tree_match(t1, t2 *DTree) float64 {
	t1_root, t1_exist := t1.get_root()
	t2_root, t2_exist := t2.get_root()
	if !t1_exist || !t2_exist {
		return 0
	}
	if t1_root != t2_root {
		return 0
	}
	rows := t1.get_children_count() + 1
	cols := t2.get_children_count() + 1
	m := create_2d_matrix(rows, cols)
	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			// m[i][j] = maxi([]int{m[i][j-1], m[i-1][j]}) // probably, meaningless line
			child1, _ := t1.get_child(i - 1)
			child2, _ := t2.get_child(j - 1)
			m[i][j] = maxf([]float64{m[i][j], m[i-1][j-1] + tree_match(child1, child2)})
		}
	}
	return 1 + m[rows-1][cols-1]
}

type BySize []*Record

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].size() > a[j].size() }

// func DeepCopyRecordSlice(el []Record) []Record {
// 	var mod bytes.Buffer
// 	enc := gob.NewEncoder(&mod)
// 	dec := gob.NewDecoder(&mod)

// 	err := enc.Encode(el)
// 	if err != nil {
// 		log.Fatal("encode error:", err)
// 	}
// 	el_result := make([]Record, 0)
// 	err = dec.Decode(&el_result)
// 	if err != nil {
// 		log.Fatal("decode error:", err)
// 	}
// 	return el_result
// }

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
