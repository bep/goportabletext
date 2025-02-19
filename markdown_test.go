// Copyright 2025 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package goportabletext_test

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/bep/goportabletext"
	"github.com/bep/goportabletext/internal/portabletext"
	"github.com/bep/goportabletext/internal/ptesting"
	"github.com/go-quicktest/qt"
)

func TestToMarkdownSamples(t *testing.T) {
	// Used during develoopment.
	pinnedName := ""
	testHelper(t, newSamplesHelper(t), pinnedName)
}

func BenchmarkToMarkdown(b *testing.B) {
	bh := func(h ptesting.BlockContentToMarkdownHelper) {
		h.ForEachGoldenPair(func(sourceName, targetName, sourceContent, targetContent string) error {
			blocks, err := portabletext.Parse(strings.NewReader(sourceContent))
			if err != nil {
				b.Fatal(err)
			}
			opts := goportabletext.ToMarkdownOptions{
				Dst: io.Discard,
				Src: blocks,
			}
			b.Run(sourceName, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					err := goportabletext.ToMarkdown(opts)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
			return nil
		})
	}

	bh(newSamplesHelper(b))
}

func newSamplesHelper(t testing.TB) ptesting.BlockContentToMarkdownHelper {
	b, err := ptesting.NewSamplesBlockContentToMarkdownHelper()
	qt.Assert(t, qt.IsNil(err))
	return b
}

func dosToUnix(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func jsonToMap(s string) any {
	var m any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		panic(err)
	}
	return m
}

func testHelper(t *testing.T, h ptesting.BlockContentToMarkdownHelper, pinnedName string) {
	t.Helper()
	if ptesting.IsCI() {
		// Make sure we run all tests in CI.
		pinnedName = ""
	}

	err := h.ForEachGoldenPair(func(sourceName, targetName, sourceContent, targetContent string) error {
		if pinnedName != "" && sourceName != pinnedName {
			return nil
		}
		checkSrc := func(src any) {
			var buff strings.Builder
			opts := goportabletext.ToMarkdownOptions{
				Dst: &buff,
				Src: src,
			}
			qt.Assert(t, qt.IsNil(goportabletext.ToMarkdown(opts)))
			result := dosToUnix(buff.String())
			targetContent = dosToUnix(targetContent)
			qt.Assert(t, qt.Equals(result, targetContent), qt.Commentf("input:\n%s\nexpected:\n%s\ngot:\n%s", sourceName, targetName, result))
		}
		checkSrc(strings.NewReader(sourceContent))
		checkSrc(jsonToMap(sourceContent))
		return nil
	})

	qt.Assert(t, qt.IsNil(err))
}
