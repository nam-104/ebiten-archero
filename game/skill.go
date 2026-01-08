package game

// Loại kỹ năng
type SkillType int

const (
	AttackBoost   SkillType = iota // Tăng sát thương
	SpeedBoost                     // Tăng tốc độ chạy
	Multishot                      // Bắn 2 tia
	PiercingShot                   // Đạn xuyên thấu
	DiagonalArrow                  // Bắn chéo
)

type Skill struct {
	Type        SkillType
	Name        string
	Description string
}

// Danh sách tất cả kỹ năng có trong game để random
var AllSkills = []Skill{
	{Type: AttackBoost, Name: "AttackBoost +", Description: "+20% ATK"},
	{Type: SpeedBoost, Name: "SpeedBoost +", Description: "SpeedBoost"},
	{Type: Multishot, Name: "Multishot", Description: "Multishot"},
	{Type: PiercingShot, Name: "PiercingShot", Description: "PiercingShot"},    // Thêm cho đủ bộ
	{Type: DiagonalArrow, Name: "DiagonalArrow", Description: "DiagonalArrow"}, // THÊM DÒNG NÀY
}
