package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
)

type icon struct {
	XMLName         xml.Name `xml:"svg"`
	XMLNS           string   `xml:"xmlns,attr"`
	Version         string   `xml:"version,attr"`
	Width           int      `xml:"width,attr"`
	Height          int      `xml:"height,attr"`
	Defs            defs     `xml:"defs"`
	BackgroundImage *bgImage `xml:"image,omitempty"`
	BackgroundRect  *bgRect  `xml:"rect,omitempty"`
	PChars          []pChar  `xml:"text"`
}

type defs struct {
	RawXML string `xml:",innerxml"`
}

type bgImage struct {
	XMLName xml.Name `xml:"image"`
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
	HRef    string   `xml:"http://www.w3.org/1999/xlink href,attr"`
}

type bgRect struct {
	XMLName xml.Name `xml:"rect"`
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
	Fill    string   `xml:"fill,attr"`
}

type pChar struct {
	XMLName    xml.Name `xml:"text"`
	Display    string   `xml:"display,attr,omitempty"`
	X          float64  `xml:"x,attr"`
	Y          float64  `xml:"y,attr"`
	Fill       string   `xml:"fill,attr"`
	FontSize   float64  `xml:"font-size,attr"`
	FontFamily string   `xml:"font-family,attr"`
	P          string   `xml:",chardata"`
}

const smallIconBreak = 60
const verySmallIconBreak = 48
const squeezeInSize = 100

type imageExtension string

const (
	png  imageExtension = "png"
	jpeg imageExtension = "jpg"
)

var sizes = map[int]imageExtension{
	600: jpeg,
	512: png,
	310: png,
	192: png,
	180: png,
	152: png,
	150: png,
	144: png,
	120: png,
	114: png,
	96:  png,
	76:  png,
	72:  png,
	70:  png,
	60:  png,
	57:  png,
	48:  png,
	36:  png,
	32:  png,
	16:  png,
}

func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Could not get executable path: %s", err)
		return
	}
	exPath := filepath.Dir(ex)
	renderPath := filepath.Join(exPath, "render")
	iconFile, err := os.Open(filepath.Join(exPath, "icon-base.svg"))
	if err != nil {
		log.Fatalf("Could not open icon SVG: %s", err)
	}
	defer iconFile.Close()
	iconBytes, err := ioutil.ReadAll(iconFile)
	if err != nil {
		log.Fatalf("Could not read icon SVG: %s", err)
	}
	smallIconFile, err := os.Open(filepath.Join(exPath, "icon-base-small.svg"))
	if err != nil {
		log.Fatalf("Could not open small icon SVG: %s", err)
	}
	defer smallIconFile.Close()
	smallIconBytes, err := ioutil.ReadAll(smallIconFile)
	if err != nil {
		log.Fatalf("Could not read small icon SVG: %s", err)
	}
	verySmallIconFile, err := os.Open(filepath.Join(exPath, "icon-base-very-small.svg"))
	if err != nil {
		log.Fatalf("Could not open small icon SVG: %s", err)
	}
	defer verySmallIconFile.Close()
	verySmallIconBytes, err := ioutil.ReadAll(verySmallIconFile)
	if err != nil {
		log.Fatalf("Could not read small icon SVG: %s", err)
	}
	altSvg := filepath.Join(exPath, "icon-small-special.svg")
	for size, imageExtension := range sizes {
		log.Printf("Handling %[1]dx%[1]d", size)
		var icon icon
		if size >= smallIconBreak {
			err = xml.Unmarshal(iconBytes, &icon)
			if err != nil {
				log.Fatalf("Could not parse icon SVG: %s", err)
			}
		} else if size >= verySmallIconBreak {
			err = xml.Unmarshal(smallIconBytes, &icon)
			if err != nil {
				log.Fatalf("Could not parse small icon SVG: %s", err)
			}
		} else {
			err = xml.Unmarshal(verySmallIconBytes, &icon)
			if err != nil {
				log.Fatalf("Could not parse very small icon SVG: %s", err)
			}
		}
		baseSize := icon.Width
		ratio := float64(size) / float64(baseSize)
		icon.Width = size
		icon.Height = size
		if icon.BackgroundImage != nil {
			icon.BackgroundImage.Width = size
			icon.BackgroundImage.Height = size
		}
		if icon.BackgroundRect != nil {
			icon.BackgroundRect.Width = size
			icon.BackgroundRect.Height = size
		}
		for i := range icon.PChars {
			icon.PChars[i].X *= ratio
			icon.PChars[i].Y *= ratio
			icon.PChars[i].FontSize *= ratio
			if size < squeezeInSize {
				switch i {
				case 0:
					icon.PChars[i].X = math.Ceil(icon.PChars[i].X) + 2
				case 2:
					icon.PChars[i].X = math.Floor(icon.PChars[i].X) - 2
				}
			}
		}
		newIcon, err := xml.MarshalIndent(icon, "", "\t")
		if err != nil {
			log.Fatalf("Could not marshal new icon SVG: %s", err)
		}
		svgOut := filepath.Join(renderPath, fmt.Sprintf("%d.svg", size))
		out := filepath.Join(renderPath, fmt.Sprintf("%d.%s", size, imageExtension))
		outFile, err := os.Create(svgOut)
		if err != nil {
			log.Fatalf("Could not create new icon SVG file: %s", err)
		}
		defer outFile.Close()
		_, err = io.WriteString(outFile, string(newIcon))
		if err != nil {
			log.Fatalf("Could not write to new icon SVG file: %s", err)
		}
		err = outFile.Sync()
		if err != nil {
			log.Fatalf("Could not finish writing to new icon SVG file: %s", err)
		}
		sizeStr := fmt.Sprintf("%d", size)
		log.Println("inkscape", "-z", "-e", out, "-w", sizeStr, "-h", sizeStr, svgOut)
		cmd := exec.Command("inkscape", "-z", "-e", out, "-w", sizeStr, "-h", sizeStr, svgOut)
		if err := cmd.Run(); err != nil {
			log.Fatalf("Could not run Inkscape: %s", err)
		}
		if size < smallIconBreak {
			out := filepath.Join(renderPath, fmt.Sprintf("%d-alt.%s", size, imageExtension))
			log.Println("inkscape", "-z", "-e", out, "-w", sizeStr, "-h", sizeStr, altSvg)
			cmd := exec.Command("inkscape", "-z", "-e", out, "-w", sizeStr, "-h", sizeStr, altSvg)
			if err := cmd.Run(); err != nil {
				log.Fatalf("Could not run Inkscape: %s", err)
			}
		}
	}
	log.Println("Done!")
}
