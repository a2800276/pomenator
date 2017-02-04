package pomenator

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

const ManifestFilename = "META-INF/MANIFEST.MF"
const ManifestDefault = `Manifest-Version: 1.0
Created-By: 0.0.1 (Pomenator)

`

// utility to create jar files.

func GenerateJarFromDirs(jarFn string, dirs ...string) error {
	jar, err := os.OpenFile(jarFn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer jar.Close()

	zw := zip.NewWriter(jar)
	foundManifest := false
	for _, dir := range dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info.Mode().IsRegular() {
				addFile(dir, path, zw)
			}
			return nil
		})
	}

	if !foundManifest {
		if err = addManifest(zw); err != nil {
			return err
		}
	}

	return zw.Close()
}

func addFile(baseDir string, fn string, jar *zip.Writer) error {
	entryName := fn[len(baseDir)+1:]

	zf, err := jar.Create(entryName)
	if err != nil {
		return err
	}

	//println(entryName)
	//println(fn)

	//open
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(zf, file)
	return err
}

func addManifest(zw *zip.Writer) error {
	zf, err := zw.Create(ManifestFilename)
	if err != nil {
		return err
	}
	if _, err = zf.Write([]byte(ManifestDefault)); err != nil {
		return err
	}
	return nil
}
