// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-27-08 22:33<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

package depta

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func ParseHTML(data io.Reader) *html.Node {
	doc, err := html.Parse(data)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func ReadHtml(n *html.Node) *DTree {
	// wrap html nodes into our DTree format
	// we might have done work on generic html parser, but, for the sake of independence
	// of other models, we shall use our data.
	dtree := DTree{}
	dtree.Attrs = make(map[string]string, 0)
	dtree.identity()
	if n.Type == html.ElementNode {
		dtree.Tag = n.Data
		for _, a := range n.Attr {
			dtree.Attrs[a.Key] = a.Val
		}
		if n.Data == "img" {
			dtree.Data = dtree.Attrs["src"]
		}
		if n.Data == "a" {
			dtree.Data = dtree.Attrs["href"]
		}
		// delete after debug
		// dtree.Attrs = make(map[string]string, 0)
	}
	if n.Type == html.TextNode {
		dtree.Tag = "text"
		dtree.Data = n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// skip empty text nodes
		if c.Type == html.TextNode {
			if strings.TrimSpace(c.Data) == "" {
				continue
			}
		}
		dtree.child_append(ReadHtml(c))
	}
	return &dtree
}

func DeptaExtract(root *DTree, k int, threshold float64) []*DataRegion {
	region_finder := MiningDataRegion{}
	region_finder.init(root, k, threshold)
	regions := region_finder.find_regions(root)

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
	return result

}

func Extract(data []byte) [][][]string {
	threshold := 0.75 // 0.75
	k := 5
	r := bytes.NewReader(data)
	dtree := ReadHtml(ParseHTML(r))
	result := DeptaExtract(dtree, k, threshold)

	return_data := make([][][]string, 0)
	for _, v := range result {
		return_data = append(return_data, v.AsPlainTexts())
	}
	return return_data
}
func ExtractToHTML(data []byte) []byte {
	threshold := 0.75
	k := 5
	r := bytes.NewReader(data)
	dtree := ReadHtml(ParseHTML(r))
	// log.WithField("dtree", dtree).Info("Initially read dtree")

	result := DeptaExtract(dtree, k, threshold)
	return as_html_tables(result, false)
}
