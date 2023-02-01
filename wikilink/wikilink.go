/*
 *
 * Copyright 2023 puzzlemarkdownserver authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package wikilink

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

const openStr = "[["
const openLen = len(openStr)
const hash = '#'
const slash = '/'
const pipe = '|'
const closeStr = "]]"
const closeLen = len(closeStr)
const priority = 150

// Manage WikiLink targeting a custom WebComponent "wiki-link" :
//   - "[[ pageName ]]" became <wiki-link title="pageName">pageName</wiki-link>
//   - "[[ pageName | linkName ]]" became <wiki-link title="pageName">linkName</wiki-link>
//   - "[[ langTag/pageName ]]" became <wiki-link lang="langTag" title="pageName">pageName</wiki-link>
//   - "[[ path/to/wiki#pageName ]]" became <wiki-link wiki="path/to/wiki" title="pageName">pageName</wiki-link>
//   - "[[ path/to/wiki#langTag/pageName ]]" became <wiki-link wiki="path/to/wiki" lang="langTag" title="pageName">pageName</wiki-link>
//
// And so on...
var Extension goldmark.Extender = wikiLinkExtender{}
var Kind = ast.NewNodeKind("WikiLink")

// check matching with interface
var _ parser.InlineParser = wikiLinkParser{}
var _ renderer.NodeRenderer = wikiLinkRenderer{}

var start = []byte{'['}
var open = []byte(openStr)
var close = []byte(closeStr)

type wikiLinkNode struct {
	ast.BaseInline
	WikiPath []byte
	Lang     []byte
	Title    []byte
}

func (*wikiLinkNode) Kind() ast.NodeKind {
	return Kind
}

func (n *wikiLinkNode) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"Lang":  string(n.Lang),
		"Title": string(n.Title),
	}, nil)
}

type wikiLinkParser struct{}

func (wikiLinkParser) Trigger() []byte {
	return start
}

func (wikiLinkParser) Parse(parent ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, seg := block.PeekLine()
	stop := bytes.Index(line, close)
	if stop < 0 || !bytes.HasPrefix(line, open) {
		return nil
	}

	seg = text.NewSegment(seg.Start+openLen, seg.Start+stop)
	var wikiPath []byte
	title := block.Value(seg)
	if index := bytes.IndexByte(title, hash); index >= 0 {
		wikiPath = title[:index]
		index++
		title = title[index:]
		seg = seg.WithStart(seg.Start + index)
	}

	var lang []byte
	if index := bytes.IndexByte(title, slash); index >= 0 {
		lang = title[:index]
		index++
		title = title[index:]
		seg = seg.WithStart(seg.Start + index)
	}

	if index := bytes.IndexByte(title, pipe); index >= 0 {
		title = title[:index]
		seg = seg.WithStart(seg.Start + index + 1)
	}

	node := &wikiLinkNode{WikiPath: trim(wikiPath), Lang: trim(lang), Title: trim(title)}
	node.AppendChild(node, ast.NewTextSegment(seg))
	block.Advance(stop + closeLen)
	return node
}

func trim(data []byte) []byte {
	start := 0
	for index, b := range data {
		if !(b == ' ' || b == '\t') {
			start = index - 1
			break
		}
	}
	index := len(data) - 1
	for {
		b := data[index]
		if !(b == ' ' || b == '\t') {
			break
		}
		index--
	}
	return data[start : index+1]
}

type wikiLinkRenderer struct{}

func (wikiLinkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, renderWikiLink)
}

func renderWikiLink(writer util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	wikiLink, ok := node.(*wikiLinkNode)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node %T, expected *wikiLinkNode", node)
	}

	if entering {
		writer.WriteString("<wiki-link")
		wikiPath := wikiLink.WikiPath
		if len(wikiPath) != 0 {
			writer.WriteString(" wiki=\"")
			writer.Write(wikiPath)
			writer.WriteByte('"')
		}
		lang := wikiLink.Lang
		if len(lang) != 0 {
			writer.WriteString(" lang=\"")
			writer.Write(lang)
			writer.WriteByte('"')
		}
		writer.WriteString(" title=\"")
		writer.Write(wikiLink.Title)
		writer.WriteString("\">")
	} else {
		writer.WriteString("</wiki-link>")
	}

	return ast.WalkContinue, nil
}

type wikiLinkExtender struct{}

func (wikiLinkExtender) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(wikiLinkParser{}, priority),
		),
	)

	md.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(wikiLinkRenderer{}, priority),
		),
	)
}
