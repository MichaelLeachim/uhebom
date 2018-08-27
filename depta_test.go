// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:34<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package depta

import (
	// "github.com/stretchr/testify/assert"
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	// "strings"
	"testing"
)

func TestNodeCopying(t *testing.T) {
	// data, _ := ioutil.ReadFile("test/1.html")
	// items := Extract(data)
	// table4 := items[4]
	// solists1 := strings.TrimSpace(table4[0][22])
	// price1 := strings.TrimSpace(table4[0][16])
	// solists2 := strings.TrimSpace(table4[1][22])
	// price2 := strings.TrimSpace(table4[1][16])
	// assert.Equal(t, solists1, "Солисты московской государственной филармонии")
	// assert.Equal(t, price1, "1500")

	// // Investigation...
	// // make_id is allright
	// for i := 0; i <= 1000; i++ {
	// 	assert.NotEqual(t, string(make_id()), string(make_id()))
	// }
	// assert.NotEqual(t, solists1, solists2)
	// assert.NotEqual(t, price1, price2)

	// Parsing is OK. We must find something that copies data

}

func TestNodeCopyingInvestigation(t *testing.T) {
	data, _ := ioutil.ReadFile("test/1.html")
	threshold := 0.75 // 0.75
	k := 5
	r := bytes.NewReader(data)
	dtree := ReadHtml(ParseHTML(r))
	region_finder := MiningDataRegion{}
	region_finder.init(dtree, k, threshold)
	regions := region_finder.find_regions(dtree)

	record_finder := MiningDataRecord{}
	record_finder.init(threshold)

	field_finder := MiningDataField{}
	field_finder.init()

	result := make([]*DataRegion, 0)
	for _, region := range regions {
		records := record_finder.find_records(region)
		items, _ := field_finder.align_records(records)
		region.Items = items
		result = append(result, region)
	}
}

func TestExamlesParsing(t *testing.T) {
	testFiles := map[string]string{
		"test/1.html": "test/result/1.html",
		"test/2.html": "test/result/2.html",
		"test/3.html": "test/result/3.html",
		"test/4.html": "test/result/4.ru",
	}
	for k, v := range testFiles {
		os.RemoveAll(v)
		file, err := ioutil.ReadFile(k)
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(v, ExtractToHTML(file), 0777)
	}

}
