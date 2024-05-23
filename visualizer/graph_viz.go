package visualizer

/*
import (
	"image"
	_ "image/png"
	"os"
	"time"

	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// renders a text and loads it to the gpu, returns the texture_id
func Pre_Render_Image(w, h int, png image.Image) uint32 {
	var texture_id uint32
	gl.GenTextures(1, &texture_id)
	img, ok := png.(*image.NRGBA)
	if ok {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture_id)
		dst := img

		w_image := img.Rect.Dx()
		h_image := img.Rect.Dy()

		//off_x := int(w/2) - int(img.Bounds().Dx())/2
		//off_y := int(h/2) + int(img.Bounds().Dy())/2
		//r := img.Bounds().Add(image.Pt(off_x, off_y))

		//draw.Draw(dst, r, img, img.Bounds().Min, draw.Src)

		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w_image), int32(h_image), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(dst.Pix))
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	} else {
		panic("png cannot be converted to an image.RGBA!")
	}
	return texture_id
}

// @Todo not needed
func ImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func Make_Plot() (image.Image, error) {
	rand.Seed(int64(0))

	p := plot.New()

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	err := plotutil.AddLinePoints(p,
		"First", randomPoints(15),
		"Second", randomPoints(15),
		"Third", randomPoints(15))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}

	var vgc vg.Canvas
	canvas := draw.NewCanvas(vgc, 6*vg.Inch, 4.8*vg.Inch)
	p.Draw(canvas)

	rect1 := vg.Rectangle{Min: vg.Point{0, 0}, Max: vg.Point{6 * vg.Inch, 4.8 * vg.Inch}}
	rect2 := image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{int(6 * vg.Inch), int(5 * vg.Inch)}}

	img := image.NewNRGBA(rect2)
	canvas.DrawImage(rect1, img)

	return img, nil
}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}

func Visualize_Graph() {

	w, h := 640, 480

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(w, h, "Data View", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetKeyCallback(KeyCallback)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	// @Todo handle resize? Fix draw bug
	// see https://gamedev.stackexchange.com/questions/49674/opengl-resizing-display-and-glortho-glviewport
	gl.Ortho(0, float64(w), float64(h), 0, -1, 1)

	img, err := ImageFromFilePath("hexiflexi1.png")
	if err != nil {
		panic(err)
	}
	texture_id := Pre_Render_Image(w, h, img)

	Make_Plot()

	for !window.ShouldClose() {
		gl.ClearColor(0.1, 0.1, 0.1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		draw_texture(0, 0, int32(w), int32(h), texture_id)

		window.SwapBuffers()
		glfw.PollEvents()
		time.Sleep(time.Millisecond * 30)
	}
}
*/
