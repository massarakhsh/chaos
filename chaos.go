package main

// TODO probably a bug in libui: changing the font away from skia leads to a crash

import (
	"github.com/massarakhsh/chaos/data"
	"github.com/massarakhsh/chaos/front"

	_ "github.com/andlabs/ui/winmanifest"
)

func main() {
	data.Generate()
	front.MainStart()
}
