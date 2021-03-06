package tictactoe

import (
	"fmt"
	"os"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	"github.com/juanvallejo/go-tictactoe/pkg/tictactoe/grid"
	"github.com/juanvallejo/go-tictactoe/pkg/tictactoe/score"
	"github.com/juanvallejo/go-tictactoe/pkg/tictactoe/shape"
)

const (
	winWidth  = 800
	winHeight = 600

	cellMargin  = 60
	shapeMargin = 25

	winTextSize   = 4
	scoreTextSize = 2

	scoreMarginX = 10
	scoreMarginY = 5
)

var winBgcolor = colornames.Darkslategrey
var winTextAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

func NewGame() {
	config := pixelgl.WindowConfig{
		Title:  "Tic Tac Toe",
		Bounds: pixel.R(0, 0, winWidth, winHeight),
		VSync:  true,
	}

	window, err := pixelgl.NewWindow(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	window.Clear(winBgcolor)

	gameWin := false
	scoreKeeper := score.ScoreKeeper(make(map[string]int))
	shapeDecider := shape.NewShapeDecider(shape.CrossShape)
	bounds := window.Bounds()
	context := imdraw.New(nil)
	winTextContext := text.New(pixel.V(bounds.Max.X/2, bounds.Max.Y/2), winTextAtlas)
	scoreTextContext := text.New(pixel.V(bounds.Min.X, bounds.Max.Y), winTextAtlas)

	scoreRenderer := score.NewScoreRenderer()
	scoreRenderer.RenderFunc(func(ctx *text.Text, scores score.ScoreKeeper) {
		ctx.Clear()

		ctx.LineHeight = 0
		ctx.Dot.X += scoreMarginX
		ctx.Dot.Y -= scoreMarginY

		text := fmt.Sprintf("%s: %d", shape.CrossShape, scores.Get(string(shape.CrossShape)))
		ctx.Dot.Y -= ctx.BoundsOf(text).H()
		fmt.Fprintf(ctx, "%s\n", text)

		text = fmt.Sprintf("%s: %d", shape.CircleShape, scores.Get(string(shape.CircleShape)))
		ctx.Dot.X = bounds.Max.X/2 - ctx.BoundsOf(text).W() - scoreMarginX
		fmt.Fprintf(ctx, "%s\n", text)
	})

	g := grid.NewGrid(pixel.V(0, 0), bounds.Max.X, bounds.Max.Y, grid.MaxCells, cellMargin)

	for !window.Closed() {
		if !gameWin {
			window.Clear(winBgcolor)
			context.Clear()
			winTextContext.Clear()
			scoreTextContext.Clear()
		}

		if window.JustPressed(pixelgl.MouseButtonLeft) {
			gameWin = handleMouseClick(window, context, winTextContext, g, shapeDecider, scoreKeeper, bounds, gameWin)
		}

		g.Render(context)
		scoreRenderer.Render(scoreTextContext, scoreKeeper)
		context.Draw(window)
		winTextContext.Draw(window, pixel.IM.Scaled(winTextContext.Orig, winTextSize))
		scoreTextContext.Draw(window, pixel.IM.Scaled(scoreTextContext.Orig, scoreTextSize))
		window.Update()
	}
}

func handleMouseClick(window *pixelgl.Window, context *imdraw.IMDraw, winTextContext *text.Text, g grid.Grid, shapeDecider *shape.ShapeDecider, scoreKeeper score.ScoreKeeper, bounds pixel.Rect, gameWin bool) bool {
	if gameWin {
		window.Clear(winBgcolor)
		g.Reset()
		return false
	}

	if cell := g.AtVector(window.MousePosition()); cell != nil {
		cell.Set(shape.NewShape(pixel.V(cell.Start().X, cell.Start().Y), shapeDecider.Next(), (bounds.Max.X-cellMargin*2)/grid.MaxCells, (bounds.Max.Y-cellMargin*2)/grid.MaxCells, shapeMargin))
		return g.CheckWin(context, winTextContext, scoreKeeper)
	}

	return false
}
