// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-29-08 19:39<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package depta

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestExamlesNWParsing(t *testing.T) {
	testFiles := map[string]string{
		"test/1.html": "test/resultnw/1.html",
		"test/2.html": "test/resultnw/2.html",
		"test/3.html": "test/resultnw/3.html",
		"test/4.html": "test/resultnw/4.ru",
	}
	for k, v := range testFiles {
		os.RemoveAll(v)
		file, err := ioutil.ReadFile(k)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(v, wunsch_processing.as_html_tables(wunsch_processing.Extract(file)), 0777)
	}
}

func TestImplementPathRecords(t *testing.T) {
	data, _ := ioutil.ReadFile("test/1.html")
	dtree := ReadHtml(ParseHTML(bytes.NewReader(data)))
	region_finder := MiningDataRegion{}
	region_finder.init(dtree, 5, 0.75)
	regions := region_finder.find_regions(dtree)
	assert.Equal(t, len(regions), 28)
	record_finder := MiningDataRecord{}
	record_finder.init(0.75)
	records := record_finder.find_records(regions[2])

	assert.Equal(t, wunsch_processing.ConvertToTabularForm(records[0].convert_to_base()),
		[]TabularForm{TabularForm{Path: "div/div/div/", Content: ""},
			TabularForm{Path: "div/div/div/a/div", Content: ""},
			TabularForm{Path: "div/div/div/div/h4/i", Content: ""},
			TabularForm{Path: "div/div/div/div/h4/text", Content: "1500"},
			TabularForm{Path: "div/div/div/div/br", Content: ""},
			TabularForm{Path: "div/div/div/div/h4/a/text", Content: "Dictionnaire encyclopedique pour tous. Nouveau Petit Larousse en couleurs."},
			TabularForm{Path: "div/div/div/div/text", Content: "Энциклопедия на французском языке. Содержит 70500 статей, 5150 иллюстраций, 245 карт."}})

}
