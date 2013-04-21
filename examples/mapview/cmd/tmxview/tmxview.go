package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mewkiz/pkg/imgutil"
	"github.com/mewmew/tmx"
	"github.com/mewmew/tmx/examples/mapview"
)

// pngPath is the path to the output png image.
var pngPath string

func init() {
	flag.StringVar(&pngPath, "o", "view.png", "Output image path.")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: mapview [OPTION]... [FILE]...")
	fmt.Fprintln(os.Stderr, "Create image representations of tmx maps.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  Create png image of tmx map.")
	fmt.Fprintln(os.Stderr, "    mapview -o map.png map.tmx")
}

func main() {
	flag.Parse()
	for _, tmxPath := range flag.Args() {
		err := view(tmxPath)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func view(tmxPath string) (err error) {
	m, err := tmx.Open(tmxPath)
	if err != nil {
		return err
	}
	view, err := mapview.NewView(m, path.Dir(tmxPath))
	if err != nil {
		return err
	}
	view.Draw()
	err = imgutil.WriteFile(pngPath, view)
	if err != nil {
		return err
	}
	return nil
}
