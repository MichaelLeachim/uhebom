// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:34<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package depta

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestNodeCopyingInvestigation(t *testing.T) {
	data, _ := ioutil.ReadFile("test/1.html")
	dtree := ReadHtml(ParseHTML(bytes.NewReader(data)))
	regions := DeptaExtract(dtree, 5, 0.75)

	first_record := regions[2].Items[0].ConvertToBase()
	second_record := regions[2].Items[1].ConvertToBase()
	assert.Equal(t, len(first_record), 13)
	assert.Equal(t, len(second_record), 13)

	assert.NotEqual(t, first_record[12].Data, second_record[12].Data)

}

func TestDeeperNodeCopyingInvestigation(t *testing.T) {
	data, _ := ioutil.ReadFile("test/1.html")
	dtree := ReadHtml(ParseHTML(bytes.NewReader(data)))
	region_finder := MiningDataRegion{}
	region_finder.init(dtree, 5, 0.75)
	regions := region_finder.find_regions(dtree)
	assert.Equal(t, len(regions), 28)
	assert.Equal(t, regions[0], "sd")
	// assert.Equal(t, len(regions[0].Items[0].ConvertToBase()), "blab")
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
