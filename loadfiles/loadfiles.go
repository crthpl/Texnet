package loadfiles

import (
	"os"
	"github.com/faiface/pixel"
	"image"
	_ "image/png"
	"github.com/kardianos/osext"
)

func LoadPicture(path string) (pixel.Picture, error) {
	path, err = osext.ExecutableFolder()+"/assets/"+path
	fmt.Println(path)
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