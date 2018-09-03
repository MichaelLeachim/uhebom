// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-01-09 23:14<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

// mining the data record from a region.
// basic assumption:
// the subtree of data records also similar. so if not any adjacent pair of them are
// similar (less than threshold), data region itself is a data record,
// otherwise children are individual data record.

package uhebom

type MiningDataRecord struct {
	threshold float64
	stm       *SimpleTreeMatch
}

func newMiningDataRecord(threshold float64) *MiningDataRecord {
	m := MiningDataRecord{}
	m.threshold = threshold
	m.stm = newSimpleTreeMatch()
	return &m
}

func (m *MiningDataRecord) sliceRegion(region *DataRegion) []*DataRecord {
	// slice every generalized node of region to a data records
	records := make([]*DataRecord, 0)
	for _, v := range region.iter(region.K) {
		vv := DataRecord(v)
		records = append(records, &vv)
	}
	return records
}

func (m *MiningDataRecord) almostSimilar(similarities map[string]float64, threshold float64) bool {
	sims := make([]int, 0)

	for _, v := range similarities {
		if float64(v) >= threshold {
			sims = append(sims, 1)
		}
	}
	return (float64(len(sims)) / float64(len(similarities))) > ALMOST_SIMILAR
}

func (m *MiningDataRecord) findRecords(region *DataRegion) []*DataRecord {
	// PROBABLY BUGS, TEST
	if region.K == 1 {
		records := make([]*DataRecord, 0)
		//if all the individual nodes of children nodes of Generalized node are similar
		for i := region.Start; i < region.Start+region.Covered; i++ {
			child, ok := region.Parent.getChild(i)
			if ok {
				for _, children := range trees_utils.pairwise(child.Children, 1, 0) {
					sim := m.stm.normalizedMatchScore(children[0], children[1])
					if float64(sim) < m.threshold {
						return m.sliceRegion(region)
					}
				}
			}
		}
		// each child of generalized node is a data record
		for _, gn := range region.iter(1) {
			for _, c := range gn {
				cc := DataRecord{c}
				records = append(records, &cc)
			}
		}
		return records
	} else {
		// if almost all the individual nodes in Generalized Node are similar
		children := make([]*DataTree, 0)
		for i := 0; i < region.Covered; i++ {
			child, ok := region.Parent.getChild(region.Start + i)
			if ok {
				children = append(children, child)
			}

		}
		sizes := make(map[int]int, 0)
		for _, child := range children {
			sizes[child.treeSize()] += 1
		}
		most_common_size_counter := 0
		most_common_size := 0
		for k, v := range sizes {
			if v > most_common_size_counter {
				most_common_size_counter = v
				most_common_size = k
			}
		}
		var most_typical_child *DataTree
		for _, v := range children {
			if v.treeSize() == most_common_size {
				most_typical_child = v
			}
		}
		similarities := make(map[string]float64, 0)
		for _, child := range children {
			similarities[child.hash()] = m.stm.normalizedMatchScore([]*DataTree{most_typical_child}, []*DataTree{child})
		}
		result := make([]*DataRecord, 0)
		if m.almostSimilar(similarities, m.threshold) {
			for _, child := range children {
				if float64(similarities[child.hash()]) >= m.threshold {
					rr := DataRecord([]*DataTree{child})
					result = append(result, &rr)
				}
			}
			return result
		} else {
			return m.sliceRegion(region)
		}
	}
}
