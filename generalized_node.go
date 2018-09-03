// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:34<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
package uhebom

type GeneralizedNode struct {
	element *DataTree
	length  int
}

type GeneralizedNodeCompareContainer struct {
	left  GeneralizedNode
	right GeneralizedNode
	score float64
}

// TODO: remove candidate
func (g GeneralizedNodeCompareContainer) identity() string {
	return g.left.element.identity() + g.right.element.identity()
}
