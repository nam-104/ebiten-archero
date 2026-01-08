package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand/v2"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"pixcel-game/game"
	g "pixcel-game/game"
	"pixcel-game/systems"
)

const (
	screenWidth  = 960
	screenHeight = 540
)

// đường dẫn assets lấy từ dự án rpg-in-golang
const assetsBase = "assets"

const TeleportGateTileID = 159 // có thể cần đổi lại đúng tile cổng trong map sau

const (
	StatePlaying = iota
	StateSkillSelect
)

type ArcheroGame struct {
	player              *g.Player
	enemies             []*g.Enemy
	projectiles         []*g.Projectile
	potions             []*Potion
	wave                *g.WaveManager
	camera              *systems.Camera
	tilemap             *g.TilemapJSON
	tilesetImg          *ebiten.Image
	playerImg           *ebiten.Image
	enemyImg            *ebiten.Image
	projectileImg       *ebiten.Image
	potionImg           *ebiten.Image
	saveData            *systems.GameData
	mapWidthPx          float64
	mapHeightPx         float64
	gameState           int          // Lưu trạng thái hiện tại
	currentSkillOptions []game.Skill // Các kỹ năng đang hiển thị để chọn
}

type Potion struct {
	X, Y   float64
	Img    *ebiten.Image
	Width  float64
	Height float64
}

func (p *Potion) draw(screen *ebiten.Image, cameraX, cameraY float64) {
	if p.Img == nil {
		return
	}
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X-cameraX, p.Y-cameraY)
	screen.DrawImage(p.Img, opts)
}

func NewArcheroGame() *ArcheroGame {
	data, err := systems.LoadGameData()
	if err != nil {
		log.Printf("khong load duoc save, dung default: %v", err)
		data = &systems.GameData{
			Level:        1,
			Experience:   0,
			Gold:         0,
			MaxHealth:    100,
			AttackDamage: 10,
			AttackSpeed:  1,
			PlayerX:      160,
			PlayerY:      120,
		}
	}

	playerImg, _, err := ebitenutil.NewImageFromFile(filepath.Join(assetsBase, "images", "ninja.png"))
	if err != nil {
		log.Fatal(err)
	}
	enemyImg, _, err := ebitenutil.NewImageFromFile(filepath.Join(assetsBase, "images", "skeleton.png"))
	if err != nil {
		log.Fatal(err)
	}
	tilesetImg, _, err := ebitenutil.NewImageFromFile(filepath.Join(assetsBase, "images", "TilesetFloor.png"))
	if err != nil {
		log.Fatal(err)
	}
	projectileImg, _, err := ebitenutil.NewImageFromFile(filepath.Join(assetsBase, "images", "shuriken.png"))
	if err != nil {
		log.Printf("khong load duoc projectile img, su dung nil: %v", err)
	}
	potionImg, _, err := ebitenutil.NewImageFromFile(filepath.Join(assetsBase, "images", "potion.png"))
	if err != nil {
		log.Printf("khong load duoc potion img: %v", err)
	}

	tilemap, err := g.NewTilemapJSON(filepath.Join(assetsBase, "maps", "spawn.json"))
	if err != nil {
		log.Fatal(err)
	}

	game := &ArcheroGame{
		playerImg:     playerImg,
		enemyImg:      enemyImg,
		projectileImg: projectileImg,
		potionImg:     potionImg,
		tilesetImg:    tilesetImg,
		tilemap:       tilemap,
		saveData:      data,
	}

	game.mapWidthPx = float64(tilemap.Width * tilemap.TileW)
	game.mapHeightPx = float64(tilemap.Height * tilemap.TileH)

	game.resetStateFromSave()
	return game
}

func (gme *ArcheroGame) resetStateFromSave() {
	gme.player = g.NewPlayer(
		gme.playerImg,
		gme.saveData.PlayerX,
		gme.saveData.PlayerY,
		gme.saveData.MaxHealth,
		3.2,
		gme.saveData.AttackDamage,
		gme.saveData.AttackSpeed,
	)
	gme.enemies = []*g.Enemy{}
	gme.potions = []*Potion{}
	gme.projectiles = []*g.Projectile{}
	gme.wave = g.NewWaveManager(gme.mapWidthPx, gme.mapHeightPx)
	gme.camera = systems.NewCamera(screenWidth, screenHeight)
}

func (gme *ArcheroGame) Update() error {
	// Nếu nhấn phím L thì hiện menu kỹ năng (để test)
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		gme.gameState = StateSkillSelect
		// Lấy đại 3 kỹ năng đầu tiên trong AllSkills để hiện
		gme.currentSkillOptions = game.AllSkills[:3]
	}

	if gme.gameState == StateSkillSelect {
		gme.handleSkillSelection() // Hàm xử lý khi người chơi bấm 1, 2, 3
		return nil                 // Dừng các logic di chuyển/bắn đạn khi đang chọn kỹ năng
	}

	gme.handleTeleportGate()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	gme.handleMovement()
	gme.player.Update()
	gme.wave.Update()

	gme.spawnEnemiesIfNeeded()

	for _, e := range gme.enemies {
		prevAlive := e.IsAlive()
		e.Update(gme.player.X, gme.player.Y, gme.mapWidthPx, gme.mapHeightPx)
		if e.IsAlive() && e.CheckCollision(gme.player.X, gme.player.Y, gme.player.Width, gme.player.Height) {
			gme.player.TakeDamage(5)
		}
		// Nếu enemy vừa chết trong frame này thì có thể drop potion
		if prevAlive && !e.IsAlive() {
			if rand.Float64() < 0.3 {
				gme.spawnPotion(e.X, e.Y)
			}
		}
	}

	gme.handleAutoAttack()
	gme.updateProjectiles()
	gme.cleanupEntities()
	gme.handlePotions()

	gme.updateCamera()
	gme.handleWaveComplete()
	gme.handleSaveLoad()

	return nil
}

func (gme *ArcheroGame) handleSkillSelection() {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		gme.player.LearnSkill(gme.currentSkillOptions[0])
		gme.gameState = StatePlaying // Chọn xong thì quay lại chơi
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		gme.player.LearnSkill(gme.currentSkillOptions[1])
		gme.gameState = StatePlaying // Chọn xong thì quay lại chơi
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		gme.player.LearnSkill(gme.currentSkillOptions[2])
		gme.gameState = StatePlaying // Chọn xong thì quay lại chơi
	}
	// Tương tự cho phím 2 và 3...
}

func (gme *ArcheroGame) handleMovement() {
	var dx, dy float64

	// 1. Nhận diện cả WASD và phím mũi tên cho nhạy
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx += 1
	}

	if dx != 0 || dy != 0 {
		// 2. Giữ nguyên chuẩn hóa (Normalization) để đi chéo không bị nhanh quá mức
		norm := math.Sqrt(dx*dx + dy*dy)
		dx /= norm
		dy /= norm

		// 3. Tăng tốc độ truyền vào (nếu p.Speed đang thấp)
		// Bạn có thể thử nhân trực tiếp ở đây để test:
		// gme.player.Move(dx * 1.5, dy * 1.5, gme.mapWidthPx, gme.mapHeightPx)

		gme.player.Move(dx, dy, gme.mapWidthPx, gme.mapHeightPx)
	}
}

func (gme *ArcheroGame) handlePotions() {
	filteredPotions := gme.potions[:0]
	for _, p := range gme.potions {
		// Kiểm tra va chạm giữa Player và Potion
		if gme.player.CheckCollision(p.X, p.Y, p.Width, p.Height) {
			// Hồi máu cho player (ví dụ 20 HP), không vượt quá MaxHealth
			gme.player.Health += 20
			if gme.player.Health > gme.player.MaxHealth {
				gme.player.Health = gme.player.MaxHealth
			}
			log.Println("Đã ăn bình máu! HP hiện tại:", gme.player.Health)
			continue // Không thêm vào danh sách mới (tương đương với việc xóa)
		}
		filteredPotions = append(filteredPotions, p)
	}
	gme.potions = filteredPotions
}

func (gme *ArcheroGame) spawnPotion(x, y float64) {
	pot := &Potion{
		X:      x,
		Y:      y,
		Img:    gme.potionImg,
		Width:  16, // Điều chỉnh kích thước tùy theo asset của bạn
		Height: 16,
	}
	gme.potions = append(gme.potions, pot)
}

func (gme *ArcheroGame) handleAutoAttack() {
	// Nếu đang di chuyển thì không bắn (đặc trưng của Archero)
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS) ||
		ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) ||
		ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyDown) ||
		ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		return
	}

	target := gme.findNearestEnemy()
	if target == nil || !gme.player.CanAttack() {
		return
	}

	px, py := gme.player.GetCenter()
	ex, ey := target.GetCenter()

	// Vector hướng
	dx := ex - px
	dy := ey - py
	len := math.Hypot(dx, dy)
	dx /= len
	dy /= len

	// Vector vuông góc
	perpX := -dy
	perpY := dx

	offset := 4.0

	if gme.player.HasSkill(game.Multishot) {
		ox := perpX * offset
		oy := perpY * offset

		// Bắn 2 viên song song
		gme.spawnOffsetProjectile(px, py, ex, ey, ox, oy)
		gme.spawnOffsetProjectile(px, py, ex, ey, -ox, -oy)

	} else {
		// Chỉ bắn 1 viên thường
		p := g.NewProjectile(
			gme.projectileImg,
			px-4, py-4,
			ex, ey,
			4.5,
			gme.player.AttackDamage,
		)
		gme.projectiles = append(gme.projectiles, p)
	}

	// Nếu có kỹ năng DiagonalArrow (Bắn chéo 3 tia)
	if gme.player.HasSkill(game.DiagonalArrow) {
		// 1. Tính góc hiện tại từ người chơi đến quái vật (Radian)
		angle := math.Atan2(ey-py, ex-px)

		// 2. Tính tọa độ mục tiêu giả định cho tia bên TRÁI (Lệch -30 độ)
		angleLeft := angle - (math.Pi / 6)     // Pi/6 tương đương 30 độ
		exLeft := px + math.Cos(angleLeft)*200 // 200 là tầm xa giả định để định hướng
		eyLeft := py + math.Sin(angleLeft)*200
		gme.spawnProjectile(exLeft, eyLeft) // Truyền đúng tên biến đã khai báo

		// 3. Tính tọa độ mục tiêu giả định cho tia bên PHẢI (Lệch +30 độ)
		angleRight := angle + (math.Pi / 6)
		exRight := px + math.Cos(angleRight)*200
		eyRight := py + math.Sin(angleRight)*200
		gme.spawnProjectile(exRight, eyRight)
	}

	// 3. Đánh dấu người chơi đã tấn công để tính cooldown (tốc độ đánh)
	gme.player.Attack()
}

func (gme *ArcheroGame) spawnOffsetProjectile(
	px, py, ex, ey, ox, oy float64,
) {
	p := g.NewProjectile(
		gme.projectileImg,
		px-4+ox,
		py-4+oy,
		ex+ox,
		ey+oy,
		4.5,
		gme.player.AttackDamage,
	)
	gme.projectiles = append(gme.projectiles, p)
}

func (gme *ArcheroGame) spawnProjectile(targetX, targetY float64) {
	px, py := gme.player.GetCenter()
	p := g.NewProjectile(gme.projectileImg, px-4, py-4, targetX, targetY, 4.5, gme.player.AttackDamage)

	// Nếu có kỹ năng xuyên thấu (Piercing), bạn có thể thiết lập ở đây
	if gme.player.HasSkill(game.PiercingShot) {
		p.IsPiercing = true // Bạn cần thêm biến này vào struct Projectile
	}

	gme.projectiles = append(gme.projectiles, p)
}

func (gme *ArcheroGame) updateProjectiles() {
	for _, p := range gme.projectiles {
		p.Update(gme.mapWidthPx, gme.mapHeightPx)
		if !p.Active {
			continue
		}
		for _, e := range gme.enemies {
			if !e.IsAlive() {
				continue
			}
			if p.CheckCollision(e.X, e.Y, e.Width, e.Height) {
				e.TakeDamage(p.Damage)
				p.Active = false
				break
			}
		}
	}
}

func (gme *ArcheroGame) cleanupEntities() {
	filteredProj := gme.projectiles[:0]
	for _, p := range gme.projectiles {
		if p.Active {
			filteredProj = append(filteredProj, p)
		}
	}
	gme.projectiles = filteredProj

	filteredEnemies := gme.enemies[:0]
	for _, e := range gme.enemies {
		if e.IsAlive() {
			filteredEnemies = append(filteredEnemies, e)
		}
	}
	gme.enemies = filteredEnemies
}

func (gme *ArcheroGame) spawnEnemiesIfNeeded() {
	// mỗi khi wave tăng EnemiesSpawned, thêm enemy mới
	for len(gme.enemies) < gme.wave.EnemiesSpawned {
		x, y := gme.wave.GetSpawnPosition(gme.player.X, gme.player.Y)
		// log.Printf("Spawned enemy tại: x=%.2f, y=%.2f", x, y)
		enemy := g.NewEnemy(gme.enemyImg, x, y, 30, 1.2, 5, 400)
		gme.enemies = append(gme.enemies, enemy)
	}
}

func (gme *ArcheroGame) findNearestEnemy() *g.Enemy {
	px, py := gme.player.GetCenter()
	var best *g.Enemy
	bestDist := math.MaxFloat64
	for _, e := range gme.enemies {
		if !e.IsAlive() {
			continue
		}
		dist := e.GetDistanceTo(px, py)
		if dist < bestDist {
			bestDist = dist
			best = e
		}
	}
	return best
}

func (gme *ArcheroGame) handleWaveComplete() {
	if gme.wave.EnemiesSpawned >= gme.wave.EnemiesPerWave && len(gme.enemies) == 0 {
		gme.wave.StartNextWave()
	}
}

func (gme *ArcheroGame) updateCamera() {
	px, py := gme.player.GetCenter()
	gme.camera.SetFollowTarget(px, py)
	gme.camera.Update()

	// Clamp camera để không lộ ra ngoài map
	gme.camera.X = clamp(gme.camera.X, 0, gme.mapWidthPx-screenWidth)
	gme.camera.Y = clamp(gme.camera.Y, 0, gme.mapHeightPx-screenHeight)
}

func (gme *ArcheroGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{80, 160, 200, 255})
	g.DrawTilemap(screen, gme.tilemap, gme.tilesetImg, gme.camera.X, gme.camera.Y)

	// 2. CHÈN VÀO ĐÂY: Nếu đang trong trạng thái chọn kỹ năng thì mới vẽ menu
	if gme.gameState == StateSkillSelect {
		// Giả sử bạn truyền vào 3 kỹ năng ngẫu nhiên
		game.DrawSkillMenu(screen, gme.currentSkillOptions)
	}

	for _, e := range gme.enemies {
		e.Draw(screen, gme.camera.X, gme.camera.Y)
	}

	// Vẽ bình máu
	for _, pot := range gme.potions {
		pot.draw(screen, gme.camera.X, gme.camera.Y)
	}

	gme.player.Draw(screen, gme.camera.X, gme.camera.Y)

	for _, p := range gme.projectiles {
		p.Draw(screen, gme.camera.X, gme.camera.Y)
	}

	gme.drawUI(screen)
	gme.drawTeleportGateHint(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	// ebitenutil.DebugPrintAt(screen, "Vui lòng tắt bộ gõ Tiếng Việt (chuyển sang E) để di chuyển mượt mà bằng WASD", 10, 500)
}

func (gme *ArcheroGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (gme *ArcheroGame) drawUI(screen *ebiten.Image) {
	// health bar
	barW := 200.0
	barH := 12.0
	x := 20.0
	y := 20.0
	ratio := gme.player.Health / gme.player.MaxHealth
	if ratio < 0 {
		ratio = 0
	}
	ebitenutil.DrawRect(screen, x, y, barW, barH, color.RGBA{120, 0, 0, 255})
	ebitenutil.DrawRect(screen, x, y, barW*ratio, barH, color.RGBA{0, 200, 80, 255})

	ebitenutil.DebugPrintAt(screen, "HP", int(x), int(y)-12)
	ebitenutil.DebugPrintAt(screen, "Wave: "+itoa(gme.wave.CurrentWave), int(x), int(y)+20)
	ebitenutil.DebugPrintAt(screen, "F5: Save | F9: Load | ESC: Quit", int(x), int(y)+36)
}

// Phát hiện cổng chuyển map và xử lý nhấn E
func (gme *ArcheroGame) handleTeleportGate() {
	gateTile := TeleportGateTileID
	playerX := int((gme.player.X + gme.player.Width/2) / float64(gme.tilemap.TileW))
	playerY := int((gme.player.Y + gme.player.Height/2) / float64(gme.tilemap.TileH))
	for _, layer := range gme.tilemap.Layers {
		if playerY >= 0 && playerY < layer.Height && playerX >= 0 && playerX < layer.Width {
			tileID := layer.Data[playerY*layer.Width+playerX]
			if tileID == gateTile {
				if inpututil.IsKeyJustPressed(ebiten.KeyE) {
					// Demo: chỉ hiện thông báo, có thể load map mới ở đây
					log.Println("Chuyển sang map mới!")
				}
			}
		}
	}
}

func (gme *ArcheroGame) handleSaveLoad() {
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		gme.saveData.PlayerX = gme.player.X
		gme.saveData.PlayerY = gme.player.Y
		gme.saveData.MaxHealth = gme.player.MaxHealth
		gme.saveData.AttackDamage = gme.player.AttackDamage
		gme.saveData.AttackSpeed = gme.player.AttackSpeed
		if err := systems.SaveGameData(gme.saveData); err != nil {
			log.Printf("save failed: %v", err)
		} else {
			log.Println("saved game")
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF9) {
		data, err := systems.LoadGameData()
		if err != nil {
			log.Printf("load failed: %v", err)
			return
		}
		gme.saveData = data
		gme.resetStateFromSave()
	}
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func itoa(v int) string {
	return fmt.Sprintf("%d", v)
}

// Vẽ gợi ý chuyển map nếu player đang đứng ở tile cổng
func (gme *ArcheroGame) drawTeleportGateHint(screen *ebiten.Image) {
	gateTile := TeleportGateTileID
	playerX := int((gme.player.X + gme.player.Width/2) / float64(gme.tilemap.TileW))
	playerY := int((gme.player.Y + gme.player.Height/2) / float64(gme.tilemap.TileH))

	for _, layer := range gme.tilemap.Layers {
		if playerY >= 0 && playerY < layer.Height && playerX >= 0 && playerX < layer.Width {
			tileID := layer.Data[playerY*layer.Width+playerX]
			if tileID == gateTile {
				msg := "Ấn E để chuyển bản đồ!"
				x := int(gme.player.X - gme.camera.X)
				y := int(gme.player.Y - gme.camera.Y - 24)
				ebitenutil.DebugPrintAt(screen, msg, x-8, y)
			}
		}
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pixcel Archero-like")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewArcheroGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
