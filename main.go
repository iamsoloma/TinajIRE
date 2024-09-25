package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"os"

	p "github.com/iamsoloma/TinajIRE/primitives"
	"github.com/iamsoloma/TinajIRE/splat"
)

const (
	nx = 400 * 5 // size of x
	ny = 200 * 5 // size of y
	ns = 1       // number of samples for aa
	c  = 255.99
)

var (
	white = p.Vector{1.0, 1.0, 1.0}
	blue  = p.Vector{0.5, 0.7, 1.0}

	camera = p.NewCamera(p.Vector{2, 0, 13})

	//sphere = p.Sphere{p.Vector{0, 0, -1}, 0.5}
	//floor  = p.Sphere{p.Vector{0, -100.5, -1}, 100}

	//world = p.World{[]p.Hitable{&sphere, &floor}}
)

func main() {

	fmt.Println("Визуализатор облака точек")
	fmt.Println("Читаем облако точек...")

	cloud, err := splat.ReadCloud("./flowerPoints.ply")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can`t read a point cloud!\n"+err.Error())
		os.Exit(1)
	}
	//check(err, "Can`t read a point cloud!\n" + err.Error())
	world := p.World{}
	for _, p := range cloud {
		world.Elements = append(world.Elements, &p)
	}
	fmt.Println("Прочитано точек: "+strconv.Itoa(len(world.Elements)))

	f, err := os.Create("out.ppm")
	check(err, "Error opening file: %v\n")
	defer f.Close()

	fmt.Println("Рендерим...")

	// http://netpbm.sourceforge.net/doc/ppm.html
	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", nx, ny)
	check(err, "Error writting to file: %v\n")

	// writes each pixel with r/g/b values
	// from top left to bottom right
	allPxls := nx * ny
	counter := 0

	//imagstrs := [][]p.Vector{} - капец долго, хороший повод ознакомиться с goroutine.

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			counter += 1
			rgb := p.Vector{}
			
			// sample rays for anti-aliasing
			for s := 0; s < ns; s++ {
				u := (float64(i) + rand.Float64()) / float64(nx)
				v := (float64(j) + rand.Float64()) / float64(ny)

				r := camera.RayAt(u, v)
				color := color(&r, &world)
				rgb = rgb.Add(color)
			}

			// average
			rgb = rgb.DivideScalar(float64(ns))

			// get intensity of colors
			ir := int(/*c * */rgb.X)
			ig := int(/*c * */rgb.Y)
			ib := int(/*c * */rgb.Z)

			if (counter%100 == 0) {
				fmt.Println("Pxl: " + strconv.Itoa(counter) + " : " + strconv.Itoa(allPxls))
			}

			_, err = fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)
			check(err, "Error writing to file: %v\n")
		}
	}

	fmt.Println("Rendering is finished!")
}

func color(r *p.Ray, h p.Hitable) p.Vector {
	hit, record := h.Hit(r, 0.0, math.MaxFloat64)

	if hit {
		return record.Color//record.Normal.AddScalar(1.0).MultiplyScalar(0.5)
	}

	// make unit vector so y is between -1.0 and 1.0
	unitDirection := r.Direction.Normalize()

	return gradient(& unitDirection )
}

func gradient(v *p.Vector) p.Vector {
	// scale t to be between 0.0 and 1.0
	t := 0.5 * (v.Y + 1.0)

	// linear blend: blended_value = (1 - t) * white + t * blue
	return white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
}

func check(e error, s string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, s, e)
		os.Exit(1)
	}
}
