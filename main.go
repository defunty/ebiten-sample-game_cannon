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

	cannonBallVelocity = 10
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
	fmt.Println("init")
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

	// CannonBall
	cbx  float64 // xPosition
	cby  float64 // yPosition
	cbv  float64 // vector速度
	cbvx float64 // x軸速度
	cbvy float64 // y軸速度

	// Cannon Degree
	cannonRadian float64
	cannonDegree int

	// sound effect
	// firePlayer   *audio.Player
	// hitPlayer   *audio.Player
	// criticalHitPlayer   *audio.Player
}

func (g *Game) setCannonDirection(x0, y0, x1, y1 int) {
	//var tan float64 = float64(y1-y0) / float64(x1-x0)
	//sita := math.Atan(tan)
	g.cannonRadian = math.Atan2(float64(y1-y0), float64(x1-x0)) + math.Pi/2
	g.cannonDegree = int(math.Round(g.cannonRadian * 180 / math.Pi)) // 0 ~ 360（中心の上位置を0として、時計回りで増加）
}

func (g *Game) drawCannon(screen *ebiten.Image) {
	w, h := cannonImage.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	// i := (g.count / 5) % frameNum
	//degree := calculateDegree(screenWidth/2, screenHeight/2, g.mx, g.my)

	op.GeoM.Rotate(float64(g.cannonDegree%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.Filter = ebiten.FilterLinear // シャギー対策
	screen.DrawImage(cannonImage, op)

	// 画像を回転させる際にspriteを採用した場合の処理。FilterLinearでspriteを利用する必要はなくなったが、参考用に残しておく
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
	//cbiW, cbiH := cannonBallImage.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.cbx, g.cby)
	// op.GeoM.Translate(float64(screenWidth/2), float64(screenHeight/2))
	screen.DrawImage(cannonBallImage, op)
}

func (g *Game) setInitialPositionCannonBall() {
	adjustX, adjustY := 2, 2
	g.cbx, g.cby = float64(screenWidth/2-adjustX), float64(screenHeight/2-adjustY)
}

func (g *Game) fireCannonBall() {
	fmt.Println("setInitialPositionCannonBall")
	g.setInitialPositionCannonBall()
	g.cbv = cannonBallVelocity
	g.cbvx = math.Sin(g.cannonRadian) * g.cbv // 中央から見て上の位置が0度地点なので、xはsinから算出する
	g.cbvy = math.Cos(g.cannonRadian) * g.cbv * -1
	fmt.Println(g.cannonDegree, g.cbvx, g.cbvy)
}

func (g *Game) updateCannonBall(g_mx int, g_my int) {
	if g.cbx < 0 || g.cby < 0 || g.cbx > screenWidth || g.cby > screenHeight {
		g.cbvx, g.cbvy = 0, 0
	} else {
		g.cbx += g.cbvx
		g.cby += g.cbvy
	}
	// fmt.Println(g.mx, g.my)
	g.setCannonDirection(screenWidth/2, screenHeight/2, g_mx, g_my)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.fireCannonBall()
	}
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

	g.updateCannonBall(g.mx, g.my)

	// if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
	// 	g.fireCannonBall()
	// }

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
	fmt.Println("main")
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Render Image")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
