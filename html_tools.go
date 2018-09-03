// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @ Copyright (c) Michael Leachim                                                      @
// @ You can find additional information regarding licensing of this work in LICENSE.md @
// @ You must not remove this notice, or any other, from this software.                 @
// @ All rights reserved.                                                               @
// @@@@@@ At 2018-01-09 22:53<mklimoff222@gmail.com> @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
package uhebom

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"io"
	"strconv"
	"strings"
)

type html_tools_ struct{}

var html_tools = html_tools_{}

// wrap HTML into DTree
func (h *html_tools_) readHTML(n *html.Node) *DataTree {
	dtree := DataTree{}
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
		dtree.childAppend(h.readHTML(c))
	}
	return &dtree
}

func (h *html_tools_) parseHTML(data io.Reader) *html.Node {
	doc, err := html.Parse(data)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func (h *html_tools_) asHTMLTables(data []*DataRegion, show_id bool) []byte {
	var template bytes.Buffer
	template.WriteString("<style>table {border-collapse: collapse;}table, th, td {border: 1px solid black;}</style>")

	for i, item := range data {
		template.WriteString("<h1>Table number: ")
		template.WriteString(strconv.Itoa(i))
		template.WriteString("</h1>")
		template.WriteString(item.asHTMLTable(show_id))
	}
	return template.Bytes()
}
