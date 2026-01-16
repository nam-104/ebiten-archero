package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Player đại diện cho nhân vật người chơi
type Player struct {
	X, Y         float64
	Health       float64
	MaxHealth    float64
	Speed        float64 // px mỗi frame
	AttackDamage float64
	AttackSpeed  float64
	AttackTimer  float64
	Img          *ebiten.Image
	Width        float64
	Height       float64
	Skills       []Skill
}

// NewPlayer tạo player mới
func NewPlayer(img *ebiten.Image, x, y, maxHealth, speed, attackDamage, attackSpeed float64) *Player {
	// Fix speed cứng ở đây nếu muốn mặc định (vd 3.2 mượt hơn, game archero thật thường > 3)
	if speed < 2.7 {
		speed = 3.2
	}
	return &Player{
		X:            x,
		Y:            y,
		Health:       maxHealth,
		MaxHealth:    maxHealth,
		Speed:        speed,
		AttackDamage: attackDamage,
		AttackSpeed:  attackSpeed,
		AttackTimer:  0.0,
		Img:          img,
		Width:        16.0,
		Height:       16.0,
	}
}

// Update cập nhật trạng thái player
func (p *Player) Update() {
	if p.AttackTimer > 0 {
		p.AttackTimer -= 1.0 / 60.0 // Giả sử 60 FPS
	}
}

// CanAttack kiểm tra xem player có thể tấn công không
func (p *Player) CanAttack() bool {
	return p.AttackTimer <= 0
}

// Attack thực hiện tấn công và reset timer
func (p *Player) Attack() {
	p.AttackTimer = 1.0 / p.AttackSpeed
}

// Move di chuyển player mượt mà hơn
func (p *Player) Move(dx, dy float64, mapWidth, mapHeight float64) {
	// Tính toán vị trí mới tiềm năng
	newX := p.X + dx*p.Speed
	newY := p.Y + dy*p.Speed

	// Kiểm tra và cập nhật X riêng biệt
	if newX >= 0 && newX <= mapWidth-p.Width {
		p.X = newX
	}

	// Kiểm tra và cập nhật Y riêng biệt
	if newY >= 0 && newY <= mapHeight-p.Height {
		p.Y = newY
	}
}

// GetCenter trả về tọa độ trung tâm của player
func (p *Player) GetCenter() (float64, float64) {
	return p.X + p.Width/2, p.Y + p.Height/2
}

// Draw vẽ player lên màn hình
func (p *Player) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X-cameraX, p.Y-cameraY)

	// Vẽ sprite từ spritesheet (16x16 đầu tiên)
	screen.DrawImage(
		p.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
		opts,
	)
	DrawHealthBar(screen, p.X-cameraX, p.Y-cameraY-8, p.Width, 3, p.Health/p.MaxHealth, false)
}

// GetDistanceTo tính khoảng cách đến một điểm
func (p *Player) GetDistanceTo(x, y float64) float64 {
	px, py := p.GetCenter()
	dx := x - px
	dy := y - py
	return math.Sqrt(dx*dx + dy*dy)
}

// TakeDamage nhận sát thương
func (p *Player) TakeDamage(amount float64) {
	p.Health -= amount
	if p.Health < 0 {
		p.Health = 0
	}
}

// IsAlive kiểm tra player còn sống không
func (p *Player) IsAlive() bool {
	return p.Health > 0
}

// CheckCollision kiểm tra va chạm giữa Player và một thực thể khác (AABB)
func (p *Player) CheckCollision(otherX, otherY, otherW, otherH float64) bool {
	return p.X < otherX+otherW &&
		p.X+p.Width > otherX &&
		p.Y < otherY+otherH &&
		p.Y+p.Height > otherY
}

// Hàm để Player học kỹ năng mới
func (p *Player) LearnSkill(s Skill) {
	p.Skills = append(p.Skills, s)

	// Áp dụng ngay lập tức nếu là kỹ năng tăng chỉ số
	switch s.Type {
	case AttackBoost:
		p.AttackDamage *= 1.2
	case SpeedBoost:
		p.Speed += 0.5
	}
}

// Hàm kiểm tra xem Player có sở hữu 1 loại kỹ năng nào đó không
func (p *Player) HasSkill(t SkillType) bool {
	for _, s := range p.Skills {
		if s.Type == t {
			return true
		}
	}
	return false
}

// GetSkillCount trả về số lượng kỹ năng cùng loại mà Player sở hữu
func (p *Player) GetSkillCount(t SkillType) int {
	count := 0
	for _, s := range p.Skills {
		if s.Type == t {
			count++
		}
	}
	return count
}
