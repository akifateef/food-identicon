package main

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/gift"
)

var (
	zero = image.Point{0, 0}
)

func getPattern(i int) func(i int) image.Point {
	return func(i int) image.Point { return image.Point{(i % 3) * 100, (i / 3) * 100} }
}

func stitch(images []image.Image) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	for i, simg := range images {
		draw.Draw(img, simg.Bounds().Add(getPattern(len(images))(i)), simg, zero, draw.Src)
	}
	return img
}

func loadImages(fileNames []string) []image.Image {
	var images []image.Image
	for _, s := range fileNames {
		f, _ := os.OpenFile(s, os.O_RDONLY, 0644)
		img, _ := jpeg.Decode(f)
		images = append(images, img)
	}
	return images
}

// exists returns whether the given file or directory exists or not
// from http://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-denoted-by-a-path-exists-in-golang
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getFileNames(ingredients []string) []string {
	var ingredientImages []string
	for _, ingredient := range ingredients {
		ingredientFolder := strings.Join(strings.Split(strings.TrimSpace(ingredient), " "), "-")
		if !exists(path.Join("resized", "ingredients", ingredientFolder)) {
			continue
		}
		fileList := []string{}
		err := filepath.Walk(path.Join("resized", "ingredients", ingredientFolder), func(path string, f os.FileInfo, err error) error {
			if strings.Contains(path, ".jpg") || strings.Contains(path, ".JPG") {
				fileList = append(fileList, path)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		if len(fileList) > 0 {
			ingredientImages = append(ingredientImages, fileList[rand.Intn(len(fileList))])
		}
	}
	dest := make([]string, len(ingredientImages))
	perm := rand.Perm(len(ingredientImages))
	for i, v := range perm {
		dest[v] = ingredientImages[i]
	}
	return dest
}

func main() {
	// resizeEverything()
	rand.Seed(time.Now().Unix())
	fileNames := getFileNames([]string{
		"italian sausage", "ground beef",
		"onion", "garlic", "tomato", "tomatoes",
		"tomato paste", "tomato sauce", "water",
		"fennel seeds", "salt", "black pepper",
		"parsley", "lasagna noodles",
	})
	images := loadImages(fileNames)
	img := stitch(images)
	g := gift.New(
		gift.Sepia(60),
		gift.Saturation(-20),
	)
	post := image.NewRGBA(img.Bounds())
	g.Draw(post, img)
	b := bytes.NewBuffer(nil)
	jpeg.Encode(b, post, nil)
	ioutil.WriteFile("./a.jpg", b.Bytes(), 0644)
}
