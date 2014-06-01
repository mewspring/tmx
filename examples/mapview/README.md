mapview
=======

The mapview package facilitates drawing image representations of [tmx][] maps.

[tmx]: https://github.com/bjorn/tiled/wiki/TMX-Map-Format

Documentation
-------------

Documentation provided by GoDoc.

- [mapview][]
	- [tile][]

[mapview]: http://godoc.org/github.com/mewmew/tmx/examples/mapview
[tile]: http://godoc.org/github.com/mewmew/tmx/examples/mapview/tile

tmxview
=======

tmxview creates image representations of [tmx][] maps.

Installation
------------

	go get github.com/mewmew/tmx/examples/mapview/cmd/tmxview

Documentation
-------------

Documentation provided by GoDoc.

- [tmxview][]: Create image representations of tmx maps.

[tmxview]: http://godoc.org/github.com/mewmew/tmx/examples/mapview/cmd/tmxview

Usage
-----

	tmxview [OPTION]... [FILE]...

Flags:

	-o (default="view.png")
		Output image path.

Examples
--------

1. Create a png image of a tmx map.

		tmxview -o map.png map.tmx

public domain
-------------

This code is hereby released into the *[public domain][]*.

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
