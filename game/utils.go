package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DrawHealthBar là hàm dùng chung để vẽ thanh máu cho bất kỳ thực thể nào
// Lưu ý: Viết hoa chữ cái đầu (D) để các file khác trong hoặc ngoài package có thể gọi được
func DrawHealthBar(screen *ebiten.Image, x, y, width, height float64, ratio float64, isEnemy bool) {
	// 1. Vẽ nền (màu đen hoặc đỏ tối)
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(width), float32(height), color.RGBA{60, 0, 0, 255}, false)

	// 2. Chọn màu dựa trên việc đó là quái hay người chơi
	healthColor := color.RGBA{0, 255, 100, 255} // Màu xanh cho Player
	if isEnemy {
		healthColor = color.RGBA{255, 200, 0, 255} // Màu vàng/cam cho Quái
	}

	// 3. Vẽ phần máu hiện tại
	if ratio > 0 {
		if ratio > 1 {
			ratio = 1
		}
		vector.DrawFilledRect(screen, float32(x), float32(y), float32(width*ratio), float32(height), healthColor, false)
	}
}
