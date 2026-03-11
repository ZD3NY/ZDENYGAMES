package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	inpututil "github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/bg_music.wav
var bgMusicWav []byte

//go:embed assets/drop.wav
var dropWav []byte

//go:embed assets/explosion.wav
var explosionWav []byte

//go:embed assets/game_over.wav
var gameOverWav []byte

const (
	boardW  = 10
	boardH  = 20
	cellPix = 24
	margin  = 12

	tickMS   = 500
	softDrop = 50
)

type cell int

var aMenu []menuItem

var shapes = [7][4][4][4]int{
	{
		{{0, 0, 0, 0}, {1, 1, 1, 1}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		{{0, 0, 1, 0}, {0, 0, 1, 0}, {0, 0, 1, 0}, {0, 0, 1, 0}},
		{{0, 0, 0, 0}, {1, 1, 1, 1}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		{{0, 0, 1, 0}, {0, 0, 1, 0}, {0, 0, 1, 0}, {0, 0, 1, 0}},
	},
	{
		{{0, 0, 0, 0}, {0, 1, 1, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
		{{0, 0, 0, 0}, {0, 1, 1, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
		{{0, 0, 0, 0}, {0, 1, 1, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
		{{0, 0, 0, 0}, {0, 1, 1, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
	},
	{
		{{0, 0, 0, 0}, {1, 1, 1, 0}, {0, 1, 0, 0}, {0, 0, 0, 0}},
		{{0, 1, 0, 0}, {1, 1, 0, 0}, {0, 1, 0, 0}, {0, 0, 0, 0}},
		{{0, 1, 0, 0}, {1, 1, 1, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		{{0, 1, 0, 0}, {0, 1, 1, 0}, {0, 1, 0, 0}, {0, 0, 0, 0}},
	},
	{
		{{0, 0, 0, 0}, {0, 1, 1, 0}, {1, 1, 0, 0}, {0, 0, 0, 0}},
		{{1, 0, 0, 0}, {1, 1, 0, 0}, {0, 1, 0, 0}, {0, 0, 0, 0}},
		{{0, 0, 0, 0}, {0, 1, 1, 0}, {1, 1, 0, 0}, {0, 0, 0, 0}},
		{{1, 0, 0, 0}, {1, 1, 0, 0}, {0, 1, 0, 0}, {0, 0, 0, 0}},
	},
	{
		{{0, 0, 0, 0}, {1, 1, 0, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
		{{0, 1, 0, 0}, {1, 1, 0, 0}, {1, 0, 0, 0}, {0, 0, 0, 0}},
		{{0, 0, 0, 0}, {1, 1, 0, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
		{{0, 1, 0, 0}, {1, 1, 0, 0}, {1, 0, 0, 0}, {0, 0, 0, 0}},
	},
	{
		{{0, 0, 0, 0}, {1, 1, 1, 0}, {0, 0, 1, 0}, {0, 0, 0, 0}},
		{{0, 1, 0, 0}, {0, 1, 0, 0}, {1, 1, 0, 0}, {0, 0, 0, 0}},
		{{1, 0, 0, 0}, {1, 1, 1, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		{{1, 1, 0, 0}, {1, 0, 0, 0}, {1, 0, 0, 0}, {0, 0, 0, 0}},
	},
	{
		{{0, 0, 0, 0}, {1, 1, 1, 0}, {1, 0, 0, 0}, {0, 0, 0, 0}},
		{{1, 1, 0, 0}, {0, 1, 0, 0}, {0, 1, 0, 0}, {0, 0, 0, 0}},
		{{0, 0, 1, 0}, {1, 1, 1, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		{{0, 1, 0, 0}, {0, 1, 0, 0}, {0, 1, 1, 0}, {0, 0, 0, 0}},
	},
}

var pieceColors = []color.RGBA{
	{255, 64, 64, 255},
	{64, 200, 64, 255},
	{255, 210, 64, 255},
	{64, 128, 255, 255},
	{180, 64, 180, 255},
	{64, 180, 180, 255},
	{255, 140, 64, 255},
}

type piece struct {
	t, rot int
	x, y   int
	color  color.RGBA
}

type game struct {
	board    [boardH][boardW]cell
	cur      piece
	next     piece
	score    int
	lines    int
	over     bool
	paused   bool
	showMenu bool

	scoreSubmitted bool

	gravityMS int
	lastFall  time.Time
	softAccum time.Duration
	rng       *rand.Rand

	moveDirX   int
	moveAccumX time.Duration

	audioCtx      *audio.Context
	bgPlayer      *audio.Player
	dropSound     *audio.Player
	clearSound    *audio.Player
	gameOverSound *audio.Player
}

type menuItem struct {
	label      string
	onActivate func()
}

type app struct {
	g       *game
	menuSel int
	lw, lh  int
	volume  float64
	muted   bool
}

func (g *game) projectDropY(p piece) int {
	q := p
	for g.canPlace(q, 0, 1, 0) {
		q.y++
	}
	return q.y
}

func (a *app) resetGame() {
	if a.g.bgPlayer != nil {
		a.g.bgPlayer.Pause()
		a.g.bgPlayer.Rewind()
	}
	if a.g.dropSound != nil {
		a.g.dropSound.Pause()
		a.g.dropSound.Rewind()
	}
	if a.g.clearSound != nil {
		a.g.clearSound.Pause()
		a.g.clearSound.Rewind()
	}

	rng := a.g.rng
	ctx := a.g.audioCtx

	*a.g = *newGame(ctx)
	a.g.rng = rng
	a.applyVolume()
}

func (a *app) applyVolume() {
	v := a.volume
	if a.muted {
		v = 0
	}
	if a.g.bgPlayer != nil {
		a.g.bgPlayer.SetVolume(v * 0.3)
	}
	if a.g.dropSound != nil {
		a.g.dropSound.SetVolume(v)
	}
	if a.g.clearSound != nil {
		a.g.clearSound.SetVolume(v)
	}
	if a.g.gameOverSound != nil {
		a.g.gameOverSound.SetVolume(v)
	}
}

const (
	initialDelay = 100 * time.Millisecond
	repeatDelay  = 50 * time.Millisecond
)

func pointInRect(px, py, rx, ry, w, h int) bool {
	return px >= rx && px < rx+w && py >= ry && py < ry+h
}

func (a *app) openMenu() {
	a.g.showMenu = true
	a.g.paused = true
	a.menuSel = 0

	items := []menuItem{
		{
			label: "Resume",
			onActivate: func() {
				a.g.showMenu = false
				a.setPaused(false)
			},
		},
		{
			label: "Restart",
			onActivate: func() {
				a.resetGame()
				a.g.showMenu = false
				a.setPaused(false)
			},
		},
		{
			label: "Fullscreen",
			onActivate: func() {
				ebiten.SetFullscreen(!ebiten.IsFullscreen())
			},
		},
		{
			label: "Exit to Site",
			onActivate: func() {
				exitGame()
			},
		},
	}
	aMenu = items
}

func (a *app) setPaused(paused bool) {
	a.g.paused = paused
	if !paused {
		a.g.lastFall = time.Now()
		if a.g.bgPlayer != nil {
			a.g.bgPlayer.Play()
		}
	} else {
		if a.g.bgPlayer != nil {
			a.g.bgPlayer.Pause()
		}
	}
}

func newGame(ctx *audio.Context) *game {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g := &game{
		gravityMS: tickMS,
		lastFall:  time.Now(),
		rng:       r,
		audioCtx:  ctx,
	}

	g.cur = g.randPiece()
	g.next = g.randPiece()

	g.bgPlayer = loadWav(ctx, bgMusicWav)
	g.bgPlayer.Play()

	g.dropSound = loadWav(ctx, dropWav)
	g.clearSound = loadWav(ctx, explosionWav)
	g.gameOverSound = loadWav(ctx, gameOverWav)

	return g
}

func drawRectOutline(screen *ebiten.Image, x, y, w, h, t int, col color.Color) {
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(t), col)
	ebitenutil.DrawRect(screen, float64(x), float64(y+h-t), float64(w), float64(t), col)
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(t), float64(h), col)
	ebitenutil.DrawRect(screen, float64(x+w-t), float64(y), float64(t), float64(h), col)
}

func (g *game) randPiece() piece {
	t := g.rng.Intn(7)
	return piece{
		t:     t,
		rot:   0,
		x:     boardW/2 - 2,
		y:     -1,
		color: pieceColors[t],
	}
}

func (g *game) canPlace(p piece, dx, dy, drot int) bool {
	rot := (p.rot + drot) & 3
	nx, ny := p.x+dx, p.y+dy
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if shapes[p.t][rot][r][c] == 0 {
				continue
			}
			x := nx + c
			y := ny + r
			if x < 0 || x >= boardW || y >= boardH {
				return false
			}
			if y >= 0 && g.board[y][x] != 0 {
				return false
			}
		}
	}
	return true
}

func (g *game) lockPiece() {
	p := g.cur
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if shapes[p.t][p.rot][r][c] == 0 {
				continue
			}
			x := p.x + c
			y := p.y + r
			if y >= 0 && y < boardH && x >= 0 && x < boardW {
				g.board[y][x] = 1
			}
		}
	}
	if g.dropSound != nil {
		g.dropSound.Rewind()
		g.dropSound.Play()
	}
	cleared := g.clearLines()
	switch cleared {
	case 1:
		g.score += 100
	case 2:
		g.score += 300
	case 3:
		g.score += 500
	case 4:
		g.score += 800
	}
	g.lines += cleared
	if g.lines > 0 && g.lines%10 == 0 && g.gravityMS > 120 {
		g.gravityMS -= 40
	}
	g.cur = g.next
	g.cur.x, g.cur.y = boardW/2-2, -1
	g.next = g.randPiece()
	if !g.canPlace(g.cur, 0, 0, 0) {
		g.over = true

		if g.bgPlayer != nil {
			g.bgPlayer.Pause()
			g.bgPlayer.Rewind()
		}

		if g.gameOverSound != nil {
			g.gameOverSound.Rewind()
			g.gameOverSound.Play()
		}

		if !g.scoreSubmitted {
			g.scoreSubmitted = true
			submitScore(g.score, g.lines)
		}
	}
}

func (g *game) clearLines() int {
	count := 0
	for y := boardH - 1; y >= 0; y-- {
		full := true
		for x := 0; x < boardW; x++ {
			if g.board[y][x] == 0 {
				full = false
				break
			}
		}
		if full {
			count++
			for yy := y; yy > 0; yy-- {
				g.board[yy] = g.board[yy-1]
			}
			for x := 0; x < boardW; x++ {
				g.board[0][x] = 0
			}
			y++
		}
	}
	if count > 0 && g.clearSound != nil {
		g.clearSound.Rewind()
		g.clearSound.Play()
	}
	return count
}

func (a *app) updateMenuInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		a.menuSel = (a.menuSel - 1 + len(aMenu)) % len(aMenu)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		a.menuSel = (a.menuSel + 1) % len(aMenu)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if len(aMenu) > 0 {
			aMenu[a.menuSel].onActivate()
		}
	}
}

func (a *app) Update() error {
	// M key to toggle mute (always active)
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		a.muted = !a.muted
		a.applyVolume()
	}

	// Volume slider + mute button mouse control (when menu is not open)
	if !a.g.showMenu && a.lw > 0 {
		mx, my := ebiten.CursorPosition()
		lx, ly := windowToLogical(mx, my, a.lw, a.lh)

		infoX := margin + boardW*cellPix + 20
		sliderX, sliderY, sliderW, sliderH := infoX, margin+316, 140, 8
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			if lx >= sliderX && lx <= sliderX+sliderW && ly >= sliderY-4 && ly <= sliderY+sliderH+4 {
				v := float64(lx-sliderX) / float64(sliderW)
				if v < 0 {
					v = 0
				}
				if v > 1 {
					v = 1
				}
				a.volume = v
				a.applyVolume()
			}
		}

		muteBtnY := margin + 334
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if lx >= infoX && lx < infoX+60 && ly >= muteBtnY && ly < muteBtnY+20 {
				a.muted = !a.muted
				a.applyVolume()
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		a.resetGame()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if a.g.paused {
			a.g.showMenu = false
			a.setPaused(false)
		} else {
			a.openMenu()
		}
		return nil
	}

	if a.g.over {
		offX, offY := margin, margin
		w := boardW * cellPix
		h := boardH * cellPix
		btnW, btnH := 140, 36
		btnX := offX + w/2 - btnW/2
		btnY := offY + h/2 + 30

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if pointInRect(x, y, btnX, btnY, btnW, btnH) {
				a.resetGame()
				return nil
			}
		}
		if a.g.showMenu {
			a.updateMenuInput()
			return nil
		}
		return nil
	}

	if a.g.paused {
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			a.menuSel = (a.menuSel - 1 + len(aMenu)) % len(aMenu)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			a.menuSel = (a.menuSel + 1) % len(aMenu)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			if len(aMenu) > 0 {
				aMenu[a.menuSel].onActivate()
			}
			return nil
		}

		mx, my := ebiten.CursorPosition()
		lx, ly := windowToLogical(mx, my, a.lw, a.lh)

		panelX, panelY, panelW, panelH := menuPanelRect(a.lw, a.lh)
		_ = panelH

		itemH, pad := 40, 12
		startY := panelY + pad + 40

		for i := range aMenu {
			iy := startY + i*(itemH+8)
			if pointInRect(lx, ly, panelX+pad, iy, panelW-2*pad, itemH) {
				a.menuSel = i
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					aMenu[i].onActivate()
					return nil
				}
			}
		}
		return nil
	}
	dir := 0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dir = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dir = 1
	}

	now := time.Now()
	if dir != 0 {
		if dir != a.g.moveDirX {
			if a.g.canPlace(a.g.cur, dir, 0, 0) {
				a.g.cur.x += dir
			}
			a.g.moveDirX = dir
			a.g.moveAccumX = 0
		} else {
			a.g.moveAccumX += 16 * time.Millisecond
			delay := initialDelay
			if a.g.moveAccumX > delay {
				for a.g.moveAccumX > repeatDelay {
					if a.g.canPlace(a.g.cur, dir, 0, 0) {
						a.g.cur.x += dir
					}
					a.g.moveAccumX -= repeatDelay
				}
			}
		}
	} else {
		a.g.moveDirX = 0
		a.g.moveAccumX = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if a.g.canPlace(a.g.cur, 0, 0, 1) {
			a.g.cur.rot = (a.g.cur.rot + 1) & 3
		} else if a.g.canPlace(a.g.cur, -1, 0, 1) {
			a.g.cur.x--
			a.g.cur.rot = (a.g.cur.rot + 1) & 3
		} else if a.g.canPlace(a.g.cur, 1, 0, 1) {
			a.g.cur.x++
			a.g.cur.rot = (a.g.cur.rot + 1) & 3
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		targetY := a.g.projectDropY(a.g.cur)
		dropped := targetY - a.g.cur.y
		if dropped > 0 {
			a.g.cur.y = targetY
			a.g.score += dropped * 2
			a.g.lockPiece()
			a.g.lastFall = time.Now()
		}
		return nil
	}

	if now.Sub(a.g.lastFall) >= time.Duration(a.g.gravityMS)*time.Millisecond {
		a.g.lastFall = now
		if a.g.canPlace(a.g.cur, 0, 1, 0) {
			a.g.cur.y++
		} else {
			a.g.lockPiece()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		a.g.softAccum += time.Duration(16) * time.Millisecond
		for a.g.softAccum >= time.Duration(softDrop)*time.Millisecond {
			a.g.softAccum -= time.Duration(softDrop) * time.Millisecond
			if a.g.canPlace(a.g.cur, 0, 1, 0) {
				a.g.cur.y++
				a.g.score++
			} else {
				a.g.lockPiece()
				break
			}
		}
	} else {
		a.g.softAccum = 0
	}

	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 24, 255})

	offX := margin
	offY := margin

	border := color.RGBA{200, 200, 220, 255}
	w := boardW * cellPix
	h := boardH * cellPix
	ebitenutil.DrawRect(screen, float64(offX-2), float64(offY-2), float64(w+4), float64(h+4), border)

	filled := color.RGBA{60, 60, 80, 255}
	for y := 0; y < boardH; y++ {
		for x := 0; x < boardW; x++ {
			if a.g.board[y][x] != 0 {
				ebitenutil.DrawRect(screen, float64(offX+x*cellPix), float64(offY+y*cellPix), float64(cellPix-1), float64(cellPix-1), filled)
			}
		}
	}

	p := a.g.cur
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if shapes[p.t][p.rot][r][c] == 0 {
				continue
			}
			x := p.x + c
			y := p.y + r
			if y >= 0 {
				ebitenutil.DrawRect(screen,
					float64(offX+x*cellPix),
					float64(offY+y*cellPix),
					float64(cellPix-1), float64(cellPix-1),
					p.color,
				)
			}
		}
	}

	infoX := offX + w + 20
	ebitenutil.DebugPrintAt(screen, "Go TETRIS", infoX, offY)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", a.g.score), infoX, offY+20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Lines: %d", a.g.lines), infoX, offY+40)
	ebitenutil.DebugPrintAt(screen, "Next:", infoX, offY+70)

	n := a.g.next
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if shapes[n.t][0][r][c] == 0 {
				continue
			}
			x := infoX + c*cellPix
			y := offY + 90 + r*cellPix
			ebitenutil.DrawRect(screen, float64(x), float64(y), float64(cellPix-1), float64(cellPix-1), n.color)
		}
	}

	ebitenutil.DebugPrintAt(screen, "Controls:", infoX, offY+200)
	ebitenutil.DebugPrintAt(screen, "←/→ move", infoX, offY+220)
	ebitenutil.DebugPrintAt(screen, "↓ soft drop", infoX, offY+240)
	ebitenutil.DebugPrintAt(screen, "↑ rotate", infoX, offY+260)
	ebitenutil.DebugPrintAt(screen, "Space hard drop", infoX, offY+280)

	// Volume slider
	ebitenutil.DebugPrintAt(screen, "Volume:", infoX, offY+300)
	sliderX, sliderY, sliderW, sliderH := infoX, offY+316, 140, 8
	ebitenutil.DrawRect(screen, float64(sliderX), float64(sliderY), float64(sliderW), float64(sliderH), color.RGBA{50, 50, 70, 255})
	fillW := int(float64(sliderW) * a.volume)
	fillCol := color.RGBA{90, 140, 220, 255}
	if a.muted {
		fillCol = color.RGBA{80, 80, 80, 255}
	}
	if fillW > 0 {
		ebitenutil.DrawRect(screen, float64(sliderX), float64(sliderY), float64(fillW), float64(sliderH), fillCol)
	}
	knobX := sliderX + int(float64(sliderW)*a.volume) - 4
	ebitenutil.DrawRect(screen, float64(knobX), float64(sliderY-3), 8, 14, color.RGBA{200, 220, 255, 255})

	// Mute button
	muteBtnY := offY + 334
	muteBtnCol := color.RGBA{60, 60, 80, 255}
	if a.muted {
		muteBtnCol = color.RGBA{160, 50, 50, 255}
	}
	ebitenutil.DrawRect(screen, float64(infoX), float64(muteBtnY), 60, 20, muteBtnCol)
	muteLabel := "Mute [M]"
	if a.muted {
		muteLabel = "Muted[M]"
	}
	ebitenutil.DebugPrintAt(screen, muteLabel, infoX+4, muteBtnY+5)

	ghostY := a.g.projectDropY(a.g.cur)

	ghostCol := color.RGBA{255, 64, 64, 200}
	thickness := 2

	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if shapes[a.g.cur.t][a.g.cur.rot][r][c] == 0 {
				continue
			}
			x := a.g.cur.x + c
			y := ghostY + r
			if y < 0 {
				continue
			}
			px := offX + x*cellPix
			py := offY + y*cellPix
			drawRectOutline(screen, px, py, cellPix-1, cellPix-1, thickness, ghostCol)
		}
	}

	if a.g.paused && a.g.showMenu {
		sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
		ebitenutil.DrawRect(screen, 0, 0, float64(sw), float64(sh), color.RGBA{0, 0, 0, 160})

		px, py, pw, ph := menuPanelRect(sw, sh)
		ebitenutil.DrawRect(screen, float64(px), float64(py), float64(pw), float64(ph), color.RGBA{30, 30, 40, 255})

		ebitenutil.DrawRect(screen, float64(px), float64(py), float64(pw), 2, color.RGBA{220, 220, 240, 255})
		ebitenutil.DrawRect(screen, float64(px), float64(py+ph-2), float64(pw), 2, color.RGBA{220, 220, 240, 255})
		ebitenutil.DrawRect(screen, float64(px), float64(py), 2, float64(ph), color.RGBA{220, 220, 240, 255})
		ebitenutil.DrawRect(screen, float64(px+pw-2), float64(py), 2, float64(ph), color.RGBA{220, 220, 240, 255})

		ebitenutil.DebugPrintAt(screen, "Pause Menu", px+16, py+14)

		itemH, pad := 40, 12
		startY := py + pad + 40
		for i, it := range aMenu {
			iy := startY + i*(itemH+8)
			rectCol := color.RGBA{60, 60, 80, 255}
			if i == a.menuSel {
				rectCol = color.RGBA{90, 120, 180, 255}
			}
			ebitenutil.DrawRect(screen, float64(px+pad), float64(iy), float64(pw-2*pad), float64(itemH), rectCol)
			ebitenutil.DebugPrintAt(screen, it.label, px+pad+12, iy+12)

			if it.label == "Fullscreen" {
				boxSize := 16
				boxX := px + pad + pw - 40
				boxY := iy + (itemH-boxSize)/2
				boxCol := color.RGBA{200, 200, 220, 255}
				ebitenutil.DrawRect(screen, float64(boxX), float64(boxY), float64(boxSize), float64(boxSize), boxCol)

				if ebiten.IsFullscreen() {
					checkCol := color.RGBA{64, 180, 64, 255}
					ebitenutil.DrawRect(screen, float64(boxX+4), float64(boxY+4), float64(boxSize-8), float64(boxSize-8), checkCol)
				}
			}
		}

		ebitenutil.DebugPrintAt(screen, "↑/↓ select • Enter/Space activate • Esc close", px+16, py+ph-28)
	}

	if a.g.over && !a.g.showMenu {
		overlay := color.RGBA{0, 0, 0, 160}
		ebitenutil.DrawRect(screen, float64(offX), float64(offY), float64(w), float64(h), overlay)

		ebitenutil.DebugPrintAt(screen, "GAME OVER", offX+w/2-40, offY+h/2-20)
		ebitenutil.DebugPrintAt(screen, "Press R to Restart", offX+w/2-70, offY+h/2)

		btnW, btnH := 140, 36
		btnX := offX + w/2 - btnW/2
		btnY := offY + h/2 + 30
		btnCol := color.RGBA{220, 80, 80, 255}
		ebitenutil.DrawRect(screen, float64(btnX), float64(btnY), float64(btnW), float64(btnH), btnCol)

		ebitenutil.DebugPrintAt(screen, "Restart", btnX+btnW/2-28, btnY+10)
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			a.resetGame()
		}
	}
}

func (a *app) Layout(outsideWidth, outsideHeight int) (int, int) {
	side := 200
	lw := margin*3 + boardW*cellPix + side
	lh := margin*2 + boardH*cellPix
	a.lw, a.lh = lw, lh
	return lw, lh
}

func menuPanelRect(lw, lh int) (x, y, w, h int) {
	w, h = 360, 270
	x = (lw - w) / 2
	y = (lh - h) / 2
	return
}

func windowToLogical(x, y, lw, lh int) (int, int) {
	ww, wh := ebiten.WindowSize()
	if ww == 0 || wh == 0 {
		return x, y
	}
	return x * lw / ww, y * lh / wh
}

func loadWav(ctx *audio.Context, data []byte) *audio.Player {
	d, err := wav.Decode(ctx, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	p, err := audio.NewPlayer(ctx, d)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func main() {
	ebiten.SetWindowTitle("Go TETRIS (Paciorek)")
	ebiten.SetWindowResizable(true)

	if m := ebiten.Monitor(); m != nil {
		mw, mh := m.Size()
		w, h := mw/2, mh/2
		ebiten.SetWindowDecorated(true)
		ebiten.SetFullscreen(false)
		ebiten.SetWindowSize(w, h)
		ebiten.SetWindowPosition(w/2, h/2)
	}

	audioCtx := audio.NewContext(44100)

	game := &app{
		g:      newGame(audioCtx),
		volume: 0.5,
	}
	game.applyVolume()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
