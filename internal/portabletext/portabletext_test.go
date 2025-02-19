// Copyright 2025 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package portabletext

import (
	"strings"
	"testing"

	"github.com/go-quicktest/qt"
)

const singleSpan = `
{
"_key": "R5FvMrjo",
"_type": "block",
"children": [
  { "_key": "cZUQGmh4", "_type": "span", "marks": [], "text": "Plain text." }
],
"markDefs": [],
"style": "normal"
}`

func TestParseBlocks(t *testing.T) {
	blocks, err := Parse(strings.NewReader(singleSpan))
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(blocks, Blocks{
		{
			BaseBlock: BaseBlock{Key: "R5FvMrjo", Type: "block"},
			Text: Text{
				Children: []Child{{Type: "span", Marks: []string{}, Text: "Plain text."}},
				MarkDefs: []MarkDef{},
				Style:    "normal",
			},
		},
	}))
}
