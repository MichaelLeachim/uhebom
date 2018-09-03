// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-29-08 19:39<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package uhebom

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestExamlesNWParsing(t *testing.T) {
	testFiles := map[string]string{
		"test/1.html": "test/resultnw/1.html",
		"test/3.html": "test/resultnw/3.html",
		"test/4.html": "test/resultnw/4.html",
		"test/5.html": "test/resultnw/5.html",
		"test/6.html": "test/resultnw/6.html",
	}
	for k, v := range testFiles {
		os.RemoveAll(v)
		file, err := ioutil.ReadFile(k)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(v, wunsch_processing.asHTMLTables(wunsch_processing.extract(file)), 0777)
	}
}

func TestImplementPathRecords(t *testing.T) {
	data, _ := ioutil.ReadFile("test/1.html")
	regions := simplified_api.findDataRegions(data)
	assert.Equal(t, len(regions), 28)

	records := simplified_api.findDataRecords(regions[2])
	assert.Equal(t,
		wunsch_processing.dereferenceListOfForms(
			wunsch_processing.convertToTabularForm(records[0].convertToBase())),
		[]TabularForm{
			TabularForm{Tag: "", Path: "div/div/div/", Content: "", IsGap: false},
			TabularForm{Tag: "div", Path: "div/div/div/a/div", Content: "", IsGap: false},
			TabularForm{Tag: "i", Path: "div/div/div/div/h4/i", Content: "", IsGap: false},
			TabularForm{Tag: "text", Path: "div/div/div/div/h4/text", Content: "1500", IsGap: false},
			TabularForm{Tag: "br", Path: "div/div/div/div/br", Content: "", IsGap: false},
			TabularForm{Tag: "text", Path: "div/div/div/div/h4/a/text", Content: "Dictionnaire encyclopedique pour tous. Nouveau Petit Larousse en couleurs.", IsGap: false},
			TabularForm{Tag: "text", Path: "div/div/div/div/text", Content: "Энциклопедия на французском языке. Содержит 70500 статей, 5150 иллюстраций, 245 карт.", IsGap: false}})
}
