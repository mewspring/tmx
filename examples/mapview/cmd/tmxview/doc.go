/*

mapview creates image representations of tmx maps.

Installation:

	go get github.com/mewmew/tmx/examples/mapview/cmd/mapview

Documentation:

Documentation provided by GoDoc.

http://godoc.org/github.com/mewmew/tmx/examples/mapview/cmd/mapview

Usage:

	mapview [OPTION]... [FILE]...

Flags:

	-o (default="view.png")
		Output image path.

Examples:

1. Create a png image of a tmx map.
	mapview -o map.png map.tmx

*/
package main
