# tmx

[![Build Status](https://travis-ci.org/mewmew/tmx.svg?branch=master)](https://travis-ci.org/mewmew/tmx)
[![Coverage Status](https://img.shields.io/coveralls/mewmew/tmx.svg)](https://coveralls.io/r/mewmew/tmx?branch=master)
[![GoDoc](https://godoc.org/github.com/mewmew/tmx?status.svg)](https://godoc.org/github.com/mewmew/tmx)

The tmx project provides support for reading [tmx][1] (Tile Map XML) files, used by the map editor [Tiled](https://github.com/bjorn/tiled/).

[1]: https://github.com/bjorn/tiled/wiki/TMX-Map-Format

## Documentation

Documentation provided by GoDoc.

- [tmx]: implements access to Tiled's tmx (Tile Map XML) files.

[tmx]: http://godoc.org/github.com/mewmew/tmx

## Examples

The `tmxview` command creates image representations of tmx maps.

    go get github.com/mewmew/tmx/examples/mapview/cmd/tmxview
    cd $GOPATH/src/github.com/mewmew/tmx/testdata
    tmxview test_csv.tmx

![Screenshot - tmxview](https://github.com/mewmew/tmx/blob/master/examples/mapview/cmd/tmxview/view.png?raw=true)

## Public domain

The source code and any original content of this repository is hereby released into the [public domain].

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/

## License

The tilesets in `testdata/` are part of the [Flare](http://flarerpg.org/) project and are licensed CC-BY-SA 3.0.

