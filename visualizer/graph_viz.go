package visualizer

import (
	"image"
	_ "image/png"
	"os"

	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// renders a text and loads it to the gpu, returns the texture_id
func Pre_Render_Image(w, h int, png image.Image) uint32 {
	var texture_id uint32
	gl.GenTextures(1, &texture_id)
	img, ok := png.(*image.RGBA)
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

func Make_Plot(w, h int, user_data *Data_Info) (image.Image, error) {
	p := plot.New()

	if user_data == nil {
		p.Title.Text = "We have no user_data. Please provide it!"
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
	} else {
		p.Title.Text = user_data.Title
		p.X.Label.Text = user_data.XLabel
		p.Y.Label.Text = user_data.YLabel

		var err error
		// I know.. this is stupid
		switch len(user_data.XY) {
		case 1:
			err = plotutil.AddLines(p,
				user_data.XY[0].Label, user_data.XY[0].XYs,
				user_data.XY[1].Label, user_data.XY[1].XYs,
			)
		case 2:
			err = plotutil.AddLines(p,
				user_data.XY[0].Label, user_data.XY[0].XYs,
				user_data.XY[1].Label, user_data.XY[1].XYs,
			)
		case 3:
			err = plotutil.AddLines(p,
				user_data.XY[0].Label, user_data.XY[0].XYs,
				user_data.XY[1].Label, user_data.XY[1].XYs,
				user_data.XY[2].Label, user_data.XY[2].XYs,
			)
		default:
			panic("Not implemented for len(user_data.XY) > 3: in Make_Plot()")
		}

		if err != nil {
			panic(err)
		}

		/*
			for i := range user_data.XY {
				err := plotutil.AddLines(p,
					user_data.XY[i].Label, user_data.XY[i].XYs,
				)

				if err != nil {
					panic(err)
				}
			} */
	}

	// Save the plot to a PNG file.
	//if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
	//	panic(err)
	//}

	dpi := 72 * 2
	dpi_o := vgimg.UseDPI(dpi)
	wh_o := vgimg.UseWH(vg.Length(w*72/dpi), vg.Length(h*72/dpi))
	canvas := vgimg.NewWith(dpi_o, wh_o)

	dc := draw.New(canvas)
	p.Draw(dc)
	img := canvas.Image()
	//fmt.Printf(img.Bounds().String())
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

func plotPoints(x []float64, y []float64) plotter.XYs {

	n := len(x)
	if n != len(y) {
		panic("x[] and y[] must have the same lenghts!")
	}

	pts := make(plotter.XYs, n)
	for i := range pts {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}
	return pts
}

type XY = plotter.XY

type XYData struct {
	XYs   plotter.XYs
	Label string
}

type Data_Info struct {
	XY     []XYData
	XLabel string
	YLabel string
	Title  string
}

type Graph_Viz struct {
	UserData *Data_Info // here the data should be set/updated by the user before Update_And_Draw()

	Window    *glfw.Window
	W, H      int
	Destroyed bool
	Initted   bool

	// internal data
	texture_id uint32
}

func (gv *Graph_Viz) Update_And_Draw() {

	if gv.Destroyed {
		// for now it's unrecoverable
		return
	}

	if !gv.Initted {
		// initialize and spawn window
		user_data := gv.UserData
		*gv = *NewGraphViz() //(will set inited to true)
		gv.UserData = user_data
	}

	if gv.Window.ShouldClose() {
		gv.Window.Destroy()
		gv.Destroyed = true
		return
	}

	img, err := Make_Plot(gv.W, gv.H, gv.UserData)
	if err != nil {
		panic(err)
	}

	gv.Window.MakeContextCurrent()
	gv.texture_id = Pre_Render_Image(gv.W, gv.H, img)

	gl.ClearColor(0.1, 0.1, 0.1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	draw_texture(0, 0, int32(gv.W), int32(gv.H), gv.texture_id)

	gv.Window.SwapBuffers()
	glfw.PollEvents()
}

func NewGraphViz() *Graph_Viz {

	w, h := 854, 480

	err := glfw.Init()
	if err != nil {
		panic(err)
	}

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

	return &Graph_Viz{
		Window:    window,
		W:         w,
		H:         h,
		Destroyed: false,
		Initted:   true,
	}
}
