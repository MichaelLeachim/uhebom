// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-01-09 22:58<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package uhebom

type MiningDataRegion struct {
	root                  *DataTree
	max_generalized_nodes int
	threshold             float64
	stm                   SimpleTreeMatch
}

func newMiningDataRegion(root *DataTree, max_generalized_nodes int, threshold float64) *MiningDataRegion {
	m := MiningDataRegion{}
	m.root = root
	m.max_generalized_nodes = max_generalized_nodes
	m.threshold = threshold
	return &m
}

func (m *MiningDataRegion) compareGeneralizedNodes(parent *DataTree, max_generalized_nodes int) map[string]GeneralizedNodeCompareContainer {
	scores := make(map[string]GeneralizedNodeCompareContainer, 0)
	for _, v := range trees_utils.pairwise(parent.Children, max_generalized_nodes, 0) {
		gn1 := GeneralizedNode{element: v[0][0], length: len(v[0])}
		gn2 := GeneralizedNode{element: v[1][0], length: len(v[1])}
		appender := GeneralizedNodeCompareContainer{left: gn1, right: gn2}
		_, ok := scores[appender.identity()]
		if ok == false {
			score := m.stm.normalizedMatchScore(v[0], v[1])
			appender.score = score
			scores[appender.identity()] = appender
		}
	}
	return scores
}

func (m *MiningDataRegion) calculateScore(region *DataRegion) float64 {
	if region.Covered == 0 {
		return 0
	}
	count := region.Covered / region.K
	return region.Score / float64(count)
}

func (m *MiningDataRegion) identifyRegions(start int,
	root *DataTree, max_generalized_nodes int,
	threshold float64,
	scores map[string]GeneralizedNodeCompareContainer) []*DataRegion {

	cur_region := newDataRegion(root, 0, 0, 0, 0)
	max_region := newDataRegion(root, 0, 0, 0, 0)
	data_regions := []*DataRegion{}

	for k := 1; k < max_generalized_nodes+1; k++ {
		for i := 0; i < max_generalized_nodes; i++ {
			flag := true
			score := 0.0
			for j := start + i; j < len(root.Children)-k; j += k {
				child_j, ok1 := root.getChild(j)
				child_jk, ok2 := root.getChild(j + k)
				if ok1 && ok2 {
					g1 := GeneralizedNode{element: child_j, length: k}
					g2 := GeneralizedNode{element: child_jk, length: k}
					container := GeneralizedNodeCompareContainer{left: g1, right: g2}
					score_item, _ := scores[container.identity()]
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
			if m.calculateScore(cur_region) > m.calculateScore(max_region) {
				max_region.K = cur_region.K
				max_region.Start = cur_region.Start
				max_region.Covered = cur_region.Covered
				max_region.Score = cur_region.Score
			}
		}
	}
	if max_region.Covered > 0 {
		data_regions = append(data_regions, max_region)
		if max_region.Start+max_region.Covered < len(max_region.Parent.Children) {
			data_regions = append(data_regions, m.identifyRegions(max_region.Start+max_region.Covered, root, max_generalized_nodes, threshold, scores)...)
		}
	}
	return data_regions
}

func (m *MiningDataRegion) findRegions(root *DataTree) []*DataRegion {
	data_regions := []*DataRegion{}
	if root.treeDepth() >= 2 {
		scores := m.compareGeneralizedNodes(root, m.max_generalized_nodes)
		data_regions = append(data_regions, m.identifyRegions(0, root, m.max_generalized_nodes, m.threshold, scores)...)
		covered := map[string]*DataTree{}

		for _, data_region := range data_regions {
			// all items that are covered by this data_region
			for i := data_region.Start; i < data_region.Covered; i++ {
				ichild, ok := data_region.Parent.getChild(i)
				if ok {
					covered[ichild.hash()] = ichild
				}
			}
		}

		// for each child that is not covered by that data region
		for _, child := range root.Children {
			_, ok := covered[child.hash()]
			if ok == false {
				data_regions = append(data_regions, m.findRegions(child)...)
			}
		}
	}
	return data_regions
}
