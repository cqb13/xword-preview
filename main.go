package main

import (
	"fmt"
	"os"

	"github.com/cqb13/puz-parser/puz"
	rlgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const PADDING = 10
const FONT_SIZE = 20

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Not enough args. xword-preview [crossword]")
	}

	bytes, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("Failed to load crossword: ", err)
		os.Exit(1)
	}

	puzzle, err := puz.DecodePuz(bytes)
	if err != nil {
		fmt.Println("Failed to decode crossword: ", err)
		os.Exit(1)
	}

	// rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.SetTraceLogLevel(rl.LogNone)
	rl.SetTargetFPS(60)
	rl.InitWindow(1440, 810, fmt.Sprintf("%s | %s", puzzle.Title, puzzle.Author))
	defer rl.CloseWindow()

	halfHeight := float32(rl.GetScreenHeight() / 2)
	thirdWidth := float32(rl.GetScreenWidth() / 3)

	downClues := puzzle.GetCluesByDirection(puz.Down)
	acrossClues := puzzle.GetCluesByDirection(puz.Across)

	var (
		downCluesPanelRec        = rl.NewRectangle(PADDING, PADDING, thirdWidth, halfHeight)
		downCluesPanelContentRec = rl.NewRectangle(PADDING, PADDING, float32(rl.GetMonitorWidth(0)), float32(FONT_SIZE*len(downClues)))
		downCluesPanelView       = rl.NewRectangle(0, 0, 0, 0)
		downViewPanelScroll      = rl.NewVector2(99, 0)

		acrossCluesPanelRec        = rl.NewRectangle(PADDING, PADDING, thirdWidth, halfHeight)
		acrossCluesPanelContentRec = rl.NewRectangle(PADDING, PADDING, float32(rl.GetMonitorWidth(0)), float32(FONT_SIZE*len(acrossClues)))
		acrossCluesPanelView       = rl.NewRectangle(0, 0, 0, 0)
		acrossViewPanelScroll      = rl.NewVector2(99, 0)
	)

	for !rl.WindowShouldClose() {
		halfHeight = float32(rl.GetScreenHeight() / 2)
		thirdWidth := float32(rl.GetScreenWidth() / 3)

		downCluesPanelRec = rl.NewRectangle(PADDING, PADDING, thirdWidth, halfHeight-PADDING*2)

		acrossCluesPanelRec = rl.NewRectangle(PADDING, halfHeight, thirdWidth, halfHeight-PADDING)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		// Down Clues
		rlgui.ScrollPanel(downCluesPanelRec, "Down", downCluesPanelContentRec, &downViewPanelScroll, &downCluesPanelView)
		rl.BeginScissorMode(int32(downCluesPanelView.X), int32(downCluesPanelView.Y), int32(downCluesPanelView.Width), int32(downCluesPanelView.Height))

		height := int32(0)
		for _, clue := range downClues {
			text := fmt.Sprintf("%d. %s", clue.Num, clue.Clue)

			width := rl.MeasureText(text, FONT_SIZE)

			_ = width

			rl.DrawText(text, int32(downCluesPanelView.X+downViewPanelScroll.X), int32(downCluesPanelView.Y+downViewPanelScroll.Y)+height, FONT_SIZE, rl.Black)
			height += FONT_SIZE
		}

		rl.EndScissorMode()

		// Across Clues
		rlgui.ScrollPanel(acrossCluesPanelRec, "Across", acrossCluesPanelContentRec, &acrossViewPanelScroll, &acrossCluesPanelView)
		rl.BeginScissorMode(int32(acrossCluesPanelView.X), int32(acrossCluesPanelView.Y), int32(acrossCluesPanelView.Width), int32(acrossCluesPanelView.Height))

		height = int32(0)
		for _, clue := range acrossClues {
			text := fmt.Sprintf("%d. %s", clue.Num, clue.Clue)

			width := rl.MeasureText(text, FONT_SIZE)

			_ = width

			rl.DrawText(text, int32(acrossCluesPanelView.X+acrossViewPanelScroll.X), int32(acrossCluesPanelView.Y+acrossViewPanelScroll.Y)+height, FONT_SIZE, rl.Black)
			height += FONT_SIZE
		}

		rl.EndScissorMode()

		rl.EndDrawing()
	}
}
