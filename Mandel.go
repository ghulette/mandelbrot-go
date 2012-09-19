package main

import "os"
import "log"
import "math/cmplx"
import "image"
import "image/png"
import "image/color"

type pt struct {
  x, y, i int
}

func scaleX (n int) float64 {
  return float64(n) / float64(SIZE) * 2.0 - 1.5
}

func scaleY (n int) float64 {
  return float64(n) / float64(SIZE) * 2.0 - 1.0
}

func scaleColor (n int) float64 {
  return (float64(n) / float64(MAX_ITER) * 255.0)
}

const SIZE = 500
const MAX_ITER = 128

func render(r chan pt) {
  file, err := os.Create("test.png")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  m := image.NewRGBA(image.Rect(0, 0, SIZE-1, SIZE-1))
  for i := 0; i < SIZE * SIZE; i ++ {
    p := <- r
    b := scaleColor(p.i)
    c := color.RGBA{0, 0, 255-uint8(b), 255}
    m.Set(p.x, p.y, c)
  }
  png.Encode(file, m)
}

func mandel(c complex128) int {
  z := complex(0,0)
  i := 0
  for ; cmplx.Abs(z) < 2.0 && i < MAX_ITER; i++ {
    z = z * z + c
  }
  return i
}

func main() {
  r := make(chan pt, SIZE * SIZE)
  go func () {
    for y := 0 ; y < SIZE ; y ++ {
      for x := 0 ; x < SIZE ; x ++ {
        xs := scaleX(x)
        ys := scaleY(y)
        c := complex(xs, ys)
        go func (x,y int) {
          i := mandel(c)
          r <- pt{x,y,i}
        }(x,y)
      }
    }
  }()
  render(r)
}
