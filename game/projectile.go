package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Projectile đại diện cho đạn
type Projectile struct {
	X, Y        float64
	VX, VY      float64
	Speed       float64
	Damage      float64
	Active      bool
	Img         *ebiten.Image
	Width       float64
	Height      float64
	LifeTime    float64
	MaxLifeTime float64
	IsPiercing  bool
}

// NewProjectile tạo projectile mới
func NewProjectile(img *ebiten.Image, x, y, targetX, targetY, speed, damage float64) *Projectile {
	dx := targetX - x
	dy := targetY - y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance == 0 {
		distance = 1
	}

	vx := (dx / distance) * speed
	vy := (dy / distance) * speed

	return &Projectile{
		X:           x,
		Y:           y,
		VX:          vx,
		VY:          vy,
		Speed:       speed,
		Damage:      damage,
		Active:      true,
		Img:         img,
		Width:       16.0,
		Height:      16.0,
		LifeTime:    0.0,
		MaxLifeTime: 5.0, // 5 giây
	}
}

// Update cập nhật trạng thái projectile
func (p *Projectile) Update(screenWidth, screenHeight float64) {
	if !p.Active {
		return
	}

	p.X += p.VX
	p.Y += p.VY
	p.LifeTime += 1.0 / 60.0

	// Kiểm tra ra ngoài màn hình hoặc hết thời gian
	if p.X < -p.Width || p.X > screenWidth+p.Width ||
		p.Y < -p.Height || p.Y > screenHeight+p.Height ||
		p.LifeTime >= p.MaxLifeTime {
		p.Active = false
	}
}

// Draw vẽ projectile lên màn hình
func (p *Projectile) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	if !p.Active || p.Img == nil {
		return
	}

	opts := &ebiten.DrawImageOptions{}

	// Tính góc xoay dựa trên vector vận tốc
	angle := math.Atan2(p.VY, p.VX)

	w, h := p.Img.Bounds().Dx(), p.Img.Bounds().Dy()

	// Dời tâm về giữa ảnh để xoay
	opts.GeoM.Translate(-float64(w)/2, -float64(h)/2)

	// Scale nhỏ lại 0.5
	opts.GeoM.Scale(0.5, 0.5)

	// Xoay ảnh (giả sử ảnh gốc mũi tên hướng sang PHẢI -> 0 độ)
	// Nếu nó hướng lên thì +Pi/2. Nếu hướng chéo thì +Pi/4.
	// User report: "nằm ngang" (sai hướng). Thử bỏ offset (giả sử asset gốc đã xoay đúng hoặc là hướng phải).
	// Nếu Asset là Arrow01(32x32), thường là chéo 45 độ (Up-Right).
	// Hãy thử -Pi/4 (để xoay nó về 0 rồi +angle) nếu nó là chéo.
	// Nhưng user bảo "nằm ngang", có thể nó đang bị xoay 90 độ.
	// Thử dùng angle thuần túy trước.
	opts.GeoM.Rotate(angle)

	// Dời về vị trí hiển thị (tâm của projectile)
	opts.GeoM.Translate(p.X+p.Width/2-cameraX, p.Y+p.Height/2-cameraY)

	screen.DrawImage(p.Img, opts)
}

// CheckCollision kiểm tra va chạm với enemy
func (p *Projectile) CheckCollision(ex, ey, ew, eh float64) bool {
	return p.X < ex+ew &&
		p.X+p.Width > ex &&
		p.Y < ey+eh &&
		p.Y+p.Height > ey
}
