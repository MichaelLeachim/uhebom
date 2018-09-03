// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-02-09 00:13<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
package uhebom

import (
	"bytes"
)

const (
	MAX_GENERALIZED_NODES = 5
	THRESHOLD             = 0.75
)

type simplified_api_ struct{}

var simplified_api = simplified_api_{}

func (s *simplified_api_) readAndParseBytes(data []byte) *DataTree {
	return html_tools.readHTML(html_tools.parseHTML(bytes.NewReader(data)))
}

func (s *simplified_api_) findDataRecords(region *DataRegion) []*DataRecord {
	return newMiningDataRecord(THRESHOLD).findRecords(region)
}

func (s *simplified_api_) findDataRegions(data []byte) []*DataRegion {
	root := s.readAndParseBytes(data)
	return newMiningDataRegion(root, MAX_GENERALIZED_NODES, THRESHOLD).findRegions(root)
}

func Extract(data []byte) [][][]string {
	result := [][][]string{}
	for _, table := range wunsch_processing.extract(data) {
		result_table := [][]string{}
		for _, row := range table {
			result_row := []string{}
			for _, col := range row {
				if col.IsGap {
					result_row = append(result_row, "")
				} else {
					result_row = append(result_row, col.Content)
				}
			}
			result_table = append(result_table, result_row)
		}
		result = append(result, result_table)
	}
	return result
}
