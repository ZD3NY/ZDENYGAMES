package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	inpututil "github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenW = 560
	screenH = 480

	playerSpeed    = 2.8
	baseWolfSpeed  = 0.85
	attackRange    = 62.0
	halfAttackArc  = math.Pi / 2.4 // ~75 deg total sweep
	attackCooldown = 480 * time.Millisecond
	attackDuration = 180 * time.Millisecond
	damageCooldown = 900 * time.Millisecond

	playerR = 9.0
	wolfR   = 8.0
)

type vec2 struct{ x, y float64 }

func (v vec2) len() float64 { return math.Sqrt(v.x*v.x + v.y*v.y) }
func (v vec2) norm() vec2 {
	l := v.len()
	if l < 0.0001 {
		return vec2{1, 0}
	}
	return vec2{v.x / l, v.y / l}
}
func (v vec2) angle() float64 { return math.Atan2(v.y, v.x) }
func sub(a, b vec2) vec2      { return vec2{a.x - b.x, a.y - b.y} }
func dist(a, b vec2) float64  { return sub(a, b).len() }

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func angleDiff(a, b float64) float64 {
	d := a - b
	for d > math.Pi {
		d -= 2 * math.Pi
	}
	for d < -math.Pi {
		d += 2 * math.Pi
	}
	return d
}

type wolf struct {
	pos   vec2
	alive bool
}

type phase int

const (
	phaseMenu phase = iota
	phasePlay
	phaseWaveClear
	phaseOver
)

type gameState struct {
	phase phase

	playerPos    vec2
	playerFacing vec2
	playerHP     int

	attackActive bool
	attackTime   time.Time
	lastAttack   time.Time
	lastDamage   time.Time

	wolves      []wolf
	wave        int
	score       int
	waveClearAt time.Time

	rng            *rand.Rand
	scoreSubmitted bool
}

func newGame() *gameState {
	now := time.Now()
	return &gameState{
		phase:        phaseMenu,
		playerPos:    vec2{screenW / 2, screenH / 2},
		playerFacing: vec2{0, -1},
		playerHP:     3,
		lastAttack:   now.Add(-attackCooldown),
		lastDamage:   now.Add(-damageCooldown),
		rng:          rand.New(rand.NewSource(now.UnixNano())),
	}
}

func (g *gameState) spawnWave() {
	g.wave++
	count := 2 + g.wave*2
	g.wolves = make([]wolf, count)
	for i := range g.wolves {
		g.wolves[i] = wolf{pos: g.edgePos(), alive: true}
	}
}

func (g *gameState) edgePos() vec2 {
	switch g.rng.Intn(4) {
	case 0:
		return vec2{g.rng.Float64() * screenW, -20}
	case 1:
		return vec2{g.rng.Float64() * screenW, screenH + 20}
	case 2:
		return vec2{-20, g.rng.Float64() * screenH}
	default:
		return vec2{screenW + 20, g.rng.Float64() * screenH}
	}
}

func (g *gameState) wolfSpeed() float64 {
	return baseWolfSpeed + float64(g.wave-1)*0.07
}

func (g *gameState) aliveWolves() int {
	n := 0
	for _, w := range g.wolves {
		if w.alive {
			n++
		}
	}
	return n
}

func (g *gameState) doAttack() {
	fa := g.playerFacing.angle()
	for i := range g.wolves {
		if !g.wolves[i].alive {
			continue
		}
		d := sub(g.wolves[i].pos, g.playerPos)
		if d.len() > attackRange+wolfR {
			continue
		}
		if math.Abs(angleDiff(d.angle(), fa)) <= halfAttackArc {
			g.wolves[i].alive = false
			g.score += 10
		}
	}
}

type app struct {
	g *gameState
}

func (a *app) Update() error {
	g := a.g
	now := time.Now()

	switch g.phase {
	case phaseMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.spawnWave()
			g.phase = phasePlay
		}

	case phaseWaveClear:
		if now.Sub(g.waveClearAt) >= 2*time.Second {
			g.spawnWave()
			g.phase = phasePlay
		}

	case phaseOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			a.g = newGame()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			exitGame()
		}

	case phasePlay:
		// Move
		dx, dy := 0.0, 0.0
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			dy -= playerSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			dy += playerSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			dx -= playerSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			dx += playerSpeed
		}
		g.playerPos.x = clamp(g.playerPos.x+dx, playerR, screenW-playerR)
		g.playerPos.y = clamp(g.playerPos.y+dy, playerR, screenH-playerR)

		// Face mouse (cursor is already in logical coords in Ebiten)
		mx, my := ebiten.CursorPosition()
		d := sub(vec2{float64(mx), float64(my)}, g.playerPos)
		if d.len() > 5 {
			g.playerFacing = d.norm()
		}

		// Attack
		canAttack := now.Sub(g.lastAttack) >= attackCooldown
		if canAttack && !g.attackActive {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
				inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				g.attackActive = true
				g.attackTime = now
				g.lastAttack = now
				g.doAttack()
			}
		}
		if g.attackActive && now.Sub(g.attackTime) >= attackDuration {
			g.attackActive = false
		}

		// Move wolves toward player
		ws := g.wolfSpeed()
		for i := range g.wolves {
			if !g.wolves[i].alive {
				continue
			}
			dir := sub(g.playerPos, g.wolves[i].pos).norm()
			g.wolves[i].pos.x += dir.x * ws
			g.wolves[i].pos.y += dir.y * ws
		}

		// Wolf hits player
		if now.Sub(g.lastDamage) >= damageCooldown {
			for _, w := range g.wolves {
				if !w.alive {
					continue
				}
				if dist(w.pos, g.playerPos) < playerR+wolfR {
					g.playerHP--
					g.lastDamage = now
					if g.playerHP <= 0 {
						g.phase = phaseOver
						if !g.scoreSubmitted {
							g.scoreSubmitted = true
							submitScore(g.score, g.wave)
						}
					}
					break
				}
			}
		}

		// Wave cleared
		if g.phase == phasePlay && g.aliveWolves() == 0 {
			g.score += g.wave * 50
			g.waveClearAt = now
			g.phase = phaseWaveClear
		}
	}

	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	g := a.g

	// Forest floor
	screen.Fill(color.RGBA{16, 30, 14, 255})

	// Decorative trees (non-collidable)
	treeDark := color.RGBA{22, 50, 18, 255}
	treeMid := color.RGBA{18, 42, 14, 255}
	treeSpots := [][2]float64{
		{55, 55}, {505, 55}, {55, 425}, {505, 425},
		{280, 35}, {280, 445}, {35, 240}, {525, 240},
		{130, 140}, {430, 140}, {130, 340}, {430, 340},
		{200, 80}, {360, 80}, {200, 400}, {360, 400},
	}
	for _, t := range treeSpots {
		drawCircle(screen, t[0], t[1], 20, treeDark)
		drawCircle(screen, t[0], t[1], 13, treeMid)
		drawCircle(screen, t[0], t[1], 6, color.RGBA{14, 35, 10, 255})
	}

	// Wolves
	for _, w := range g.wolves {
		if !w.alive {
			continue
		}
		drawCircle(screen, w.pos.x, w.pos.y, wolfR, color.RGBA{185, 55, 35, 255})
		// glowing eyes
		drawCircle(screen, w.pos.x-3, w.pos.y-2, 2, color.RGBA{255, 210, 40, 255})
		drawCircle(screen, w.pos.x+3, w.pos.y-2, 2, color.RGBA{255, 210, 40, 255})
	}

	// Attack arc
	if g.attackActive {
		fa := g.playerFacing.angle()
		arcSteps := 18
		for i := 0; i <= arcSteps; i++ {
			t := float64(i) / float64(arcSteps)
			angle := fa - halfAttackArc + t*halfAttackArc*2
			ex := g.playerPos.x + math.Cos(angle)*attackRange
			ey := g.playerPos.y + math.Sin(angle)*attackRange
			ebitenutil.DrawLine(screen, g.playerPos.x, g.playerPos.y, ex, ey, color.RGBA{255, 220, 80, 160})
		}
	}

	// Player
	drawCircle(screen, g.playerPos.x, g.playerPos.y, playerR, color.RGBA{235, 175, 55, 255})
	// Facing dot
	fx := g.playerPos.x + g.playerFacing.x*14
	fy := g.playerPos.y + g.playerFacing.y*14
	drawCircle(screen, fx, fy, 2, color.RGBA{255, 240, 120, 255})

	// Damage flash: red tint on screen edges when recently hit
	now := time.Now()
	if now.Sub(g.lastDamage) < 300*time.Millisecond {
		ebitenutil.DrawRect(screen, 0, 0, screenW, screenH, color.RGBA{180, 0, 0, 60})
	}

	// HUD
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Wave %d", g.wave), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score), 10, 26)
	hp := ""
	for i := 0; i < g.playerHP; i++ {
		hp += "* "
	}
	for i := g.playerHP; i < 3; i++ {
		hp += ". "
	}
	ebitenutil.DebugPrintAt(screen, "HP: "+hp, 10, 42)
	if g.phase == phasePlay {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Wolves: %d", g.aliveWolves()), 10, 58)
	}

	// Controls hint (bottom)
	ebitenutil.DebugPrintAt(screen, "WASD: move  Mouse: aim  Click/Space: attack", 10, screenH-20)

	// Overlays
	switch g.phase {
	case phaseMenu:
		drawOverlay(screen,
			"WOLFPACK",
			"Survive the forest night",
			"Click or Space to begin",
			"")
	case phaseWaveClear:
		drawOverlay(screen,
			fmt.Sprintf("WAVE %d CLEARED", g.wave),
			fmt.Sprintf("Score: %d", g.score),
			"Next wave coming...",
			"")
	case phaseOver:
		drawOverlay(screen,
			"YOU FELL",
			fmt.Sprintf("Wave %d  |  Score: %d", g.wave, g.score),
			"R - Play Again",
			"Esc - Exit")
	}
}

func drawOverlay(screen *ebiten.Image, title, sub, line1, line2 string) {
	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	ebitenutil.DrawRect(screen, 0, 0, sw, sh, color.RGBA{0, 0, 0, 170})
	cx := int(sw / 2)
	cy := int(sh / 2)
	ebitenutil.DebugPrintAt(screen, title, cx-len(title)*3, cy-44)
	ebitenutil.DebugPrintAt(screen, sub, cx-len(sub)*3, cy-22)
	if line1 != "" {
		ebitenutil.DebugPrintAt(screen, line1, cx-len(line1)*3, cy+8)
	}
	if line2 != "" {
		ebitenutil.DebugPrintAt(screen, line2, cx-len(line2)*3, cy+28)
	}
}

func drawCircle(screen *ebiten.Image, cx, cy, r float64, col color.Color) {
	ri := int(r)
	r2 := ri * ri
	for dy := -ri; dy <= ri; dy++ {
		for dx := -ri; dx <= ri; dx++ {
			if dx*dx+dy*dy <= r2 {
				ebitenutil.DrawRect(screen, cx+float64(dx), cy+float64(dy), 1, 1, col)
			}
		}
	}
}

func (a *app) Layout(_, _ int) (int, int) {
	return screenW, screenH
}

func main() {
	ebiten.SetWindowTitle("Wolfpack")
	ebiten.SetWindowResizable(true)

	if m := ebiten.Monitor(); m != nil {
		mw, mh := m.Size()
		ebiten.SetWindowSize(mw/2, mh/2)
		ebiten.SetWindowPosition(mw/4, mh/4)
	}

	a := &app{g: newGame()}
	if err := ebiten.RunGame(a); err != nil {
		log.Fatal(err)
	}
}
