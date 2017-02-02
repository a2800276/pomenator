package pomenator

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// utility to create jar files.

func GenerateJarFromDirs(fn string, dirs ...string)                              {}
func GenerateJarFromDirsWithManifest(fn string, manifest string, dirs ...string) {}

func GenerateJarFromFiles(jarFn string, baseDir string, fns ...string) error {
	baseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return err
	}

	jar, err := os.OpenFile(jarFn, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer jar.Close()

	w := zip.NewWriter(jar)

	for _, fn := range fns {
		println(fn)
		// check startswith baseDir
		fn, err := filepath.Abs(fn)
		if err != nil {
			return err
		}
		println(fn)
		if !strings.HasPrefix(fn, baseDir) {
			return fmt.Errorf("%s is not located in baseDir %s", fn, baseDir)
		}
		//open
		file, err := os.Open(fn)
		if err != nil {
			return err
		}
		defer file.Close()

		zf, err := w.Create(fn[len(baseDir)+1:])
		if err != nil {
			return err
		}

		_, err = io.Copy(zf, file)
		if err != nil {
			return err
		}

	}
	return w.Close()
}
