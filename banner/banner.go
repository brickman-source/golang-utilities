package banner

import (
	"github.com/arsham/figurine/figurine"
	_ "github.com/brickman-source/golang-utilities/banner/statik"
	"os"
)

func Print(s string) {
	_ = figurine.Write(os.Stdout, s, "Calvin S.flf")
}
