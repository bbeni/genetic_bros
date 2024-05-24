package visualizer

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"math"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"

	"github.com/bbeni/genetic_bros/game"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

var fontFace48 font.Face
var fontFace64 font.Face
var fontFace128 font.Face

func init_fonts() {
	font, err := truetype.Parse(gomono.TTF)
	if err != nil {
		panic(err)
	}

	fontFace48 = truetype.NewFace(font, &truetype.Options{Size: 48})
	fontFace64 = truetype.NewFace(font, &truetype.Options{Size: 64})
	fontFace128 = truetype.NewFace(font, &truetype.Options{Size: 128})

}

func render_text(text string, text_color color.RGBA, fontFace font.Face) image.Image {
	drawer := &font.Drawer{
		Src:  &image.Uniform{text_color},
		Face: fontFace,
		Dot:  fixed.P(0, 0),
	}

	b26_6, _ := drawer.BoundString(text)
	bounds := image.Rect(
		b26_6.Min.X.Floor(),
		b26_6.Min.Y.Floor(),
		b26_6.Max.X.Ceil(),
		b26_6.Max.Y.Ceil(),
	)

	drawer.Dst = image.NewRGBA(bounds)
	drawer.DrawString(text)
	return drawer.Dst
}

func draw_rect(x, y, w, h int32, color [3]float32) {
	// Triangle
	gl.Begin(gl.TRIANGLES)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3i(x, y, 0)
	gl.Vertex3i(x, y+h, 0)
	gl.Vertex3i(x+w, y, 0)
	gl.Vertex3i(x+w, y+h, 0)
	gl.Vertex3i(x, y+h, 0)
	gl.Vertex3i(x+w, y, 0)
	gl.End()
}

const NUMBER_OF_NUMBERS = 16

var textures_intialzed bool
var texture_ids [NUMBER_OF_NUMBERS]uint32

func pre_render_numbers(w, h int) {

	for i := range NUMBER_OF_NUMBERS {
		n := (2 << i)
		text := fmt.Sprint(n)

		if n > 999 {
			texture_ids[i] = Pre_Render_Text_Centered(w, h, text, color.RGBA{24, 24, 24, 255}, fontFace48)
		} else {
			texture_ids[i] = Pre_Render_Text_Centered(w, h, text, color.RGBA{24, 24, 24, 255}, fontFace64)
		}

	}
}

// renders a text and loads it to the gpu, returns the texture_id
func Pre_Render_Text_Centered(w, h int, text string, c color.RGBA, fontFace font.Face) uint32 {
	var texture_id uint32
	gl.GenTextures(1, &texture_id)
	img, ok := render_text(text, c, fontFace).(*image.RGBA)
	if ok {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture_id)
		dst := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

		off_x := int(w/2) - int(img.Bounds().Dx())/2
		off_y := int(h/2) + int(img.Bounds().Dy())/2
		r := img.Bounds().Add(image.Pt(off_x, off_y))

		draw.Draw(dst, r, img, img.Bounds().Min, draw.Src)

		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(dst.Pix))
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	} else {
		panic("expected a image.RGBA!")
	}
	return texture_id
}

func draw_power_of_2(x, y, w, h int32, power int) {

	if !textures_intialzed {
		textures_intialzed = true
		pre_render_numbers(int(w), int(h))
	}

	texture_idx := power - 2
	texture_id := texture_ids[texture_idx]
	draw_texture(x, y, w, h, texture_id)
}

func draw_texture(x, y, w, h int32, texture_id uint32) {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture_id)
	gl.Enable(gl.TEXTURE_2D)

	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2i(x, y+h)
	gl.TexCoord2f(0, 0)
	gl.Vertex2i(x, y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2i(x+w, y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2i(x+w, y+h)
	gl.End()

	gl.Disable(gl.TEXTURE_2D)
}

func Draw_Tile(x, y, w, h int32, n int) {

	target := n
	power_of_2 := 0
	for target >= 1 {
		power_of_2++
		target = target >> 1
	}

	r := 0.25 + 0.04*float32(power_of_2)
	g := 0.4 + 0.024*float32(power_of_2)
	b := 0.2 + 0.06*float32(power_of_2)
	colr := [3]float32{r, g, b}

	draw_rect(x, y, w, h, colr)
	if power_of_2 == 0 {
		return
	}
	draw_power_of_2(x, y, w, h, power_of_2)
}

type Input_State struct {
	Pressed    bool // gets reset to false when handeled
	Dir        game.Direction
	Quit       bool
	Interacted bool
}

var input_state = Input_State{}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyEscape:
			input_state.Quit = true
		case glfw.KeyLeft:
			input_state.Pressed = true
			input_state.Dir = game.West
		case glfw.KeyRight:
			input_state.Pressed = true
			input_state.Dir = game.East
		case glfw.KeyUp:
			input_state.Pressed = true
			input_state.Dir = game.North
		case glfw.KeyDown:
			input_state.Pressed = true
			input_state.Dir = game.South
		}

	}
}

type Animation_State struct {
	State1 game.GameState
	State2 game.GameState
	T      float64 // form 0 to 1
	Dt     float64
}

func (ani_state *Animation_State) did_combine_lately() [4][4]bool {
	var result [4][4]bool
	for j := range 4 {
		for i := range 4 {
			if ani_state.State1.Board[j][i] != 0 && ani_state.State2.Board[j][i] != 0 && ani_state.State1.Board[j][i] != ani_state.State2.Board[j][i] {
				result[j][i] = true
			}
		}
	}
	return result
}

func (ani_state *Animation_State) update() bool {
	if ani_state.T == 1 {
		return false
	}
	ani_state.T += ani_state.Dt
	if ani_state.T > 1 {
		ani_state.T = 1
	}
	return true
}

// for Visulize_Game()
const (
	W            = 820
	H            = 820
	MARGIN       = 120
	TILE_SIZE    = 150
	TILE_PADDING = 8
)

var texture_game_over uint32

func Play_Game(gs *game.GameState) {
	Visualize_Game(gs, []game.Direction{}, -1, -1)
}

// if passed an empty array of moves we are just playing the game normaly
func Visualize_Game(gs *game.GameState, driver_moves []game.Direction, move_time float32, delay_on_gameover float32) {

	if len(driver_moves) == 0 {
		input_state.Interacted = true
	}

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(W, H, "2048", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetKeyCallback(KeyCallback)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Ortho(0, W, H, 0, -1, 1)

	init_fonts()

	g := *gs

	ani_state := Animation_State{}
	ani_state.State1 = g
	ani_state.State2 = g
	ani_state.T = 1
	ani_state.Dt = 0.1

	bloat_size := 20

	game_over := false
	game_over_timer := 3.0

	driver_index := 0

	game_time := 0.0

	for !window.ShouldClose() {

		// Handle input
		if input_state.Quit {
			break
		}

		if !game_over {
			if input_state.Pressed {
				input_state.Pressed = false
				input_state.Interacted = true
				ani_state.State1 = g
				if g.Update(input_state.Dir) {
					// it's game over
					if texture_game_over == 0 {
						texture_game_over = Pre_Render_Text_Centered(W, H, "Game Over!", color.RGBA{255, 255, 255, 255}, fontFace128)
					}
					game_over = true
					g = game.MakeGame()
				}
				ani_state.State2 = g
				ani_state.T = 0
			}
		} else {
			game_over_timer -= 0.02
			if game_over_timer <= 0 {
				game_over_timer = 3.0
				game_over = false
			}
		}

		// if using the slice of directions driver_moves to drive the game
		if !input_state.Interacted {

			proposed_index := int(math.Floor(game_time / float64(move_time)))

			if driver_index < proposed_index && driver_index < len(driver_moves) && !game_over {
				input_state.Pressed = false
				ani_state.State1 = g
				if g.Update(driver_moves[driver_index]) {
					// it's game over but we dont draw end screen
					game_over = true
				}
				ani_state.State2 = g
				ani_state.T = 0
				driver_index++
			}

			if game_over || driver_index >= len(driver_moves) {
				delay_on_gameover -= 0.015 //millis
				if delay_on_gameover <= 0 {
					input_state.Quit = true
				}
			}
		}
		// drawing the game

		need_animate := ani_state.update()
		did_combine := ani_state.did_combine_lately()

		window.MakeContextCurrent()

		gl.ClearColor(0.1, 0.1, 0.1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		if !game_over || !input_state.Interacted {

			for j := range 4 {
				offset_y := int32(TILE_SIZE * j)
				for i := range 4 {
					real_size := int32(TILE_SIZE - TILE_PADDING*2)
					offset_x := int32(TILE_SIZE * i)

					if need_animate && did_combine[j][i] {
						// blow it up when animating
						var bloat int32
						if ani_state.T > 0.5 {
							bloat = int32((1 - ani_state.T) * float64(bloat_size))
						} else {
							bloat = int32(ani_state.T * float64(bloat_size))
						}
						real_size += bloat * 2
						Draw_Tile(
							MARGIN+offset_x+TILE_PADDING-bloat,
							MARGIN+offset_y+TILE_PADDING-bloat,
							real_size,
							real_size,
							g.Board[j][i])
					} else {
						Draw_Tile(MARGIN+offset_x+TILE_PADDING, MARGIN+offset_y+TILE_PADDING, real_size, real_size, g.Board[j][i])
					}
				}
			}
		} else {
			draw_texture(0, 0, W, H, texture_game_over)
		}
		window.SwapBuffers()
		glfw.PollEvents()
		time.Sleep(time.Millisecond * 10)
		game_time += 0.015 // 15 millisec neuristic assumnig 5 millisec per iteration.. should calculate a real dt per frame!
	}
}
