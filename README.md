tmx
===

This package provides support for reading [tmx][1] (Tile Map XML) files, used by
the map editor [Tiled][].

[1]: https://github.com/bjorn/tiled/wiki/TMX-Map-Format
[Tiled]: https://github.com/bjorn/tiled/

Documentation
-------------

Documentation provided by GoDoc.

   - [tmx][]: implements access to Tiled's tmx (Tile Map XML) files.

[tmx]: http://godoc.org/github.com/mewmew/tmx

Examples
--------

tmxview creates image representations of [tmx][] maps.

	go get github.com/mewmew/tmx/examples/mapview/cmd/tmxview
	cd $GOPATH/src/github.com/mewmew/tmx/testdata
	tmxview test_csv.tmx

![Screenshot - tmxview](https://github.com/mewmew/tmx/blob/master/examples/mapview/cmd/tmxview/view.png?raw=true)

public domain
-------------

This code is hereby released into the *[public domain][]*.

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/

license
-------

The tilesets in `testdata/` are part of the [Flare][] project and are licensed
CC-BY-SA 3.0.

[Flare]: http://flarerpg.org/
