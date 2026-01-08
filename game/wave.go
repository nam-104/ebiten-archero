package game

import (
	"math"
	"math/rand"
)

// WaveManager quản lý các wave quái
type WaveManager struct {
	CurrentWave    int
	EnemiesPerWave int
	EnemiesSpawned int
	WaveComplete   bool
	SpawnTimer     float64
	SpawnInterval  float64
	ScreenWidth    float64
	ScreenHeight   float64
}

// NewWaveManager tạo wave manager mới
func NewWaveManager(screenWidth, screenHeight float64) *WaveManager {
	return &WaveManager{
		CurrentWave:    1,
		EnemiesPerWave: 5,
		EnemiesSpawned: 0,
		WaveComplete:   false,
		SpawnTimer:     0.0,
		SpawnInterval:  1.0, // Spawn mỗi 1 giây
		ScreenWidth:    screenWidth,
		ScreenHeight:   screenHeight,
	}
}

// Update cập nhật wave manager
func (wm *WaveManager) Update() {
	if wm.WaveComplete {
		return
	}

	wm.SpawnTimer += 1.0 / 60.0

	if wm.SpawnTimer >= wm.SpawnInterval && wm.EnemiesSpawned < wm.EnemiesPerWave {
		wm.SpawnTimer = 0.0
		wm.EnemiesSpawned++
	}

	// Kiểm tra wave hoàn thành (khi đã spawn hết quái)
	if wm.EnemiesSpawned >= wm.EnemiesPerWave {
		// Wave chỉ hoàn thành khi tất cả quái đã bị tiêu diệt (được check ở game loop)
	}
}

// GetSpawnPosition trả về vị trí spawn quái ngẫu nhiên xung quanh player
func (w *WaveManager) GetSpawnPosition(playerX, playerY float64) (float64, float64) {
	// Góc ngẫu nhiên (0 đến 360 độ)
	angle := rand.Float64() * 2 * math.Pi

	// Khoảng cách từ player (ví dụ: cách player từ 300 đến 500 pixel)
	// Khoảng cách này phải lớn hơn một nửa màn hình để quái xuất hiện từ rìa
	distance := 150.0 + rand.Float64()*100.0

	spawnX := playerX + math.Cos(angle)*distance
	spawnY := playerY + math.Sin(angle)*distance

	return spawnX, spawnY
}

// StartNextWave bắt đầu wave tiếp theo
func (wm *WaveManager) StartNextWave() {
	wm.CurrentWave++
	wm.EnemiesPerWave = 5 + wm.CurrentWave*2 // Tăng số quái mỗi wave
	wm.EnemiesSpawned = 0
	wm.WaveComplete = false
	wm.SpawnTimer = 0.0
	wm.SpawnInterval = math.Max(0.3, 1.0-float64(wm.CurrentWave)*0.05) // Spawn nhanh hơn theo wave
}

// Reset reset về wave 1
func (wm *WaveManager) Reset() {
	wm.CurrentWave = 1
	wm.EnemiesPerWave = 5
	wm.EnemiesSpawned = 0
	wm.WaveComplete = false
	wm.SpawnTimer = 0.0
	wm.SpawnInterval = 1.0
}
