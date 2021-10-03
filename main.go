package main

import (
	"fmt" // for debug
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240

	frameOX     = 0  // cannon_spriteの関連値
	frameOY     = 0  // cannon_spriteの関連値
	frameWidth  = 26 // cannon_spriteの関連値
	frameHeight = 26 // cannon_spriteの関連値
	frameNum    = 11 // cannon_spriteの関連値
)

var (
	cannonImage     *ebiten.Image
	cannonBallImage *ebiten.Image
)

type mouseClick struct {
	originX, originY int
	currX, currY     int
	duration         int
}

// for debug
func OutputLog(s string) {
	fmt.Println(s)
	// log.Fatal(err)
}

func init() {
	var err error
	//cannonImage, _, err = ebitenutil.NewImageFromFile("src/dist/image/sprite_cannon.png")
	cannonImage, _, err = ebitenutil.NewImageFromFile("src/dist/image/cannon.png")
	cannonBallImage, _, err = ebitenutil.NewImageFromFile("src/dist/image/cannonball.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	count int

	// mousePosition
	mx, my int

	// CannonBall's position
	cbx int
	cby int

	// sound effect
	// firePlayer   *audio.Player
	// hitPlayer   *audio.Player
	// criticalHitPlayer   *audio.Player
}

func calculateDegree(x0, y0, x1, y1 int) int {
	//var tan float64 = float64(y1-y0) / float64(x1-x0)
	//sita := math.Atan(tan)
	radian := math.Atan2(float64(y1-y0), float64(x1-x0))
	var degree int = int(math.Round(radian*180/math.Pi)) + 90 // 0 ~ 360（中心の上位置を0として、時計回りで増加）

	return degree
}

func (g *Game) drawCannon(screen *ebiten.Image) {
	w, h := cannonImage.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	// i := (g.count / 5) % frameNum
	degree := calculateDegree(screenWidth/2, screenHeight/2, g.mx, g.my)
	fmt.Println(degree)

	op.GeoM.Rotate(float64(degree%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.Filter = ebiten.FilterLinear // シャギー対策
	screen.DrawImage(cannonImage, op)

	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	// op.GeoM.Translate(screenWidth/2, screenHeight/2)
	// // i := (g.count / 5) % frameNum
	// degree := calculateDegree(screenWidth/2, screenHeight/2, g.mx, g.my)
	// i := degree
	// sx, sy := frameOX+i*frameWidth, frameOY
	// screen.DrawImage(cannonImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) drawCannonBall(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	// i := (g.count / 5) % frameNum
	screen.DrawImage(cannonBallImage, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawCannon(screen)
	g.drawCannonBall(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// error型を返す
func (g *Game) Update() error {
	g.count++
	g.mx, g.my = ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		fmt.Println("hassya")

	}

	// What touches have just ended?
	// for id, t := range g.touches {
	// 	if inpututil.IsTouchJustReleased(id) {
	// 		if g.pinch != nil && (id == g.pinch.id1 || id == g.pinch.id2) {
	// 			g.pinch = nil
	// 		}
	// 		if g.pan != nil && id == g.pan.id {
	// 			g.pan = nil
	// 		}

	// 		// If this one has not been touched long (30 frames can be assumed
	// 		// to be 500ms), or moved far, then it's a tap.
	// 		diff := distance(t.originX, t.originY, t.currX, t.currY)
	// 		if !t.wasPinch && !t.isPan && (t.duration <= 30 || diff < 2) {
	// 			g.taps = append(g.taps, tap{
	// 				X: t.currX,
	// 				Y: t.currY,
	// 			})
	// 		}

	// 		delete(g.touches, id)
	// 	}
	// }
	//
	return nil
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Render Image")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
