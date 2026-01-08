package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Enemy đại diện cho quái vật
type Enemy struct {
	X, Y       float64
	Health     float64
	MaxHealth  float64
	Speed      float64
	Damage     float64
	Img        *ebiten.Image
	Width      float64
	Height     float64
	Active     bool
	FollowDist float64 // Khoảng cách bắt đầu đuổi theo player
	State      int     // 0: Đứng nghỉ, 1: Lao tới
	Timer      float64 // Bộ đếm thời gian cho trạng thái hiện tại
}

const (
	StateRest = 0
	StateDash = 1
)

// NewEnemy tạo enemy mới
func NewEnemy(img *ebiten.Image, x, y, maxHealth, speed, damage, followDist float64) *Enemy {
	return &Enemy{
		X:          x,
		Y:          y,
		Health:     maxHealth,
		MaxHealth:  maxHealth,
		Speed:      speed,
		Damage:     damage,
		Img:        img,
		Width:      16.0,
		Height:     16.0,
		Active:     true,
		FollowDist: followDist,
		State:      0,
		Timer:      1.0, // 1 giây sau khi sinh ra mới bắt đầu lao tới
	}
}

func (e *Enemy) Update(playerX, playerY float64, screenWidth, screenHeight float64) {
	if !e.Active || e.Health <= 0 {
		e.Active = false
		return
	}

	// 1. Cập nhật bộ đếm thời gian
	e.Timer -= 1.0 / 60.0 // Giả sử game chạy 60 FPS

	// 2. Kiểm tra đổi trạng thái
	if e.Timer <= 0 {
		if e.State == 0 { // Nếu đang nghỉ -> Chuyển sang Lao tới
			e.State = 1
			e.Timer = 1.2 // Lao tới trong 1.2 giây
		} else { // Nếu đang lao tới -> Chuyển sang Nghỉ
			e.State = 0
			e.Timer = 0.8 // Nghỉ trong 0.8 giây
		}
	}

	// 3. Xử lý di chuyển dựa trên trạng thái
	if e.State == 1 { // Chỉ di chuyển khi ở trạng thái Lao tới
		dx := playerX - e.X
		dy := playerY - e.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// Chỉ lao tới nếu trong tầm nhìn (FollowDist)
		if distance < e.FollowDist && distance > 0 {
			dx /= distance
			dy /= distance

			// Khi lao tới, có thể tăng tốc độ lên một chút (ví dụ e.Speed * 1.5)
			speedMultiplier := 1.8
			newX := e.X + dx*e.Speed*speedMultiplier
			newY := e.Y + dy*e.Speed*speedMultiplier

			// Giới hạn trong bản đồ
			if newX >= 0 && newX <= screenWidth-e.Width {
				e.X = newX
			}
			if newY >= 0 && newY <= screenHeight-e.Height {
				e.Y = newY
			}
		}
	}
	// Nếu e.State == 0 (Nghỉ), quái sẽ đứng yên tại chỗ vì code không thay đổi e.X, e.Y
}

// GetCenter trả về tọa độ trung tâm của enemy
func (e *Enemy) GetCenter() (float64, float64) {
	return e.X + e.Width/2, e.Y + e.Height/2
}

// TakeDamage nhận sát thương
func (e *Enemy) TakeDamage(amount float64) {
	e.Health -= amount
	if e.Health <= 0 {
		e.Health = 0
		e.Active = false
	}
}

// IsAlive kiểm tra enemy còn sống không
func (e *Enemy) IsAlive() bool {
	return e.Active && e.Health > 0
}

// Draw vẽ enemy lên màn hình
func (e *Enemy) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	if !e.IsAlive() {
		return
	}

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(e.X-cameraX, e.Y-cameraY)

	// Vẽ sprite từ spritesheet (16x16 đầu tiên)
	screen.DrawImage(
		e.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
		opts,
	)
	DrawHealthBar(screen, e.X-cameraX, e.Y-cameraY-5, e.Width, 2, e.Health/30.0, true)
}

// GetDistanceTo tính khoảng cách đến một điểm
func (e *Enemy) GetDistanceTo(x, y float64) float64 {
	ex, ey := e.GetCenter()
	dx := x - ex
	dy := y - ey
	return math.Sqrt(dx*dx + dy*dy)
}

// CheckCollision kiểm tra va chạm với player
func (e *Enemy) CheckCollision(px, py, pw, ph float64) bool {
	return e.X < px+pw &&
		e.X+e.Width > px &&
		e.Y < py+ph &&
		e.Y+e.Height > py
}
