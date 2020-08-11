package loadfiles

import (
	"os"
	"github.com/faiface/pixel"
	"image"
	_ "image/png"
	"github.com/kardianos/osext"
)

func LoadPicture(path string) (pixel.Picture, error, string) {
	path = osext.ExecutableFolder()+"/assets/"+path
	file, err := os.Open(path)
	if err != nil {
		return nil, err, osext.ExecutableFolder()
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err, osext.ExecutableFolder()
	}
	return pixel.PictureDataFromImage(img), nil, osext.ExecutableFolder()
}