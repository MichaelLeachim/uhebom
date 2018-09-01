// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:34<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
package depta

import (
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

	// We are inserting item at a point (with shifting)
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
