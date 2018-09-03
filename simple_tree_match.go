// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-01-09 23:04<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package uhebom

type SimpleTreeMatch struct{}

func newSimpleTreeMatch() *SimpleTreeMatch {
	return &SimpleTreeMatch{}
}

func (s *SimpleTreeMatch) match(l1, l2 []*DataTree) float64 {
	// match(self, l1, l2)
	// match two trees list.
	rows := len(l1) + 1
	cols := len(l2) + 1
	m := trees_utils.create2dMatrix(rows, cols)
	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			m[i][j] = utils.maxf([]float64{m[i][j], m[i-1][j-1] + trees_utils.treeMatch(l1[i-1], l2[j-1])})
		}
	}
	return m[rows-1][cols-1]
}

func (s *SimpleTreeMatch) normalizedMatchScore(l1, l2 []*DataTree) float64 {
	l1size := 1
	l2size := 1
	for _, v := range l1 {
		l1size += v.treeSize()
	}
	for _, v := range l2 {
		l2size += v.treeSize()
	}
	return float64(s.match(l1, l2)) / (float64(l1size+l2size) / float64(2))
}
