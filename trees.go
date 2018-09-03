// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:34<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
package depta

import ()

type DataTree struct {
	Tag      string
	Data     string
	Children []*DataTree
	_id      string
	Index    int
	Parent   *DataTree
	Attrs    map[string]string
}

func (d *DataTree) display(delim string) string {
	result := ""
	result += utils.joinString(delim, d.Tag, ":", d.Data, "\n")
	for _, v := range d.Children {
		result += v.display(delim + delim)
	}
	return result
}

func (d *DataTree) identity() string {
	if d._id == "" {
		d._id = utils.makeIdString()
	}
	return d._id
}

func (d *DataTree) childAppend(child *DataTree) {
	last_index := len(d.Children) - 1
	if last_index < 0 {
		last_index = 0
	}
	child.Parent = d
	child.Index = last_index
	d.Children = append(d.Children, child)
}

func (d *DataTree) elementRepresentation() string {
	return utils.joinString("<", d.Tag, "#", d.Attrs["id"], ".", d.Attrs["class"], ">")
}

func (d *DataTree) hash() string {
	return string(d.identity())
}

func (d *DataTree) treeSize() int {
	// tree_size(root)
	if len(d.Children) == 0 {
		return 1
	}

	size := 0
	for _, v := range d.Children {
		size += v.treeSize()
	}
	return size + 1
}

func (d *DataTree) treeDepth() int {
	// tree_depth(root)
	if len(d.Children) == 0 {
		return 1
	}
	size := 0

	for _, v := range d.Children {
		tsize := v.treeSize()
		if tsize > size {
			size = tsize
		}
	}
	return size + 1
}

func (d *DataTree) getRoot() (string, bool) {
	// determines node uniqueness for alignment
	// for the future, add another tree metric
	if d.Tag == "" {
		return "", false
	}

	return d.Tag, true
}

func (d *DataTree) getChild(i int) (*DataTree, bool) {
	// _get_child(e,i)
	if i < len(d.Children) {
		return d.Children[i], true
	}
	return &DataTree{}, false
}

func (d *DataTree) getChildrenCount() int {
	// _get_child(e,i)
	return len(d.Children)
}
