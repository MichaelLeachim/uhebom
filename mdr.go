// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:34<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
package depta

import (
	"bytes"
	"fmt"
)

const (
	ALMOST_SIMILAR = 0.8
)

type DataRegion struct {
	Parent  *DTree
	Start   int
	K       int
	Covered int
	Score   float64
	Items   []*Record
}

func (d *DataRegion) init(parent *DTree, start, k, covered int, score float64) {
	d.Parent = parent
	d.Start = start
	d.K = k
	d.Score = score
	d.Covered = covered
}

func (d *DataRegion) iter(k int) [][]*DTree {
	result := make([][]*DTree, 0)

	for i := d.Start; i < d.Start+d.Covered; i += k {
		result = append(result, d.Parent.Children[i:i+k])
	}
	return result
}

type Record []*DTree

func (r *Record) Display(delim string) string {
	result := ""
	for _, v := range *r {
		result += "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n"
		result += v.display(delim)
	}
	return result
}

func (r *Record) convert_to_base() []*DTree {
	return *r
}

func (r *Record) id() string {
	var buffer bytes.Buffer

	for _, v := range *r {
		buffer.WriteString(v.identity())
	}
	return buffer.String()
}

func (r *Record) size() int {
	size := 0
	for _, v := range *r {
		size += v.tree_size()
	}
	return size
}
func (r *Record) str() string {
	var buffer bytes.Buffer
	buffer.WriteString("DataRecord: ")
	for _, v := range *r {
		buffer.WriteString(v.element_repr())
		buffer.WriteString(",")
	}
	return buffer.String()
}

type GeneralizedNode struct {
	element *DTree
	length  int
}

type GeneralizedNodeCompareContainer struct {
	left  GeneralizedNode
	right GeneralizedNode
	score float64
}

func (g GeneralizedNodeCompareContainer) hash() string {
	// oh. gosh, just kill me now..
	// The problem is, when you instantiate element, you've got new address.
	// So, instead of deep comparison of newly created elements,
	// we just build hash from a string sum of pointers.
	// In case, you know how to do better. Welcome!
	return fmt.Sprintf("%v%v", g.left.element.identity(), g.right.element.identity())
}

type MiningDataRegion struct {
	root                  *DTree
	max_generalized_nodes int
	threshold             float64
}

func (m *MiningDataRegion) init(root *DTree, max_generalized_nodes int, threshold float64) {
	m.root = root
	m.max_generalized_nodes = max_generalized_nodes
	m.threshold = threshold
}

type MiningDataRecord struct {
	// mining the data record from a region.
	// basic assumption:
	// the subtree of data records also similar. so if not any adjacent pair of them are
	// similar (less than threshold), data region itself is a data record,
	// otherwise children are individual data record.
	threshold float64
	stm       *SimpleTreeMatch
}

func newMiningDataRecord(threshold float64) *MiningDataRecord {
	m := MiningDataRecord{}
	m.threshold = threshold
	m.stm = newSimpleTreeMatch()
	return &m
}

func (m *MiningDataRecord) slice_region(region *DataRegion) []*Record {
	// slice every generalized node of region to a data records
	records := make([]*Record, 0)
	for _, v := range region.iter(region.K) {
		vv := Record(v)
		records = append(records, &vv)
	}
	return records
}

func (m *MiningDataRecord) almost_similar(similarities map[string]float64, threshold float64) bool {
	sims := make([]int, 0)

	for _, v := range similarities {
		if float64(v) >= threshold {
			sims = append(sims, 1)
		}
	}
	return (float64(len(sims)) / float64(len(similarities))) > ALMOST_SIMILAR
}

func (m *MiningDataRecord) find_records(region *DataRegion) []*Record {
	// PROBABLY BUGS, TEST
	if region.K == 1 {
		records := make([]*Record, 0)
		//if all the individual nodes of children nodes of Generalized node are similar
		for i := region.Start; i < region.Start+region.Covered; i++ {
			child, ok := region.Parent.get_child(i)
			if ok {
				for _, children := range trees_utils.Pairwise(child.Children, 1, 0) {
					sim := m.stm.normalized_match_score(children[0], children[1])
					if float64(sim) < m.threshold {
						return m.slice_region(region)
					}
				}
			}
		}
		// each child of generalized node is a data record
		for _, gn := range region.iter(1) {
			for _, c := range gn {
				cc := Record{c}
				records = append(records, &cc)
			}
		}
		return records
	} else {
		// if almost all the individual nodes in Generalized Node are similar
		children := make([]*DTree, 0)
		for i := 0; i < region.Covered; i++ {
			child, ok := region.Parent.get_child(region.Start + i)
			if ok {
				children = append(children, child)
			}

		}
		sizes := make(map[int]int, 0)
		for _, child := range children {
			sizes[child.tree_size()] += 1
		}
		most_common_size_counter := 0
		most_common_size := 0
		for k, v := range sizes {
			if v > most_common_size_counter {
				most_common_size_counter = v
				most_common_size = k
			}
		}
		var most_typical_child *DTree
		for _, v := range children {
			if v.tree_size() == most_common_size {
				most_typical_child = v
			}
		}
		similarities := make(map[string]float64, 0)
		for _, child := range children {
			similarities[child.hash()] = m.stm.normalized_match_score([]*DTree{most_typical_child}, []*DTree{child})
		}
		result := make([]*Record, 0)
		if m.almost_similar(similarities, m.threshold) {
			for _, child := range children {
				if float64(similarities[child.hash()]) >= m.threshold {
					rr := Record([]*DTree{child})
					result = append(result, &rr)
				}
			}
			return result
		} else {
			return m.slice_region(region)
		}
	}
}

func (d *DataRegion) as_html_table(show_id bool) string {
	var buffer bytes.Buffer
	buffer.WriteString("<table>")
	for i, item := range d.Items {
		buffer.WriteString("<tr>")
		if show_id {
			buffer.WriteString(fmt.Sprintf("<td>%d</td>", i+1))
		}
		for _, field := range *item {
			if field.Tag == "a" {
				buffer.WriteString(fmt.Sprintf("<td><a href='%s'>%s</a></td>", field.Data, field.Data))
			} else if field.Tag == "img" {
				buffer.WriteString(fmt.Sprintf("<td><img src='%s'></img></td>", field.Data))
			} else {
				buffer.WriteString(fmt.Sprintf("<td>%s</td>", field.Data))
			}
		}
		buffer.WriteString("</tr>")
	}
	buffer.WriteString("</table>")
	return buffer.String()
}

func (d *DataRegion) AsPlainTexts() [][]string {
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
