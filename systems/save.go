package systems

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// GameData lưu trữ dữ liệu game
type GameData struct {
	Level        int     `json:"level"`
	Experience   int     `json:"experience"`
	Gold         int     `json:"gold"`
	MaxHealth    float64 `json:"maxHealth"`
	AttackDamage float64 `json:"attackDamage"`
	AttackSpeed  float64 `json:"attackSpeed"`
	PlayerX      float64 `json:"playerX"`
	PlayerY      float64 `json:"playerY"`
}

const saveFilePath = "save.json"

// LoadGameData tải dữ liệu game từ file
func LoadGameData() (*GameData, error) {
	data := &GameData{
		Level:        1,
		Experience:   0,
		Gold:         0,
		MaxHealth:    100.0,
		AttackDamage: 10.0,
		AttackSpeed:  1.0,
		PlayerX:      160.0,
		PlayerY:      120.0,
	}

	file, err := os.ReadFile(saveFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Nếu file không tồn tại, trả về dữ liệu mặc định
			return data, nil
		}
		return nil, err
	}

	err = json.Unmarshal(file, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SaveGameData lưu dữ liệu game vào file
func SaveGameData(data *GameData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Đảm bảo thư mục tồn tại
	dir := filepath.Dir(saveFilePath)
	if dir != "." && dir != "" {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(saveFilePath, jsonData, 0644)
}
