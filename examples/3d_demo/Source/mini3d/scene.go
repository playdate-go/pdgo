package mini3d

import (
	"math"
	"sort"
)

type RenderStyle int

const (
	RenderInheritStyle   RenderStyle = 0
	RenderFilled         RenderStyle = 1 << 0
	RenderWireframe      RenderStyle = 1 << 1
	RenderWireframeBack  RenderStyle = 1 << 2
	RenderWireframeWhite RenderStyle = 1 << 3
)

type FaceInstance struct {
	P1, P2, P3, P4 *Point3D // pointers into transformed points
	Normal         Vector3D
	ColorBias      float32
}

type ShapeInstance struct {
	Prototype   *Shape3D
	Transform   Matrix3D
	Points      []Point3D
	Faces       []FaceInstance
	Center      Point3D
	ColorBias   float32
	RenderStyle RenderStyle
	Inverted    bool
}

type Scene3DNode struct {
	Transform   Matrix3D
	Parent      *Scene3DNode
	Children    []*Scene3DNode
	Shapes      []*ShapeInstance
	ColorBias   float32
	RenderStyle RenderStyle
	IsVisible   bool
	NeedsUpdate bool
}

func NewScene3DNode() *Scene3DNode {
	return &Scene3DNode{
		Transform:   IdentityMatrix,
		Children:    make([]*Scene3DNode, 0),
		Shapes:      make([]*ShapeInstance, 0),
		RenderStyle: RenderInheritStyle,
		IsVisible:   true,
		NeedsUpdate: true,
	}
}

func (n *Scene3DNode) SetTransform(xform Matrix3D) {
	n.Transform = xform
	node := n
	for node != nil {
		node.NeedsUpdate = true
		node = node.Parent
	}
}

func (n *Scene3DNode) AddTransform(xform Matrix3D) {
	m := n.Transform.Multiply(xform)
	n.SetTransform(m)
}

func (n *Scene3DNode) AddShape(shape *Shape3D, offsetX, offsetY, offsetZ float32) {
	n.AddShapeWithTransform(shape, NewTranslateMatrix(offsetX, offsetY, offsetZ))
}

func (n *Scene3DNode) AddShapeWithTransform(shape *Shape3D, transform Matrix3D) {
	instance := &ShapeInstance{
		Prototype:   shape,
		Transform:   transform,
		Points:      make([]Point3D, len(shape.Points)),
		Faces:       make([]FaceInstance, len(shape.Faces)),
		RenderStyle: RenderInheritStyle,
	}

	for i, face := range shape.Faces {
		instance.Faces[i] = FaceInstance{
			P1:        &instance.Points[face.P1],
			P2:        &instance.Points[face.P2],
			P3:        &instance.Points[face.P3],
			ColorBias: face.ColorBias,
		}
		if !face.IsTriangle() {
			instance.Faces[i].P4 = &instance.Points[face.P4]
		}
	}

	n.Shapes = append(n.Shapes, instance)
	n.NeedsUpdate = true
}

func (n *Scene3DNode) AddChildNode() *Scene3DNode {
	child := NewScene3DNode()
	child.Parent = n
	n.Children = append(n.Children, child)
	return child
}

func (n *Scene3DNode) SetColorBias(bias float32) {
	n.ColorBias = bias
	n.NeedsUpdate = true
}

func (n *Scene3DNode) SetRenderStyle(style RenderStyle) {
	n.RenderStyle = style
	n.NeedsUpdate = true
}

func (n *Scene3DNode) SetVisible(visible bool) {
	n.IsVisible = visible
	n.NeedsUpdate = true
}

type Scene3D struct {
	HasPerspective bool
	Camera         Matrix3D
	Light          Vector3D
	CenterX        float32
	CenterY        float32
	Scale          float32
	Root           *Scene3DNode
}

func NewScene3D() *Scene3D {
	scene := &Scene3D{
		HasPerspective: true,
		Root:           NewScene3DNode(),
		CenterX:        0.5,
		CenterY:        0.5,
	}
	scene.SetCamera(NewPoint3D(0, 0, 0), NewPoint3D(0, 0, 1), 1.0, NewVector3D(0, 1, 0))
	scene.SetLight(NewVector3D(0, -1, 0))
	return scene
}

func (s *Scene3D) GetRootNode() *Scene3DNode {
	return s.Root
}

func (s *Scene3D) SetCenter(x, y float32) {
	s.CenterX = x
	s.CenterY = y
}

func (s *Scene3D) SetLight(light Vector3D) {
	s.Light = light.Normalize()
}

func (s *Scene3D) SetCamera(origin, lookAt Point3D, scale float32, up Vector3D) {
	camera := IdentityMatrix
	camera.IsIdentity = false

	camera.DX = -origin.X
	camera.DY = -origin.Y
	camera.DZ = -origin.Z

	dir := NewVector3D(lookAt.X-origin.X, lookAt.Y-origin.Y, lookAt.Z-origin.Z)
	l := dir.Length()

	dir.DX /= l
	dir.DY /= l
	dir.DZ /= l

	s.Scale = 240 * scale

	h := float32(0)
	if dir.DX != 0 || dir.DZ != 0 {
		h = float32(math.Sqrt(float64(dir.DX*dir.DX + dir.DZ*dir.DZ)))
		yaw := NewMatrix3D(dir.DZ/h, 0, -dir.DX/h, 0, 1, 0, dir.DX/h, 0, dir.DZ/h, false)
		camera = camera.Multiply(yaw)
	}

	pitch := NewMatrix3D(1, 0, 0, 0, h, -dir.DY, 0, dir.DY, h, false)
	camera = camera.Multiply(pitch)

	if up.DX != 0 || up.DY != 0 {
		l = float32(math.Sqrt(float64(up.DX*up.DX + up.DY*up.DY)))
		roll := NewMatrix3D(up.DY/l, up.DX/l, 0, -up.DX/l, up.DY/l, 0, 0, 0, 1, false)
		s.Camera = camera.Multiply(roll)
	} else {
		s.Camera = camera
	}

	s.Root.NeedsUpdate = true
}

func (s *Scene3D) SetCameraOrigin(x, y, z float32) {
	s.SetCamera(NewPoint3D(x, y, z), NewPoint3D(0, 0, 0), 1.0, NewVector3D(0, 1, 0))
}

func (s *Scene3D) updateShapeInstance(shape *ShapeInstance, xform Matrix3D, colorBias float32, style RenderStyle) {
	proto := shape.Prototype

	for i, p := range proto.Points {
		shape.Points[i] = xform.Apply(shape.Transform.Apply(p))
	}

	shape.Center = xform.Apply(shape.Transform.Apply(proto.Center))
	shape.ColorBias = proto.ColorBias + colorBias
	shape.RenderStyle = style
	shape.Inverted = xform.Inverting

	for i := range shape.Faces {
		face := &shape.Faces[i]
		face.Normal = PNormal(face.P1, face.P2, face.P3)
	}

	for i := range shape.Points {
		p := &shape.Points[i]
		if p.Z > 0 {
			if s.HasPerspective {
				p.X = s.Scale * (p.X/p.Z + 1.6666666*s.CenterX)
				p.Y = s.Scale * (p.Y/p.Z + s.CenterY)
			} else {
				p.X = s.Scale * (p.X + 1.6666666*s.CenterX)
				p.Y = s.Scale * (p.Y + s.CenterY)
			}
		}
	}
}

func (s *Scene3D) updateNode(node *Scene3DNode, xform Matrix3D, colorBias float32, style RenderStyle, update bool) {
	if !node.IsVisible {
		return
	}

	if node.NeedsUpdate {
		update = true
		node.NeedsUpdate = false
	}

	if update {
		xform = node.Transform.Multiply(xform)
		colorBias += node.ColorBias

		if node.RenderStyle != RenderInheritStyle {
			style = node.RenderStyle
		}

		for _, shape := range node.Shapes {
			s.updateShapeInstance(shape, xform, colorBias, style)
		}

		for _, child := range node.Children {
			s.updateNode(child, xform, colorBias, style, update)
		}
	}
}

func (s *Scene3D) collectShapes(node *Scene3DNode, shapes []*ShapeInstance) []*ShapeInstance {
	if !node.IsVisible {
		return shapes
	}

	shapes = append(shapes, node.Shapes...)

	for _, child := range node.Children {
		shapes = s.collectShapes(child, shapes)
	}

	return shapes
}

func (s *Scene3D) drawShapeFace(shape *ShapeInstance, bitmap []uint8, rowstride int, face *FaceInstance) {
	if face.P1.Z <= 0 || face.P2.Z <= 0 || face.P3.Z <= 0 {
		return
	}
	if face.P4 != nil && face.P4.Z <= 0 {
		return
	}

	x1, y1 := face.P1.X, face.P1.Y
	x2, y2 := face.P2.X, face.P2.Y
	x3, y3 := face.P3.X, face.P3.Y

	var x4, y4 float32
	if face.P4 != nil {
		x4, y4 = face.P4.X, face.P4.Y
	}

	if (x1 < 0 && x2 < 0 && x3 < 0 && (face.P4 == nil || x4 < 0)) ||
		(x1 >= LCDWidth && x2 >= LCDWidth && x3 >= LCDWidth && (face.P4 == nil || x4 >= LCDWidth)) ||
		(y1 < 0 && y2 < 0 && y3 < 0 && (face.P4 == nil || y4 < 0)) ||
		(y1 >= LCDHeight && y2 >= LCDHeight && y3 >= LCDHeight && (face.P4 == nil || y4 >= LCDHeight)) {
		return
	}

	if shape.Prototype.IsClosed {
		var d float32
		if s.HasPerspective {
			d = (x2-x1)*(y3-y1) - (y2-y1)*(x3-x1)
		} else {
			d = face.Normal.DZ
		}

		inverted := shape.Inverted
		if (d >= 0) != inverted {
			return
		}
	}

	c := face.ColorBias + shape.ColorBias
	var v float32

	if c <= -1 {
		v = 0
	} else if c >= 1 {
		v = 1
	} else {
		if shape.Inverted {
			v = (1.0 + face.Normal.Dot(s.Light)) / 2
		} else {
			v = (1.0 - face.Normal.Dot(s.Light)) / 2
		}

		if c > 0 {
			v = c + (1-c)*v
		} else if c < 0 {
			v *= 1 + c
		}
	}

	vi := int(32.99 * v)
	if vi > 32 {
		vi = 32
	} else if vi < 0 {
		vi = 0
	}

	pattern := patterns[vi]

	if face.P4 != nil {
		FillQuad(bitmap, rowstride, face.P1, face.P2, face.P3, face.P4, pattern)
	} else {
		FillTriangle(bitmap, rowstride, face.P1, face.P2, face.P3, pattern)
	}
}

func (s *Scene3D) drawFilledShape(shape *ShapeInstance, bitmap []uint8, rowstride int) {
	for i := range shape.Faces {
		s.drawShapeFace(shape, bitmap, rowstride, &shape.Faces[i])
	}
}

func (s *Scene3D) drawWireframe(shape *ShapeInstance, bitmap []uint8, rowstride int) {
	style := shape.RenderStyle
	color := patterns[32] // white

	if style&RenderWireframeWhite == 0 {
		color = patterns[0] // black
	}

	for _, face := range shape.Faces {
		if face.P1.Z <= 0 || face.P2.Z <= 0 || face.P3.Z <= 0 {
			continue
		}
		if face.P4 != nil && face.P4.Z <= 0 {
			continue
		}

		x1, y1 := face.P1.X, face.P1.Y
		x2, y2 := face.P2.X, face.P2.Y
		x3, y3 := face.P3.X, face.P3.Y

		if style&RenderWireframeBack == 0 {
			var d float32
			if s.HasPerspective {
				d = (x2-x1)*(y3-y1) - (y2-y1)*(x3-x1)
			} else {
				d = face.Normal.DZ
			}

			if (d > 0) != shape.Inverted {
				continue
			}
		}

		DrawLine(bitmap, rowstride, face.P1, face.P2, 1, color)
		DrawLine(bitmap, rowstride, face.P2, face.P3, 1, color)

		if face.P4 != nil {
			DrawLine(bitmap, rowstride, face.P3, face.P4, 1, color)
			DrawLine(bitmap, rowstride, face.P4, face.P1, 1, color)
		} else {
			DrawLine(bitmap, rowstride, face.P3, face.P1, 1, color)
		}
	}
}

func (s *Scene3D) Draw(bitmap []uint8, rowstride int) {
	s.updateNode(s.Root, s.Camera, 0, RenderFilled, false)

	shapes := s.collectShapes(s.Root, nil)

	sort.Slice(shapes, func(i, j int) bool {
		return shapes[i].Center.Z > shapes[j].Center.Z
	})

	for _, shape := range shapes {
		style := shape.RenderStyle

		if style&RenderFilled != 0 {
			s.drawFilledShape(shape, bitmap, rowstride)
		}

		if style&RenderWireframe != 0 {
			s.drawWireframe(shape, bitmap, rowstride)
		}
	}
}
