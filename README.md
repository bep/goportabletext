[![Tests on Linux, MacOS and Windows](https://github.com/bep/goportabletext/workflows/Test/badge.svg)](https://github.com/bep/goportabletext/actions?query=workflow:Test)
[![Go Report Card](https://goreportcard.com/badge/github.com/bep/goportabletext)](https://goreportcard.com/report/github.com/bep/goportabletext)
[![GoDoc](https://godoc.org/github.com/bep/goportabletext?status.svg)](https://godoc.org/github.com/bep/goportabletext)

WORK IN PROGRESS.

Converts [Portable Text](https://www.portabletext.org/) to Markdown.

Note that the image handling is currently very simple; we link to the `asset.url` using `asset.altText` as the image alt text and `asset.title` as the title.