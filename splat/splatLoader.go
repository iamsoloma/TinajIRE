package splat

import (
	"fmt"
	"io"
	"os"

	"github.com/black40x/plyfile/plyfile"
	p "github.com/iamsoloma/TinajIRE/primitives"
)

func ReadCloud(file string) (spheres []p.Sphere, err error) {
	spheres = []p.Sphere{}

	ply, err := plyfile.Open(file)
	if err != nil {
		return spheres, err
	}
	defer ply.Close()

	f, err := os.Create("cloudOut.dot")
	check(err, "Error opening file: %v\n")
	defer f.Close()

	r, err := ply.GetElementReader("vertex")
	if err == nil {
		point := plyfile.Point{}

		for {
			_, err := r.ReadNext(&point)
			if err == io.EOF {
				break
			}
			//fmt.Println(point.String())
			spheres = append(spheres,
				p.Sphere{
					Center: p.Vector{X: point.X, Y: point.Y, Z: point.Z},
					Radius: 0.01,
					Color: p.Vector{
						X: float64(point.R),
						Y: float64(point.G),
						Z: float64(point.B),
					}})
			_, err = fmt.Fprintf(f, "%s\n", point.String())
			check(err, "Error writing to file: %v\n")
		}

		return spheres, nil

	} else {
		return spheres, err
	}
}

func check(e error, s string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, s, e)
		os.Exit(1)
	}
}
