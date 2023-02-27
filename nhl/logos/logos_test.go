package logos

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadLogo_2(t *testing.T) {
	baseDir := "." // change this to the actual directory path
	files, err := filepath.Glob(filepath.Join(baseDir, "*.png"))
	if err != nil {
		fmt.Println(err)
	}

	images := make(map[string][]byte)
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
		}

		fileName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		images[fileName] = data
	}

	fmt.Println("Images:", len(images))
}
