// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-02-09 00:13<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
package depta

import (
	"bytes"
)

const (
	MAX_GENERALIZED_NODES = 5
	THRESHOLD             = 0.75
)

type simplified_api_ struct{}

var simplified_api = simplified_api_{}

func (s *simplified_api_) ReadAndParseBytes(data []byte) *DataTree {
	return html_tools.ReadHtml(html_tools.ParseHTML(bytes.NewReader(data)))
}

func (s *simplified_api_) FindDataRecords(region *DataRegion) []*DataRecord {
	return newMiningDataRecord(THRESHOLD).find_records(region)
}

func (s *simplified_api_) FindDataRegions(data []byte) []*DataRegion {
	root := s.ReadAndParseBytes(data)
	return newMiningDataRegion(root, MAX_GENERALIZED_NODES, THRESHOLD).find_regions(root)
}
