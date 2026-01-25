package mini3d

type Face3D struct {
	P1, P2, P3, P4 uint16 // indices into points array, P4 = 0xFFFF for triangles
	ColorBias      float32
}

func (f Face3D) IsTriangle() bool {
	return f.P4 == 0xFFFF
}

type Shape3D struct {
	Points    []Point3D
	Faces     []Face3D
	Center    Point3D
	ColorBias float32
	IsClosed  bool
}

func NewShape3D() *Shape3D {
	return &Shape3D{
		Points: make([]Point3D, 0),
		Faces:  make([]Face3D, 0),
	}
}

func (s *Shape3D) addPoint(p Point3D) uint16 {
	for i, existing := range s.Points {
		if existing.Equals(p) {
			return uint16(i)
		}
	}
	s.Points = append(s.Points, p)
	return uint16(len(s.Points) - 1)
}

func (s *Shape3D) AddFace(a, b, c Point3D) {
	s.AddFaceWithBias(a, b, c, 0)
}

func (s *Shape3D) AddFaceWithBias(a, b, c Point3D, colorBias float32) {
	face := Face3D{
		P1:        s.addPoint(a),
		P2:        s.addPoint(b),
		P3:        s.addPoint(c),
		P4:        0xFFFF, // triangle marker
		ColorBias: colorBias,
	}
	s.Faces = append(s.Faces, face)

	n := float32(len(s.Faces))
	s.Center.X += (a.X + b.X + c.X) / 3 / n
	s.Center.Y += (a.Y + b.Y + c.Y) / 3 / n
	s.Center.Z += (a.Z + b.Z + c.Z) / 3 / n
}

func (s *Shape3D) AddQuad(a, b, c, d Point3D) {
	s.AddQuadWithBias(a, b, c, d, 0)
}

func (s *Shape3D) AddQuadWithBias(a, b, c, d Point3D, colorBias float32) {
	face := Face3D{
		P1:        s.addPoint(a),
		P2:        s.addPoint(b),
		P3:        s.addPoint(c),
		P4:        s.addPoint(d),
		ColorBias: colorBias,
	}
	s.Faces = append(s.Faces, face)

	n := float32(len(s.Faces))
	s.Center.X += (a.X + b.X + c.X + d.X) / 4 / n
	s.Center.Y += (a.Y + b.Y + c.Y + d.Y) / 4 / n
	s.Center.Z += (a.Z + b.Z + c.Z + d.Z) / 4 / n
}

func (s *Shape3D) SetClosed(closed bool) {
	s.IsClosed = closed
}

func NewCube() *Shape3D {
	s := NewShape3D()

	// Front face
	s.AddQuad(
		NewPoint3D(-0.5, -0.5, 0.5),
		NewPoint3D(0.5, -0.5, 0.5),
		NewPoint3D(0.5, 0.5, 0.5),
		NewPoint3D(-0.5, 0.5, 0.5),
	)
	// Back face
	s.AddQuad(
		NewPoint3D(0.5, -0.5, -0.5),
		NewPoint3D(-0.5, -0.5, -0.5),
		NewPoint3D(-0.5, 0.5, -0.5),
		NewPoint3D(0.5, 0.5, -0.5),
	)
	// Top face
	s.AddQuad(
		NewPoint3D(-0.5, 0.5, 0.5),
		NewPoint3D(0.5, 0.5, 0.5),
		NewPoint3D(0.5, 0.5, -0.5),
		NewPoint3D(-0.5, 0.5, -0.5),
	)
	// Bottom face
	s.AddQuad(
		NewPoint3D(-0.5, -0.5, -0.5),
		NewPoint3D(0.5, -0.5, -0.5),
		NewPoint3D(0.5, -0.5, 0.5),
		NewPoint3D(-0.5, -0.5, 0.5),
	)
	// Right face
	s.AddQuad(
		NewPoint3D(0.5, -0.5, 0.5),
		NewPoint3D(0.5, -0.5, -0.5),
		NewPoint3D(0.5, 0.5, -0.5),
		NewPoint3D(0.5, 0.5, 0.5),
	)
	s.AddQuad(
		NewPoint3D(-0.5, -0.5, -0.5),
		NewPoint3D(-0.5, -0.5, 0.5),
		NewPoint3D(-0.5, 0.5, 0.5),
		NewPoint3D(-0.5, 0.5, -0.5),
	)

	s.SetClosed(true)
	return s
}

func NewIcosahedron() *Shape3D {
	s := NewShape3D()

	// Golden ratio
	p := float32((1.618033988749895 - 1) / 2) // (sqrt(5) - 1) / 2

	x1 := NewPoint3D(0, -p, 1)
	x2 := NewPoint3D(0, p, 1)
	x3 := NewPoint3D(0, p, -1)
	x4 := NewPoint3D(0, -p, -1)

	y1 := NewPoint3D(1, 0, p)
	y2 := NewPoint3D(1, 0, -p)
	y3 := NewPoint3D(-1, 0, -p)
	y4 := NewPoint3D(-1, 0, p)

	z1 := NewPoint3D(-p, 1, 0)
	z2 := NewPoint3D(p, 1, 0)
	z3 := NewPoint3D(p, -1, 0)
	z4 := NewPoint3D(-p, -1, 0)

	// Top cap
	s.AddFace(z1, y3, y4)
	s.AddFace(z1, x3, y3)
	s.AddFace(z1, z2, x3)
	s.AddFace(z1, x2, z2)
	s.AddFace(z1, y4, x2)

	// Middle band
	s.AddFace(y4, y3, z4)
	s.AddFace(z4, y3, x4)
	s.AddFace(y3, x3, x4)
	s.AddFace(x4, x3, y2)
	s.AddFace(x3, z2, y2)
	s.AddFace(y2, z2, y1)
	s.AddFace(z2, x2, y1)
	s.AddFace(y1, x2, x1)
	s.AddFace(x2, y4, x1)
	s.AddFace(x1, y4, z4)

	// Bottom cap
	s.AddFace(z3, y2, y1)
	s.AddFace(z3, y1, x1)
	s.AddFace(z3, x1, z4)
	s.AddFace(z3, z4, x4)
	s.AddFace(z3, x4, y2)

	s.SetClosed(true)
	return s
}
