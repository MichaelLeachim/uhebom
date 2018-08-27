// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:34<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
package depta

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type DTree struct {
	Tag      string
	Data     string
	Children []*DTree
	_id      string
	Index    int
	Parent   *DTree
	Attrs    map[string]string
}

func (d *DTree) display(delim string) string {
	result := ""
	result += fmt.Sprintf("%v%v:%v\n", delim, d.Tag, d.Data)
	for _, v := range d.Children {
		result += v.display(delim + delim)
	}
	return result
}
func (d *DTree) display_text() string {
	result := d.Data
	for _, v := range d.Children {
		result += " " + v.display_text()
	}
	return result
}

func (d *DTree) identity() string {
	if d._id == "" {
		d._id = make_id_string()
	}
	return d._id
}

func (d *DTree) flatten() []string {
	// flattens children tree
	result := make([]string, 0)
	result = append(result, d.Data)
	for _, v := range d.Children {
		result = append(result, v.flatten()...)
	}
	return result
}

func (d *DTree) getprevious() (*DTree, bool) {
	if d.Index-1 < 0 {
		return &DTree{}, false
	}
	return d.Parent.Children[d.Index-1], true
}

func (d *DTree) getnext() (*DTree, bool) {
	last_child_index := len(d.Parent.Children) - 1
	if d.Index >= last_child_index {
		return &DTree{}, false
	} else {
		return d.Parent.Children[d.Index+1], true
	}
}

func (d *DTree) deepcopy(parent DTree) DTree {
	// TODO: test it
	element := *d
	element.Parent = &parent
	children := make([]*DTree, 0)
	for _, v := range d.Children {
		child := v.deepcopy(element)
		children = append(children, &child)
	}
	element.Children = children
	return element
}

func (d *DTree) child_insert(child *DTree, index int) {
	// if index > len(d.Children)-1
	//    index = len(d.Children)-1
	// if index < 0:
	//    index = 0
	log.Println("OOOOOOOOOOOOOOOOOO")
	if index > len(d.Children) {
		index = len(d.Children) - 1
	}
	if index < 0 {
		index = 0
	}

	child.Parent = d
	child.Index = index

	// a little bit of mumbo jumbo.
	// We are inserting item at point(with shifting)
	result := make([]*DTree, 0)
	for k, v := range d.Children {
		if k == index {
			result = append(result, child)
		}
		result = append(result, v)
	}
	d.Children = result
	child.identity()
}

func (d *DTree) child_append(child *DTree) {
	last_index := len(d.Children) - 1
	if last_index < 0 {
		last_index = 0
	}
	child.Parent = d
	child.Index = last_index
	d.Children = append(d.Children, child)
}

func (d *DTree) parent_add(parent *DTree) {
	parent.child_append(d)
}

func (d *DTree) element_repr() string {
	return fmt.Sprintf("<%s #%s .%s>", d.Tag, d.Attrs["class"], d.Attrs["id"])
}

func (d *DTree) nilp() bool {
	if d.Tag == "" && d.Data == "" && len(d.Children) == 0 {
		return true
	}
	return false
}

func (d *DTree) hash() string {
	return string(d.identity())
}

func (d *DTree) tree_size() int {
	// tree_size(root)
	if len(d.Children) == 0 {
		return 1
	}

	size := 0
	for _, v := range d.Children {
		size += v.tree_size()
	}
	return size + 1
}

func (d *DTree) tree_depth() int {
	// tree_depth(root)
	if len(d.Children) == 0 {
		return 1
	}
	size := 0

	for _, v := range d.Children {
		tsize := v.tree_size()
		if tsize > size {
			size = tsize
		}
	}
	return size + 1
}
func (d *DTree) get_root() (string, bool) {
	// determines node uniqueness for alignment
	// for the future, add another tree metrics
	if d.Tag == "" {
		return "", false
	}

	return d.Tag, true
}

func (d *DTree) get_child(i int) (*DTree, bool) {
	// _get_child(e,i)
	if i < len(d.Children) {
		return d.Children[i], true
	}
	return &DTree{}, false
}

func (d *DTree) get_children_count() int {
	// _get_child(e,i)
	return len(d.Children)
}

type SimpleTreeMatch struct{}

func (s *SimpleTreeMatch) match(l1, l2 []*DTree) float64 {
	// match(self, l1, l2)
	// match two trees list.
	rows := len(l1) + 1
	cols := len(l2) + 1
	m := create_2d_matrix(rows, cols)
	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			m[i][j] = maxf([]float64{m[i][j], m[i-1][j-1] + s._single_match(l1[i-1], l2[j-1])})
		}
	}
	return m[rows-1][cols-1]
}

func (s *SimpleTreeMatch) normalized_match_score(l1, l2 []*DTree) float64 {
	l1size := 1
	l2size := 1
	for _, v := range l1 {
		l1size += v.tree_size()
	}
	for _, v := range l2 {
		l2size += v.tree_size()
	}

	return float64(s.match(l1, l2)) / (float64(l1size+l2size) / float64(2))
}

// probably, meaningless function
func (s *SimpleTreeMatch) _single_match(t1, t2 *DTree) float64 {
	return tree_match(t1, t2)
}

type TreeAlignment struct {
	TRACE_LEFT float64
	TRACE_UP   float64
	TRACE_DIAG float64
	score      float64
	subs       []*TreeAlignment
	first      *DTree
	second     *DTree
}

func (s *TreeAlignment) init(first, second *DTree, score float64) {
	// first=None, second=None, score=0

	s.first = first
	s.second = second
	s.score = score
	s.subs = make([]*TreeAlignment, 0)
	s.TRACE_LEFT = 1
	s.TRACE_UP = 2
	s.TRACE_DIAG = 3
}

func (s *TreeAlignment) add(alignment *TreeAlignment) {
	if s.first.nilp() && s.second.nilp() {
		s.first = alignment.first
		s.second = alignment.second
	} else {
		s.subs = append(s.subs, alignment)
	}
	s.subs = append(s.subs, alignment.subs...)
}

func (s *TreeAlignment) tag() string {
	if s.first.Tag != s.second.Tag {
		log.Fatal(errors.New(s.first.Tag + "!=" + s.second.Tag))
	}
	return s.first.Tag
}

type SimpleTreeAligner struct{}

func (s *SimpleTreeAligner) align(l1, l2 []*DTree) *TreeAlignment {
	alignment := TreeAlignment{}

	alignment.init(&DTree{}, &DTree{}, 0)
	matrix := create_2d_matrix(len(l1)+1, len(l2)+1)
	alignment_matrix := create_2d_matrix_TreeAlignment(len(l1), len(l2))

	trace := create_2d_matrix(len(l1), len(l2))
	for i := 1; i < len(matrix); i++ {
		for j := 1; j < len(matrix[0]); j++ {
			if matrix[i][j-1] > matrix[i-1][j] {
				matrix[i][j] = matrix[i][j-1]
				trace[i-1][j-1] = float64(alignment.TRACE_LEFT)
			} else {
				matrix[i][j] = matrix[i-1][j]
				trace[i-1][j-1] = float64(alignment.TRACE_UP)
			}
			alignment_matrix[i-1][j-1] = s.single_align(l1[i-1], l2[j-1])
			score := float64(matrix[i-1][j-1]) + alignment_matrix[i-1][j-1].score
			if score > matrix[i][j] {
				matrix[i][j] = score
				trace[i-1][j-1] = float64(alignment.TRACE_DIAG)
			}
		}
	}
	row := len(trace) - 1
	col := len(trace[0]) - 1

	for row >= 0 && col >= 0 {
		if trace[row][col] == alignment.TRACE_DIAG {
			alignment.add(alignment_matrix[row][col])
			row -= 1
			col -= 1
		} else if trace[row][col] == alignment.TRACE_UP {
			row -= 1

		} else if trace[row][col] == alignment.TRACE_LEFT {
			col -= 1
		}
	}
	alignment.score = 0
	alignment.score += matrix[len(matrix)-1][len(matrix[0])-1]
	return &alignment
}

func (s *SimpleTreeAligner) single_align(t1, t2 *DTree) *TreeAlignment {
	// differs
	t1root, t1root_ok := t1.get_root()
	t2root, t2root_ok := t2.get_root()
	if (!t1root_ok || !t2root_ok) || t1root != t2root {
		alignment := TreeAlignment{}
		alignment.init(&DTree{}, &DTree{}, 0)
		return &alignment
	}
	alignment := TreeAlignment{}
	alignment.init(t1, t2, 0)

	t1_len := t1.get_children_count()
	t2_len := t2.get_children_count()
	matrix := create_2d_matrix(t1_len+1, t2_len+1)
	alignment_matrix := create_2d_matrix_TreeAlignment(t1_len, t2_len)
	trace := create_2d_matrix(t1_len, t2_len)
	// end differs

	for i := 1; i < len(matrix); i++ {
		for j := 1; j < len(matrix[0]); j++ {
			if matrix[i][j-1] > matrix[i-1][j] {
				matrix[i][j] = matrix[i][j-1]
				trace[i-1][j-1] = alignment.TRACE_LEFT
			} else {
				matrix[i][j] = matrix[i-1][j]
				trace[i-1][j-1] = alignment.TRACE_UP
			}

			// differs
			child1, _ := t1.get_child(i - 1)
			child2, _ := t2.get_child(j - 1)
			alignment_matrix[i-1][j-1] = s.single_align(child1, child2)
			// end differs

			score := matrix[i-1][j-1] + alignment_matrix[i-1][j-1].score
			if score > matrix[i][j] {
				matrix[i][j] = score
				trace[i-1][j-1] = alignment.TRACE_DIAG
			}
		}
	}
	// differs
	row := len(trace) - 1
	col := -1
	if row >= 0 {
		col = len(trace[0]) - 1
	} else {
		col = -1
	}
	// end differs
	for row >= 0 && col >= 0 {
		if trace[row][col] == alignment.TRACE_DIAG {
			alignment.add(alignment_matrix[row][col])
			row -= 1
			col -= 1
		} else if trace[row][col] == alignment.TRACE_UP {
			row -= 1

		} else if trace[row][col] == alignment.TRACE_LEFT {
			col -= 1
		}
	}
	// differs
	alignment.score = 1
	// end differs
	alignment.score += matrix[len(matrix)-1][len(matrix[0])-1]
	return &alignment
}

func find_subsequence(iterable []*DTree, predicate func(*DTree) bool) [][]*DTree {

	seqs := make([][]*DTree, 0)
	seq := make([]*DTree, 0)
	continuous := false
	for _, i := range iterable {
		if predicate(i) {
			if continuous {
				seq = append(seq, i)
			} else {
				seq = []*DTree{i}
				continuous = true
			}
		} else if continuous {
			seqs = append(seqs, seq)
			seq = make([]*DTree, 0)
			continuous = false
		}
	}

	if len(seq) > 0 {
		seqs = append(seqs, seq)
	}
	return seqs
}

type PartialTreeAligner struct {
	sta     SimpleTreeAligner
	options map[string]string
}

func (p *PartialTreeAligner) find_unaligned_elements(aligned, reverse_aligned map[string]*DTree, elements []*DTree) [][]*DTree {
	unaligned := make([][]*DTree, 0)
	predicate := func(d *DTree) bool {
		_, ok := reverse_aligned[d.identity()]
		return !ok
	}

	for _, element := range elements {
		current_level := find_subsequence(element.Children, predicate)
		unaligned = append(unaligned, current_level...)
		for _, child := range element.Children {
			unaligned = append(unaligned, find_subsequence(child.Children, predicate)...)
		}
	}
	return unaligned
}

func (p *PartialTreeAligner) align(l1, l2 []*DTree) (bool, bool, map[string]*DTree) {

	alignment := p.sta.align(l1, l2)

	aligned := make(map[string]*DTree, 0)
	reverse_aligned := make(map[string]*DTree, 0)
	modified := false

	aligned[alignment.first.identity()] = alignment.second
	reverse_aligned[alignment.second.identity()] = alignment.first
	for _, sub := range alignment.subs {
		aligned[sub.first.identity()] = sub.second
		aligned[sub.second.identity()] = sub.first
	}

	unaligned_elements := p.find_unaligned_elements(aligned, reverse_aligned, l2)
	for _, l := range unaligned_elements {
		left_most := l[0]
		right_most := l[len(l)-1]

		prev_sibling, prev_ok := left_most.getprevious()
		next_sibling, next_ok := right_most.getnext()

		if !prev_ok {
			if next_ok {
				// leftmost alignment
				next_sibling_match, ok := reverse_aligned[next_sibling.identity()]
				if ok {
					for i, element := range l {
						element_copy := element.deepcopy(DTree{})
						next_sibling_match.Parent.child_insert(&element_copy, i)
						aligned[element_copy.identity()] = element
					}
					modified = true
				}
			}
		} else if !next_ok {
			//  rightmost alignment
			prev_sibling_match := reverse_aligned[prev_sibling.identity()]
			previous_match_index := prev_sibling_match.Index
			// unique insertion
			for i, element := range l {
				element_copy := element.deepcopy(DTree{})
				prev_sibling_match.Parent.child_insert(&element_copy, previous_match_index+1+i)
				aligned[element_copy.identity()] = element
			}
			modified = true
		} else {
			// flanked by two sibling elements
			prev_sibling_match := reverse_aligned[prev_sibling.identity()]
			next_sibling_match := reverse_aligned[next_sibling.identity()]
			if !prev_sibling_match.nilp() && !next_sibling_match.nilp() {
				next_match_index := next_sibling_match.Index
				previous_match_index := prev_sibling_match.Index
				if (next_match_index - previous_match_index) == 1 {
					//unique insertion
					for i, element := range l {
						element_copy := element.deepcopy(DTree{})
						prev_sibling_match.Parent.child_insert(&element_copy, previous_match_index+1+i)
						aligned[element_copy.identity()] = element
					}
					modified = true
				}
			}
		}
	}

	return modified, len(unaligned_elements) > 0, aligned
}
