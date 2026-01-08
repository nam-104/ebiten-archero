package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawSkillMenu(screen *ebiten.Image, skills []Skill) {
	// Vẽ lớp phủ tối
	vector.DrawFilledRect(
		screen,
		0, 0,
		960, 540,
		color.RGBA{0, 0, 0, 180},
		false,
	)

	for i, s := range skills {
		x := float32(150 + i*250)
		y := float32(150)

		// Vẽ khung ô kỹ năng
		vector.DrawFilledRect(
			screen,
			x, y,
			200, 250,
			color.RGBA{40, 40, 80, 255},
			false,
		)

		// Vẽ text
		ebitenutil.DebugPrintAt(screen, s.Name, int(x+50), int(y+20))
		ebitenutil.DebugPrintAt(screen, s.Description, int(x+20), int(y+100))
		ebitenutil.DebugPrintAt(
			screen,
			"Nhan phim "+string(rune('1'+i)),
			int(x+60),
			int(y+200),
		)
	}
}
