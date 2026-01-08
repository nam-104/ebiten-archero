package game

import (
	"image"
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
		Width:       8.0,
		Height:      8.0,
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
	if !p.Active {
		return
	}

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X-cameraX, p.Y-cameraY)

	if p.Img != nil {
		// Vẽ sprite nếu có
		screen.DrawImage(
			p.Img.SubImage(image.Rect(0, 0, 8, 8)).(*ebiten.Image),
			opts,
		)
	}
}

// CheckCollision kiểm tra va chạm với enemy
func (p *Projectile) CheckCollision(ex, ey, ew, eh float64) bool {
	return p.X < ex+ew &&
		p.X+p.Width > ex &&
		p.Y < ey+eh &&
		p.Y+p.Height > ey
}
