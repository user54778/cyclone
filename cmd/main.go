package main

import (
	"fmt"

	"github.com/user54778/cyclone/internal/physics"
)

func main() {
	u := physics.NewVector3(0, 1, 1)
	v := physics.NewVector3(0, -1, 0)
	fmt.Println(u.Cross(v))
}
