package pomenator

import (
	"fmt"
	"os"
)

func GenerateAllArtifacts(fn string, pgpCfg PGPConfig) error {
	// loadJsonConfig
	cfgs, err := LoadConfig(fn)
	if err != nil {
		return err
	}
	for _, cfg := range cfgs {
		if err := generateArtifacts(cfg, pgpCfg); err != nil {
			return err
		}
	}
	return nil
}

func generateArtifacts(cfg POMConfig, pgpCfg PGPConfig) (err error) {
	if err = makeDirIfNotExists(cfg.OutputDir); err != nil {
		return
	}
	// will generate all artifacts in:
	// {output}/{artifactId-version}
	// and then zip them all up into a jar:
	// {output}/bundle-{artifactId-version}.jar
	a_v := fmt.Sprintf("%s-%s", cfg.ArtifactID, cfg.Version)

	tempDir := fmt.Sprintf("%s/%s", cfg.OutputDir, a_v)
	if err = makeDirIfNotExists(tempDir); err != nil {
		return
	}

	pomFn := fmt.Sprintf("%s/%s.pom", tempDir, a_v)
	if err = GeneratePOM(pomFn, cfg); err != nil {
		return
	}
	if err = Sign(pomFn, pgpCfg); err != nil {
		return
	}

	srcFn := fmt.Sprintf("%s/%s-sources.jar", tempDir, a_v)
	if err = GenerateJarFromDirs(srcFn, cfg.SourceDirs...); err != nil {
		return
	}
	if err = Sign(srcFn, pgpCfg); err != nil {
		return
	}

	classesFn := fmt.Sprintf("%s/%s.jar", tempDir, a_v)
	if err = GenerateJarFromDirs(classesFn, cfg.ClassDirs...); err != nil {
		return
	}
	if err = Sign(classesFn, pgpCfg); err != nil {
		return
	}

	javadocDir := fmt.Sprintf("%s/javadoc", tempDir)
	if err = makeDirIfNotExists(javadocDir); err != nil {
		return
	}

	if err = GenerateJavadoc(cfg.SourceDirs, javadocDir); err != nil {
		return
	}
	javadocFn := fmt.Sprintf("%s/%s-javadoc.jar", tempDir, a_v)
	if err = GenerateJarFromDirs(javadocFn, javadocDir); err != nil {
		return
	}
	if err = Sign(javadocFn, pgpCfg); err != nil {
		return
	}

	bundleFn := fmt.Sprintf("%s/bundle-%s.jar", cfg.OutputDir, a_v)
	return GenerateJarFromDirs(bundleFn, tempDir)
}

func generatePOM() {}
func makeDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		println(dir)
		return os.Mkdir(dir, os.ModePerm)
	}
	return nil
}
