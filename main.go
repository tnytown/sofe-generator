package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/paint"
	"github.com/anthonynsimon/bild/transform"
)

type info struct {
	style    style
	rotation int
	color1   color.Color
	color2   color.Color
}

type style byte

const (
	inverted style = 0
	normal   style = 1
)

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/generate", generate)
	http.ListenAndServe(":"+os.Getenv("HTTP_PLATFORM_PORT"), nil)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Nothing to see here.")
}

func generate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	i, err := parseForm(r.Form)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err)
		return
	}

	t, err := imgio.Open("sofe_template.png")

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "An internal server error occurred.")
		return
	}

	// flood fill the colors: color1 = outlying area, color2 = inner shapes, TODO find a less stupid method
	t = paint.FloodFill(t, image.Point{0, 0}, i.color1, 0)
	t = paint.FloodFill(t, image.Point{100, 100}, i.color2, 0)
	t = paint.FloodFill(t, image.Point{300, 100}, i.color2, 0)
	t = paint.FloodFill(t, image.Point{100, 320}, i.color2, 0)

	// invert the image if specified
	if i.style == inverted {
		t = effect.Invert(t)
	}

	// rotate the image if specified
	if i.rotation != 0 && i.rotation > -360 && i.rotation < 360 {
		t = transform.Rotate(t, float64(i.rotation), nil)
	}

	w.Header().Add("Content-Type", "image/png")
	png.Encode(w, t)
}

func parseForm(f url.Values) (info, error) {
	i := info{}
	rng := rand.New(rand.NewSource(time.Now().Unix()))

	switch f.Get("style") {
	case "inverted":
		i.style = inverted
	case "normal":
		i.style = normal
	default:
		i.style = style(rand.Intn(2))
	}

	var err error
	i.rotation, err = strconv.Atoi(f.Get("rot"))
	if err != nil {
		i.rotation = 0
	}

	if f.Get("color1") != "" {
		i.color1, err = hexToRGBA(f.Get("color1"))
		if err != nil {
			return i, errors.New("color1 is invalid")
		}
	} else {
		i.color1 = color.RGBA{uint8(rng.Uint32()), uint8(rng.Uint32()), uint8(rng.Uint32()), 100}
	}

	if f.Get("color2") != "" {
		i.color2, err = hexToRGBA(f.Get("color2"))
		if err != nil {
			return i, errors.New("color2 is invalid")
		}
	} else {
		i.color2 = color.RGBA{uint8(rng.Uint32()), uint8(rng.Uint32()), uint8(rng.Uint32()), 100}
	}
	return i, nil
}

func hexToRGBA(h string) (color.Color, error) {
	if len(h) != 6 {
		return nil, errors.New("hex code is of invalid length")
	}

	r, _ := strconv.ParseUint(h[0:2], 16, 8)
	g, _ := strconv.ParseUint(h[2:4], 16, 8)
	b, _ := strconv.ParseUint(h[4:6], 16, 8)

	return color.RGBA{uint8(r), uint8(g), uint8(b), 100}, nil
}
