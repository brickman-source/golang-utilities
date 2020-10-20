package banner

import (
	"bytes"
	"fmt"
	"github.com/brickman-source/golang-utilities/banner/figure"
	_ "github.com/brickman-source/golang-utilities/banner/statik"
	"io"
	"math/rand"
	"os"

	"github.com/arsham/rainbow/rainbow"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
)

func Print(s string) {
	_ = Write(os.Stdout, s, "Calvin S.flf")
}

func Write(out io.Writer, msg, fontName string) error {
	fileS, err := fs.New()
	if err != nil {
		return err
	}
	font, err := fileS.Open(fmt.Sprintf("/%s", fontName))
	if err != nil {
		return errors.Wrap(err, fontName)
	}
	buf := &bytes.Buffer{}
	myFigure := figure.NewFigureWithFont(msg, font, true)
	figure.Write(buf, myFigure)
	l := &rainbow.Light{
		Writer: out,
		Seed:   rand.Int63n(256),
	}
	if _, err := io.Copy(l, buf); err != nil {
		return err
	}
	return nil
}
