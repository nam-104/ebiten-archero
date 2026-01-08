package systems

// Camera quản lý viewport
type Camera struct {
	X, Y     float64
	Width    float64
	Height   float64
	FollowX  float64
	FollowY  float64
}

// NewCamera tạo camera mới
func NewCamera(width, height float64) *Camera {
	return &Camera{
		X:       0,
		Y:       0,
		Width:   width,
		Height:  height,
		FollowX: 0,
		FollowY: 0,
	}
}

// Update cập nhật vị trí camera (smooth follow)
func (c *Camera) Update() {
	// Smooth camera follow
	targetX := c.FollowX - c.Width/2
	targetY := c.FollowY - c.Height/2

	// Linear interpolation cho smooth movement
	c.X += (targetX - c.X) * 0.1
	c.Y += (targetY - c.Y) * 0.1
}

// SetFollowTarget đặt target để camera follow
func (c *Camera) SetFollowTarget(x, y float64) {
	c.FollowX = x
	c.FollowY = y
}
