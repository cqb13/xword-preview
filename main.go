package main

import (
	"fmt"
	"os"
	"strconv"

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

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.SetTraceLogLevel(rl.LogNone)
	rl.SetTargetFPS(60)
	rl.InitWindow(1440, 810, fmt.Sprintf("%s | %s", puzzle.Title, puzzle.Author))
	rl.SetWindowMinSize(720, 405)
	defer rl.CloseWindow()

	halfHeight := float32(rl.GetScreenHeight() / 2)
	thirdWidth := float32(rl.GetScreenWidth() / 3)
	startBoardDrawX := thirdWidth + PADDING

	widthForBoard := rl.GetScreenWidth() - int(startBoardDrawX)
	cellDim := min(widthForBoard, rl.GetScreenHeight()) / max(puzzle.Board.Width(), puzzle.Board.Height())

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
		thirdWidth = float32(rl.GetScreenWidth() / 3)
		startBoardDrawX = thirdWidth + PADDING*2

		widthForBoard = rl.GetScreenWidth() - int(startBoardDrawX)
		cellDim = min(widthForBoard-PADDING, rl.GetScreenHeight()-PADDING*2) / max(puzzle.Board.Width(), puzzle.Board.Height())

		downCluesPanelRec = rl.NewRectangle(PADDING, PADDING, thirdWidth, halfHeight-PADDING*2)

		acrossCluesPanelRec = rl.NewRectangle(PADDING, halfHeight, thirdWidth, halfHeight-PADDING)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		// Down Clues
		rlgui.ScrollPanel(downCluesPanelRec, "Down", downCluesPanelContentRec, &downViewPanelScroll, &downCluesPanelView)
		rl.BeginScissorMode(int32(downCluesPanelView.X), int32(downCluesPanelView.Y), int32(downCluesPanelView.Width), int32(downCluesPanelView.Height))

		maxWidth := int32(0)
		height := int32(0)
		for _, clue := range downClues {
			text := fmt.Sprintf("%d. %s", clue.Num, clue.Clue)

			width := rl.MeasureText(text, FONT_SIZE)
			maxWidth = max(width, maxWidth)

			_ = width

			rl.DrawText(text, int32(downCluesPanelView.X+downViewPanelScroll.X), int32(downCluesPanelView.Y+downViewPanelScroll.Y)+height, FONT_SIZE, rl.Black)
			height += FONT_SIZE
		}
		downCluesPanelContentRec.Width = float32(maxWidth + PADDING)

		rl.EndScissorMode()

		// Across Clues
		rlgui.ScrollPanel(acrossCluesPanelRec, "Across", acrossCluesPanelContentRec, &acrossViewPanelScroll, &acrossCluesPanelView)
		rl.BeginScissorMode(int32(acrossCluesPanelView.X), int32(acrossCluesPanelView.Y), int32(acrossCluesPanelView.Width), int32(acrossCluesPanelView.Height))

		maxWidth = 0
		height = int32(0)
		for _, clue := range acrossClues {
			text := fmt.Sprintf("%d. %s", clue.Num, clue.Clue)

			width := rl.MeasureText(text, FONT_SIZE)
			maxWidth = max(width, maxWidth)

			_ = width

			rl.DrawText(text, int32(acrossCluesPanelView.X+acrossViewPanelScroll.X), int32(acrossCluesPanelView.Y+acrossViewPanelScroll.Y)+height, FONT_SIZE, rl.Black)
			height += FONT_SIZE
		}

		acrossCluesPanelContentRec.Width = float32(maxWidth + PADDING)

		rl.EndScissorMode()

		// Draw board
		x := startBoardDrawX
		y := PADDING

		letterFontSize := cellDim / 2
		numberFontSize := cellDim / 4

		nextNum := 1

		// so the edge thickness is the same as lines inside the grid
		rl.DrawRectangleLinesEx(rl.NewRectangle(startBoardDrawX, PADDING, float32(cellDim*puzzle.Board.Width()+1), float32(cellDim*puzzle.Board.Height())+1), 2, rl.Black)

		for cy, row := range puzzle.Board {
			for cx, cell := range row {
				if cell.Value == puz.DIAGRAMLESS_SOLID_SQUARE {
					x += float32(cellDim)
					continue
				}

				if cell.Value == puz.SOLID_SQUARE {
					rl.DrawRectangle(int32(x), int32(y), int32(cellDim), int32(cellDim), rl.Black)
					x += float32(cellDim)
					continue
				}

				rl.DrawRectangleLines(int32(x), int32(y), int32(cellDim), int32(cellDim), rl.Black)

				w := rl.MeasureText(string(cell.Value), int32(letterFontSize))

				rl.DrawText(string(cell.Value), int32(x)+int32(cellDim/2)-w/2, int32(y)+int32(letterFontSize)/2, int32(letterFontSize), rl.Black)

				if puzzle.Board.StartsAcrossWord(cx, cy) || puzzle.Board.StartsDownWord(cx, cy) {
					rl.DrawText(strconv.Itoa(nextNum), int32(x+3), int32(y+3), int32(numberFontSize), rl.Gray)
					nextNum++
				}

				x += float32(cellDim)
			}
			x = startBoardDrawX
			y += cellDim
		}

		rl.EndDrawing()
	}
}
