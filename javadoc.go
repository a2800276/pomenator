package pomenator

import (
	//"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var GenerateJavadocVerbosely = false

func GenerateJavadoc(sourceDirs []string, output string) error {
	args := []string{"-d", output}
	javaFiles := findJavaFiles(sourceDirs)
	args = append(args, javaFiles...)
	//for i, j := range args {
	//	fmt.Printf("%d : %s\n", i, j)
	//}
	javadocCmd := exec.Command("javadoc", args...)
	if GenerateJavadocVerbosely {
		javadocCmd.Stdout = os.Stdout
		javadocCmd.Stderr = os.Stderr
	}
	return javadocCmd.Run()

}

func findJavaFiles(sourceDirs []string) []string {
	ret := []string{}
	for _, dir := range sourceDirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			//println(path)
			if err == nil && info.Mode().IsRegular() && strings.HasSuffix(path, ".java") {
				//println(path)
				ret = append(ret, path)
			}
			return nil
		})
	}
	return ret
}
