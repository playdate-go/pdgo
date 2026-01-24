// Package mini3d provides a simple 3D rendering library for Playdate
package mini3d

import "math"

// Point3D represents a point in 3D space
type Point3D struct {
	X, Y, Z float32
}

// NewPoint3D creates a new Point3D
func NewPoint3D(x, y, z float32) Point3D {
	return Point3D{X: x, Y: y, Z: z}
}

// Equals checks if two points are equal
func (p Point3D) Equals(other Point3D) bool {
	return p.X == other.X && p.Y == other.Y && p.Z == other.Z
}

// Vector3D represents a direction/displacement in 3D space
type Vector3D struct {
	DX, DY, DZ float32
}

// NewVector3D creates a new Vector3D
func NewVector3D(dx, dy, dz float32) Vector3D {
	return Vector3D{DX: dx, DY: dy, DZ: dz}
}

// Cross computes the cross product of two vectors
func (v Vector3D) Cross(other Vector3D) Vector3D {
	return Vector3D{
		DX: v.DY*other.DZ - v.DZ*other.DY,
		DY: v.DZ*other.DX - v.DX*other.DZ,
		DZ: v.DX*other.DY - v.DY*other.DX,
	}
}

// Dot computes the dot product of two vectors
func (v Vector3D) Dot(other Vector3D) float32 {
	return v.DX*other.DX + v.DY*other.DY + v.DZ*other.DZ
}

// LengthSquared returns the squared length of the vector
func (v Vector3D) LengthSquared() float32 {
	return v.DX*v.DX + v.DY*v.DY + v.DZ*v.DZ
}

// Length returns the length of the vector
func (v Vector3D) Length() float32 {
	return float32(math.Sqrt(float64(v.LengthSquared())))
}

// Normalize returns a normalized version of the vector
func (v Vector3D) Normalize() Vector3D {
	d := v.Length()
	if d == 0 {
		return v
	}
	return Vector3D{DX: v.DX / d, DY: v.DY / d, DZ: v.DZ / d}
}

// AddVector adds a vector to a point
func (p Point3D) AddVector(v Vector3D) Point3D {
	return Point3D{X: p.X + v.DX, Y: p.Y + v.DY, Z: p.Z + v.DZ}
}

// PNormal computes the normal of a triangle defined by three points
func PNormal(p1, p2, p3 *Point3D) Vector3D {
	v := Vector3D{
		DX: p2.X - p1.X,
		DY: p2.Y - p1.Y,
		DZ: p2.Z - p1.Z,
	}.Cross(Vector3D{
		DX: p3.X - p1.X,
		DY: p3.Y - p1.Y,
		DZ: p3.Z - p1.Z,
	})
	return v.Normalize()
}

// Matrix3D represents a 3x3 transformation matrix with translation
type Matrix3D struct {
	IsIdentity bool
	Inverting  bool
	M          [3][3]float32
	DX, DY, DZ float32
}

// IdentityMatrix is the identity transformation
var IdentityMatrix = Matrix3D{
	IsIdentity: true,
	Inverting:  false,
	M:          [3][3]float32{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
	DX:         0, DY: 0, DZ: 0,
}

// NewMatrix3D creates a new transformation matrix
func NewMatrix3D(m11, m12, m13, m21, m22, m23, m31, m32, m33 float32, inverting bool) Matrix3D {
	return Matrix3D{
		IsIdentity: false,
		Inverting:  inverting,
		M:          [3][3]float32{{m11, m12, m13}, {m21, m22, m23}, {m31, m32, m33}},
		DX:         0, DY: 0, DZ: 0,
	}
}

// NewTranslateMatrix creates a translation matrix
func NewTranslateMatrix(dx, dy, dz float32) Matrix3D {
	return Matrix3D{
		IsIdentity: true,
		Inverting:  false,
		M:          [3][3]float32{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
		DX:         dx, DY: dy, DZ: dz,
	}
}

// NewRotationMatrix creates a rotation matrix around an arbitrary axis
func NewRotationMatrix(angleDeg float32, ax, ay, az float32) Matrix3D {
	angle := float64(angleDeg) * math.Pi / 180.0
	c := float32(math.Cos(angle))
	s := float32(math.Sin(angle))

	// Normalize axis
	l := float32(math.Sqrt(float64(ax*ax + ay*ay + az*az)))
	if l == 0 {
		return IdentityMatrix
	}
	ax, ay, az = ax/l, ay/l, az/l

	// Rodrigues' rotation formula
	return Matrix3D{
		IsIdentity: false,
		Inverting:  false,
		M: [3][3]float32{
			{c + ax*ax*(1-c), ax*ay*(1-c) - az*s, ax*az*(1-c) + ay*s},
			{ay*ax*(1-c) + az*s, c + ay*ay*(1-c), ay*az*(1-c) - ax*s},
			{az*ax*(1-c) - ay*s, az*ay*(1-c) + ax*s, c + az*az*(1-c)},
		},
		DX: 0, DY: 0, DZ: 0,
	}
}

// Multiply multiplies two matrices
func (l Matrix3D) Multiply(r Matrix3D) Matrix3D {
	m := Matrix3D{
		IsIdentity: false,
		Inverting:  l.Inverting != r.Inverting,
	}

	if l.IsIdentity {
		if r.IsIdentity {
			m = IdentityMatrix
		} else {
			m.M = r.M
		}
		m.DX = l.DX + r.DX
		m.DY = l.DY + r.DY
		m.DZ = l.DZ + r.DZ
	} else {
		if !r.IsIdentity {
			m.M[0][0] = l.M[0][0]*r.M[0][0] + l.M[1][0]*r.M[0][1] + l.M[2][0]*r.M[0][2]
			m.M[1][0] = l.M[0][0]*r.M[1][0] + l.M[1][0]*r.M[1][1] + l.M[2][0]*r.M[1][2]
			m.M[2][0] = l.M[0][0]*r.M[2][0] + l.M[1][0]*r.M[2][1] + l.M[2][0]*r.M[2][2]

			m.M[0][1] = l.M[0][1]*r.M[0][0] + l.M[1][1]*r.M[0][1] + l.M[2][1]*r.M[0][2]
			m.M[1][1] = l.M[0][1]*r.M[1][0] + l.M[1][1]*r.M[1][1] + l.M[2][1]*r.M[1][2]
			m.M[2][1] = l.M[0][1]*r.M[2][0] + l.M[1][1]*r.M[2][1] + l.M[2][1]*r.M[2][2]

			m.M[0][2] = l.M[0][2]*r.M[0][0] + l.M[1][2]*r.M[0][1] + l.M[2][2]*r.M[0][2]
			m.M[1][2] = l.M[0][2]*r.M[1][0] + l.M[1][2]*r.M[1][1] + l.M[2][2]*r.M[1][2]
			m.M[2][2] = l.M[0][2]*r.M[2][0] + l.M[1][2]*r.M[2][1] + l.M[2][2]*r.M[2][2]
		} else {
			m.M = l.M
		}

		m.DX = l.DX*r.M[0][0] + l.DY*r.M[0][1] + l.DZ*r.M[0][2] + r.DX
		m.DY = l.DX*r.M[1][0] + l.DY*r.M[1][1] + l.DZ*r.M[1][2] + r.DY
		m.DZ = l.DX*r.M[2][0] + l.DY*r.M[2][1] + l.DZ*r.M[2][2] + r.DZ
	}

	m.Inverting = l.Inverting != r.Inverting
	return m
}

// Apply applies the matrix transformation to a point
func (l Matrix3D) Apply(p Point3D) Point3D {
	if l.IsIdentity {
		return Point3D{X: p.X + l.DX, Y: p.Y + l.DY, Z: p.Z + l.DZ}
	}

	x := p.X*l.M[0][0] + p.Y*l.M[0][1] + p.Z*l.M[0][2] + l.DX
	y := p.X*l.M[1][0] + p.Y*l.M[1][1] + p.Z*l.M[1][2] + l.DY
	z := p.X*l.M[2][0] + p.Y*l.M[2][1] + p.Z*l.M[2][2] + l.DZ

	return Point3D{X: x, Y: y, Z: z}
}
