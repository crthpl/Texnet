package loadfiles

import (
	"os"
	"github.com/faiface/pixel"
	"image"
	_ "image/png"
	"github.com/kardianos/osext"
	"github.com/golang/freetype/truetype"
)

func LoadPicture(path string) (pixel.Picture, error) {
	curPath, err :=osext.ExecutableFolder()
	path = curPath+"/assets/"+path
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}


func LoadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}