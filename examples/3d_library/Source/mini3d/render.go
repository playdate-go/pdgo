package mini3d

import "unsafe"

const (
	LCDWidth  = 400
	LCDHeight = 240
)

var patterns = [33][8]uint8{
	{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	{0x80, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00},
	{0x88, 0x00, 0x00, 0x00, 0x88, 0x00, 0x00, 0x00},
	{0x88, 0x00, 0x20, 0x00, 0x88, 0x00, 0x02, 0x00},
	{0x88, 0x00, 0x22, 0x00, 0x88, 0x00, 0x22, 0x00},
	{0xa8, 0x00, 0x22, 0x00, 0x8a, 0x00, 0x22, 0x00},
	{0xaa, 0x00, 0x22, 0x00, 0xaa, 0x00, 0x22, 0x00},
	{0xaa, 0x00, 0xa2, 0x00, 0xaa, 0x00, 0x2a, 0x00},
	{0xaa, 0x00, 0xaa, 0x00, 0xaa, 0x00, 0xaa, 0x00},
	{0xaa, 0x40, 0xaa, 0x00, 0xaa, 0x04, 0xaa, 0x00},
	{0xaa, 0x44, 0xaa, 0x00, 0xaa, 0x44, 0xaa, 0x00},
	{0xaa, 0x44, 0xaa, 0x10, 0xaa, 0x44, 0xaa, 0x01},
	{0xaa, 0x44, 0xaa, 0x11, 0xaa, 0x44, 0xaa, 0x11},
	{0xaa, 0x54, 0xaa, 0x11, 0xaa, 0x45, 0xaa, 0x11},
	{0xaa, 0x55, 0xaa, 0x11, 0xaa, 0x55, 0xaa, 0x11},
	{0xaa, 0x55, 0xaa, 0x51, 0xaa, 0x55, 0xaa, 0x15},
	{0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55},
	{0xba, 0x55, 0xaa, 0x55, 0xab, 0x55, 0xaa, 0x55},
	{0xbb, 0x55, 0xaa, 0x55, 0xbb, 0x55, 0xaa, 0x55},
	{0xbb, 0x55, 0xea, 0x55, 0xbb, 0x55, 0xae, 0x55},
	{0xbb, 0x55, 0xee, 0x55, 0xbb, 0x55, 0xee, 0x55},
	{0xfb, 0x55, 0xee, 0x55, 0xbf, 0x55, 0xee, 0x55},
	{0xff, 0x55, 0xee, 0x55, 0xff, 0x55, 0xee, 0x55},
	{0xff, 0x55, 0xfe, 0x55, 0xff, 0x55, 0xef, 0x55},
	{0xff, 0x55, 0xff, 0x55, 0xff, 0x55, 0xff, 0x55},
	{0xff, 0x55, 0xff, 0xd5, 0xff, 0x55, 0xff, 0x5d},
	{0xff, 0x55, 0xff, 0xdd, 0xff, 0x55, 0xff, 0xdd},
	{0xff, 0x75, 0xff, 0xdd, 0xff, 0x57, 0xff, 0xdd},
	{0xff, 0x77, 0xff, 0xdd, 0xff, 0x77, 0xff, 0xdd},
	{0xff, 0x77, 0xff, 0xfd, 0xff, 0x77, 0xff, 0xdf},
	{0xff, 0x77, 0xff, 0xff, 0xff, 0x77, 0xff, 0xff},
	{0xff, 0xf7, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff},
	{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minf(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func maxf(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func getRow32(bitmap []uint8, rowstride, y int) []uint32 {
	rowOffset := y * rowstride
	ptr := unsafe.Pointer(&bitmap[rowOffset])
	return unsafe.Slice((*uint32)(ptr), rowstride/4)
}

func swap(n uint32) uint32 {
	return ((n & 0xff000000) >> 24) | ((n & 0xff0000) >> 8) | ((n & 0xff00) << 8) | (n << 24)
}

func drawMaskPattern(row []uint32, idx int, mask, color uint32) {
	if idx >= len(row) {
		return
	}
	if mask == 0xffffffff {
		row[idx] = color
	} else {
		row[idx] = (row[idx] & ^mask) | (color & mask)
	}
}

func drawFragment(bitmap []uint8, rowstride int, y int, x1, x2 int, pattern [8]uint8) {
	if y < 0 || y >= LCDHeight {
		return
	}
	if x2 < 0 || x1 >= LCDWidth {
		return
	}
	if x1 < 0 {
		x1 = 0
	}
	if x2 > LCDWidth {
		x2 = LCDWidth
	}
	if x1 >= x2 {
		return
	}

	p := pattern[y%8]
	color := uint32(p)<<24 | uint32(p)<<16 | uint32(p)<<8 | uint32(p)

	row := getRow32(bitmap, rowstride, y)

	startbit := x1 % 32
	startmask := swap((1 << (32 - startbit)) - 1)
	endbit := x2 % 32
	var endmask uint32
	if endbit > 0 {
		endmask = swap(((1 << endbit) - 1) << (32 - endbit))
	}

	col := x1 / 32
	endcol := x2 / 32

	if col == endcol {
		var mask uint32
		if startbit > 0 && endbit > 0 {
			mask = startmask & endmask
		} else if startbit > 0 {
			mask = startmask
		} else if endbit > 0 {
			mask = endmask
		} else {
			mask = 0xffffffff
		}
		drawMaskPattern(row, col, mask, color)
	} else {
		x := x1

		if startbit > 0 {
			drawMaskPattern(row, col, startmask, color)
			col++
			x += 32 - startbit
		}

		for x+32 <= x2 {
			drawMaskPattern(row, col, 0xffffffff, color)
			col++
			x += 32
		}

		if endbit > 0 && col < len(row) {
			drawMaskPattern(row, col, endmask, color)
		}
	}
}

func slope(x1, y1, x2, y2 float32) int32 {
	dx := x2 - x1
	dy := y2 - y1

	if dy < 1 {
		return int32(dx * (1 << 16))
	}
	return int32(dx / dy * (1 << 16))
}

func sortTri(p1, p2, p3 *Point3D) (*Point3D, *Point3D, *Point3D) {
	y1, y2, y3 := p1.Y, p2.Y, p3.Y

	if y1 <= y2 && y1 <= y3 {
		if y2 <= y3 {
			return p1, p2, p3
		}
		return p1, p3, p2
	} else if y2 <= y1 && y2 <= y3 {
		if y1 <= y3 {
			return p2, p1, p3
		}
		return p2, p3, p1
	} else {
		if y1 <= y2 {
			return p3, p1, p2
		}
		return p3, p2, p1
	}
}

func FillTriangle(bitmap []uint8, rowstride int, p1, p2, p3 *Point3D, pattern [8]uint8) {
	p1, p2, p3 = sortTri(p1, p2, p3)

	endy := min(LCDHeight, int(p3.Y))

	if int(p1.Y) >= LCDHeight || endy < 0 {
		return
	}

	x1 := int32(p1.X * (1 << 16))
	x2 := x1

	sb := slope(p1.X, p1.Y, p2.X, p2.Y)
	sc := slope(p1.X, p1.Y, p3.X, p3.Y)

	var dx1, dx2 int32
	if sb < sc {
		dx1, dx2 = sb, sc
	} else {
		dx1, dx2 = sc, sb
	}

	y := int(p1.Y)
	if y < 0 {
		x1 += int32(-y) * dx1
		x2 += int32(-y) * dx2
		y = 0
	}

	midY := min(LCDHeight, int(p2.Y))
	for y < midY {
		left := int(x1 >> 16)
		right := int(x2 >> 16)
		if left > right {
			left, right = right, left
		}
		drawFragment(bitmap, rowstride, y, left, right+1, pattern)
		x1 += dx1
		x2 += dx2
		y++
	}

	dx := slope(p2.X, p2.Y, p3.X, p3.Y)

	if sb < sc {
		x1 = int32(p2.X * (1 << 16))
		dx1 = dx
	} else {
		x2 = int32(p2.X * (1 << 16))
		dx2 = dx
	}

	for y < endy {
		left := int(x1 >> 16)
		right := int(x2 >> 16)
		if left > right {
			left, right = right, left
		}
		drawFragment(bitmap, rowstride, y, left, right+1, pattern)
		x1 += dx1
		x2 += dx2
		y++
	}
}

func FillQuad(bitmap []uint8, rowstride int, p1, p2, p3, p4 *Point3D, pattern [8]uint8) {
	FillTriangle(bitmap, rowstride, p1, p2, p3, pattern)
	FillTriangle(bitmap, rowstride, p1, p3, p4, pattern)
}

func DrawLine(bitmap []uint8, rowstride int, p1, p2 *Point3D, thick int, pattern [8]uint8) {
	if p1.Y > p2.Y {
		p1, p2 = p2, p1
	}

	y := int(p1.Y)
	endy := int(p2.Y)

	if y >= LCDHeight || endy < 0 || minf(p1.X, p2.X) >= LCDWidth || maxf(p1.X, p2.X) < 0 {
		return
	}

	x := int32(p1.X * (1 << 16))
	dx := slope(p1.X, p1.Y, p2.X, p2.Y)
	py := p1.Y

	if y < 0 {
		x += int32(-p1.Y) * dx
		y = 0
		py = 0
	}

	x1 := x + dx*int32(float32(y)+1-py)

	for y <= endy {
		if y == endy {
			x1 = int32(p2.X * (1 << 16))
		}

		if dx < 0 {
			drawFragment(bitmap, rowstride, y, int(x1>>16), int(x>>16)+thick, pattern)
		} else {
			drawFragment(bitmap, rowstride, y, int(x>>16), int(x1>>16)+thick, pattern)
		}

		y++
		if y == LCDHeight {
			break
		}

		x = x1
		x1 += dx
	}
}
