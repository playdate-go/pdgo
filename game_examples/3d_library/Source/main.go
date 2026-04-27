package main

import (
	"3d_demo/mini3d"

	"github.com/playdate-go/pdgo"
)

var (
	pd      *pdgo.PlaydateAPI
	scene   *mini3d.Scene3D
	root    *mini3d.Scene3DNode
	rot1    mini3d.Matrix3D
	rot2    mini3d.Matrix3D
	cameraZ float32 = -4
	cameraX float32 = 0
)

// initGame is called once when the game starts
func initGame() {
	pd.Display.SetRefreshRate(50)

	scene = mini3d.NewScene3D()
	scene.SetCameraOrigin(0, 0, cameraZ)
	scene.SetLight(mini3d.NewVector3D(0.2, 0.8, 0.4))

	root = scene.GetRootNode()

	// Create an icosahedron shape
	shape := mini3d.NewIcosahedron()

	// Add six copies of the shape to the scene
	// Right - normal
	n1 := root.AddChildNode()
	n1.AddShape(shape, 2, 0, 0)

	// Left - bright
	n2 := root.AddChildNode()
	n2.AddShape(shape, -2, 0, 0)
	n2.SetColorBias(0.8)

	// Top - wireframe only (front faces)
	n3 := root.AddChildNode()
	n3.AddShape(shape, 0, 2, 0)
	n3.SetRenderStyle(mini3d.RenderWireframe | mini3d.RenderWireframeWhite)

	// Bottom - wireframe only (all faces)
	n4 := root.AddChildNode()
	n4.AddShape(shape, 0, -2, 0)
	n4.SetRenderStyle(mini3d.RenderWireframe | mini3d.RenderWireframeBack | mini3d.RenderWireframeWhite)

	// Front - filled + wireframe
	n5 := root.AddChildNode()
	n5.AddShape(shape, 0, 0, 2)
	n5.SetColorBias(0.8)
	n5.SetRenderStyle(mini3d.RenderFilled | mini3d.RenderWireframe)

	// Back - dark + wireframe
	n6 := root.AddChildNode()
	n6.AddShape(shape, 0, 0, -2)
	n6.SetColorBias(-0.8)
	n6.SetRenderStyle(mini3d.RenderFilled | mini3d.RenderWireframe | mini3d.RenderWireframeWhite)

	// Rotation matrices
	rot1 = mini3d.NewRotationMatrix(5, 0, 0, 1)
	rot2 = mini3d.NewRotationMatrix(3, 0, 1, 0)

	pd.System.LogToConsole("3D Demo initialized!")
}

// update is called every frame
func update() int {
	// Handle input for camera movement
	current, _, _ := pd.System.GetButtonState()

	var dx, dz float32 = 0, 0

	if current&pdgo.ButtonUp != 0 {
		dz = 0.1
	}
	if current&pdgo.ButtonDown != 0 {
		dz = -0.1
	}
	if current&pdgo.ButtonLeft != 0 {
		dx = -0.1
	}
	if current&pdgo.ButtonRight != 0 {
		dx = 0.1
	}

	if dx != 0 || dz != 0 {
		cameraX += dx
		cameraZ += dz
		scene.SetCameraOrigin(-cameraX, 0, cameraZ)
	}

	// Handle crank for camera roll
	if !pd.System.IsCrankDocked() {
		angle := pd.System.GetCrankAngle()
		rad := float64(angle) * 3.14159265 / 180.0
		sinA := float32(sin(rad))
		cosA := float32(cos(rad))
		scene.SetCamera(
			mini3d.NewPoint3D(-cameraX, 0, cameraZ),
			mini3d.NewPoint3D(0, 0, 0),
			1.0,
			mini3d.NewVector3D(sinA, cosA, 0),
		)
	}

	// Rotate the scene
	root.AddTransform(rot2)
	root.AddTransform(rot1)

	// Clear screen to black
	pd.Graphics.Clear(pdgo.SolidBlack)

	// Get frame buffer and draw
	bitmap := pd.Graphics.GetFrame()
	if bitmap != nil {
		scene.Draw(bitmap, pdgo.LCDRowSize)
		pd.Graphics.MarkUpdatedRows(0, pdgo.LCDRows-1)
	}

	return 1
}

// Simple sin function (Taylor series approximation)
func sin(x float64) float64 {
	// Normalize to -π to π
	for x > 3.14159265 {
		x -= 2 * 3.14159265
	}
	for x < -3.14159265 {
		x += 2 * 3.14159265
	}
	// Taylor series
	x3 := x * x * x
	x5 := x3 * x * x
	x7 := x5 * x * x
	return x - x3/6 + x5/120 - x7/5040
}

// Simple cos function
func cos(x float64) float64 {
	return sin(x + 3.14159265/2)
}

func main() {}
