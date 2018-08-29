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
	"sort"
)

const (
	almost_similar = 0.8
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

func (r *Record) ConvertToBase() []*DTree {
	return *r
}

func (r *Record) id() string {
	var buffer bytes.Buffer

	for _, v := range *r {
		buffer.WriteString(v.identity())
	}
	return buffer.String()
}

func (r *Record) deepcopy() *Record {
	result := Record{}
	for _, v := range *r {
		el := v.deepcopy(DTree{})
		result = append(result, &el)
	}
	return &result
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
	stm                   SimpleTreeMatch
}

func (m *MiningDataRegion) init(root *DTree, max_generalized_nodes int, threshold float64) {
	m.root = root
	m.max_generalized_nodes = max_generalized_nodes
	m.threshold = threshold
}

func pairwise(data []*DTree, K, start int) [][][]*DTree {
	// TODO: check index sizes
	result := make([][][]*DTree, 0)
	// _ = "breakpoint"
	for k := 1; k < K+1; k++ {
		for i := 0; i < K; i++ {
			for j := start + i; j < len(data); j += k {
				slice_ax, slice_ay := bind_slice(j, j+k, len(data))
				slice_bx, slice_by := bind_slice(j+k, j+2*k, len(data))

				slice_a := data[slice_ax:slice_ay]
				slice_b := data[slice_bx:slice_by]
				if len(slice_a) >= k && len(slice_b) >= k {
					result = append(result, [][]*DTree{slice_a, slice_b})
				}
			}
		}
	}
	return result
}

func (m *MiningDataRegion) compare_generalized_nodes(parent *DTree, max_generalized_nodes int) map[string]GeneralizedNodeCompareContainer {
	// POTENTIAL BUGS. TEST

	//   compare the adjacent children generalized nodes similarity of a given element
	// Arguments:
	// `parent`: the lxml element to compare children of.
	// `k`: the maximum length of generalized node.
	scores := make(map[string]GeneralizedNodeCompareContainer, 0)
	for _, v := range pairwise(parent.Children, max_generalized_nodes, 0) {
		gn1 := GeneralizedNode{element: v[0][0], length: len(v[0])}
		gn2 := GeneralizedNode{element: v[1][0], length: len(v[1])}
		appender := GeneralizedNodeCompareContainer{left: gn1, right: gn2}
		_, ok := scores[appender.hash()]
		if ok == false {
			score := m.stm.normalized_match_score(v[0], v[1])
			appender.score = score
			scores[appender.hash()] = appender
		}
	}
	return scores
}

func (m *MiningDataRegion) calculate_score(region DataRegion) float64 {
	if region.Covered == 0 {
		return 0
	}
	count := region.Covered / region.K
	return region.Score / float64(count)
}

func (m *MiningDataRegion) identify_regions(start int,
	root *DTree,
	max_generalized_nodes int,
	threshold float64,
	scores map[string]GeneralizedNodeCompareContainer) []*DataRegion {

	cur_region := DataRegion{}
	max_region := DataRegion{}
	cur_region.init(root, 0, 0, 0, 0)
	max_region.init(root, 0, 0, 0, 0)
	data_regions := make([]*DataRegion, 0)

	for k := 1; k < max_generalized_nodes+1; k++ {
		for i := 0; i < max_generalized_nodes; i++ {
			flag := true
			score := 0.0
			for j := start + i; j < len(root.Children)-k; j += k {
				child_j, ok1 := root.get_child(j)
				child_jk, ok2 := root.get_child(j + k)
				if ok1 && ok2 {
					g1 := GeneralizedNode{element: child_j, length: k}
					g2 := GeneralizedNode{element: child_jk, length: k}
					container := GeneralizedNodeCompareContainer{left: g1, right: g2}
					score_item, _ := scores[container.hash()]
					score = score_item.score
				} else {
					score = 0
				}
				if score >= threshold {
					if flag {
						cur_region.K = k
						cur_region.Start = j
						cur_region.Score = score
						cur_region.Covered = 2 * k
						flag = false
					} else {
						cur_region.Covered += k
						cur_region.Score += score
					}
				} else if !flag {
					break
				}
			}
			if m.calculate_score(cur_region) > m.calculate_score(max_region) {
				max_region.K = cur_region.K
				max_region.Start = cur_region.Start
				max_region.Covered = cur_region.Covered
				max_region.Score = cur_region.Score
			}
		}
	}
	if max_region.Covered > 0 {
		data_regions = append(data_regions, &max_region)
		if max_region.Start+max_region.Covered < len(max_region.Parent.Children) {
			data_regions = append(data_regions, m.identify_regions(max_region.Start+max_region.Covered, root, max_generalized_nodes, threshold, scores)...)
		}
	}
	return data_regions
}

func (m *MiningDataRegion) find_regions(root *DTree) []*DataRegion {
	data_regions := make([]*DataRegion, 0)
	if root.tree_depth() >= 2 {
		scores := m.compare_generalized_nodes(root, m.max_generalized_nodes)
		data_regions = append(data_regions, m.identify_regions(0, root, m.max_generalized_nodes, m.threshold, scores)...)
		covered := make(map[string]*DTree, 0)

		for _, data_region := range data_regions {
			// all items that are covered by this data_region
			for i := data_region.Start; i < data_region.Covered; i++ {
				ichild, ok := data_region.Parent.get_child(i)
				if ok {
					covered[ichild.hash()] = ichild
				}
			}
		}

		// for each child that is not covered by that data region
		for _, child := range root.Children {

			_, ok := covered[child.hash()]
			if ok == false {
				data_regions = append(data_regions, m.find_regions(child)...)
			}
		}
	}
	return data_regions
}

type MiningDataRecord struct {
	// mining the data record from a region.
	// basic assumption:
	// the subtree of data records also similar. so if not any adjacent pair of them are
	// similar (less than threshold), data region itself is a data record,
	// otherwise children are individual data record.
	stm       SimpleTreeMatch
	threshold float64
}

func (m *MiningDataRecord) init(threshold float64) {
	m.stm = SimpleTreeMatch{}
	m.threshold = threshold
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
	return (float64(len(sims)) / float64(len(similarities))) > almost_similar
}

func (m *MiningDataRecord) find_records(region *DataRegion) []*Record {
	// PROBABLY BUGS, TEST
	if region.K == 1 {
		records := make([]*Record, 0)
		//if all the individual nodes of children nodes of Generalized node are similar
		for i := region.Start; i < region.Start+region.Covered; i++ {
			child, ok := region.Parent.get_child(i)
			if ok {
				for _, children := range pairwise(child.Children, 1, 0) {
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

type MiningDataField struct {
	pta PartialTreeAligner
	sta SimpleTreeAligner
}

func (m *MiningDataField) init() {
	m.sta = SimpleTreeAligner{}
	m.pta = PartialTreeAligner{}
}

func (m *MiningDataField) align_records(records []*Record) ([]*Record, *Record) {
	// partial align multiple records.

	// sort by the tree size
	srecords := records
	sort.Sort(BySize(records))
	// seed is the largest tree
	seed, srecords := srecords[len(srecords)-1], srecords[:len(srecords)-1]

	R := make([]*Record, 0)
	items := make([]*Record, 0)
	mappings := make(map[string]map[string]*DTree, 0)
	seed_copy := seed.deepcopy()
	mappings[seed.id()] = m._create_seed_mapping(seed_copy, seed)
	// mappings.setdefault(seed, self._create_seed_mapping(seed_copy, seed))
	var next *Record
	for len(srecords) > 0 {
		next, srecords = srecords[len(srecords)-1], srecords[:len(srecords)-1]
		modified, partial_match, aligned := m.pta.align([]*DTree(*seed_copy), []*DTree(*next))
		mappings[next.id()] = aligned
		if modified {
			records = append(records, R...)
			R = make([]*Record, 0)
		} else {
			if partial_match {
				R = append(R, next)
			}
		}
	}
	for _, record := range records {
		aligned := mappings[record.id()]
		result := m._extract_item(seed_copy, aligned)
		items = append(items, result)
	}
	return items, seed_copy
}
func (m *MiningDataField) align_record(seed, record *Record) *Record {
	// simple align the given record with given seed
	alignment := m.sta.align([]*DTree(*seed), []*DTree(*record))
	aligned := make(map[string]*DTree, 0)
	aligned[alignment.first.identity()] = alignment.second
	for _, sub := range alignment.subs {
		aligned[sub.first.identity()] = sub.second
	}
	return m._extract_item(seed, aligned)

}
func (m *MiningDataField) _extract_item(seed *Record, mapping map[string]*DTree) *Record {
	// extract data item from the tree.
	// r = []
	// for element in seed:
	//     r.extend(self._extract_element(element, mapping))
	// return r
	r := Record{}
	for _, element := range []*DTree(*seed) {
		r = append(r, *m._extract_element(element, mapping)...)
	}
	return &r
}

func (m *MiningDataField) _extract_element(seed *DTree, mapping map[string]*DTree) *Record {
	// actually, this is a simple flattening function
	r := Record{}
	e, ok := mapping[seed.identity()]
	if ok {
		r = append(r, e)
		if len(e.Children) > 0 {
			for _, child := range e.Children {
				r = append(r, *m._extract_element(child, mapping)...)
			}
		}
	}
	return &r
}

func (m *MiningDataField) _create_seed_mapping(seed, record *Record) map[string]*DTree {
	// create a mapping from seed record to data record.
	// for s, e in zip(seed, record):
	//     d[s] = e
	//     d.update(self._create_seed_mapping(s, e))
	// return d
	result := make(map[string]*DTree, 0)
	// log.Println(seed)
	seed_ := *record
	record_ := *seed
	for k := 0; k <= len(*seed)-1 && k <= len(*record)-1; k++ {
		e := record_[k]
		s := seed_[k]
		result[s.identity()] = e
		chh := Record(s.Children)
		che := Record(e.Children)
		for k1, v1 := range m._create_seed_mapping(&chh, &che) {
			result[k1] = v1
		}
	}
	return result
}

func (d *DataRegion) as_html_table(show_id bool) string {
	// convert region to a HTML table
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
	// convert the region to a two dim plain texts.
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
