package avatar

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"github.com/brickman-source/golang-utilities/compress"
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
)

type Gender int

const (
	MALE Gender = iota
	FEMALE
)

func init() {
}

func GenerateFromUsername(username string) (image.Image, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(username))
	if err != nil {
		return nil, err
	}
	return randomAvatar(int64(h.Sum32()))
}

func randomAvatar(seed int64) (image.Image, error) {
	rnd := rand.New(rand.NewSource(seed))
	avatar := image.NewRGBA(image.Rect(0, 0, 256, 256))
	var err error
	//20:background 30:skin 70:Hair 75:body 80:Mouth 90:FacialHair 100:Nose 110:Eyes
	//backgroundColor := backgroundColors[rnd.Intn(len(backgroundColors))]
	//background := color.RGBA{backgroundColor[0], backgroundColor[1], backgroundColor[2], 0xFF}
	r := backgroundColors[rnd.Intn(len(backgroundColors))]
	g := backgroundColors[rnd.Intn(len(backgroundColors))]
	for g == r {
		g = backgroundColors[rnd.Intn(len(backgroundColors))]
	}
	b := backgroundColors[rnd.Intn(len(backgroundColors))]
	for b == g || b == r {
		b = backgroundColors[rnd.Intn(len(backgroundColors))]
	}
	background := color.RGBA{backgroundColors[rnd.Intn(len(backgroundColors))], backgroundColors[rnd.Intn(len(backgroundColors))], backgroundColors[rnd.Intn(len(backgroundColors))], 0xFF}

	draw.Draw(avatar, avatar.Bounds(), &image.Uniform{background}, image.Point{X: 0, Y: 0}, draw.Src)

	err = drawImg(avatar, 6, rnd.Intn(categoryCount[6]), err) ///*2 Skin*/
	err = drawImg(avatar, 3, rnd.Intn(categoryCount[3]), err) ///*3 Hair*/
	err = drawImg(avatar, 0, rnd.Intn(categoryCount[0]), err) ///*4 Body*/
	err = drawImg(avatar, 4, rnd.Intn(categoryCount[4]), err) ///*5 Mouth*/
	err = drawImg(avatar, 2, rnd.Intn(categoryCount[2]), err) ///*6 FacialHair*/
	err = drawImg(avatar, 5, rnd.Intn(categoryCount[5]), err) ///*7 Nose*/
	err = drawImg(avatar, 1, rnd.Intn(categoryCount[1]), err) ///*8 Eyes*/
	return avatar, err
}

func drawImg(dst draw.Image, i, j int, err error) error {
	if err != nil {
		return err
	}

	data, err := hex.DecodeString(avatarBinData[i][j])
	if err != nil {
		return err
	}

	zipReader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	defer func() {
		_ = zipReader.Close()
	}()

	//data, err = compress.DecompressGzip(data)
	data, err = compress.Gunzip(data)
	if err != nil {
		return err
	}

	src, _, err := image.Decode(zipReader)
	if err != nil {
		return err
	}
	draw.Draw(dst, dst.Bounds(), src, image.Point{X: 0, Y: 0}, draw.Over)
	return nil
}
