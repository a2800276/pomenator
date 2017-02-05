package pomenator

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"text/template"
)

type Dependency struct {
	XMLName    struct{} `xml:"dependency"`
	GroupID    string   `json:"groupId" xml:"groupId"`
	ArtifactID string   `json:"artifactId" xml:"artifactId"`
	Version    string   `json:"version" xml:"version"`
}

// pom config
type POMConfig struct {
	GroupID        string       `json:"groupId"`
	ArtifactID     string       `json:"artifactId"`
	Version        string       `json:"version"`
	ProjectName    string       `json:"projectName"`
	Description    string       `json:"description"`
	URL            string       `json:"url"`
	LicenseName    string       `json:"licenseName"`
	LicenseURL     string       `json:"licenseURL"`
	ScmURL         string       `json:"scmURL"`
	DeveloperName  string       `json:"developerName"`
	DeveloperEmail string       `json:"developerEmail"`
	DeveloperURL   string       `json:"developerURL"`
	DeveloperID    string       `json:"developerId"`
	Dependencies   []Dependency `json:"dependencies"`
	SourceDirs     []string     `json:"sources"`
	ClassDirs      []string     `json:"classes"`
	OutputDir      string       `json:"output"`
}

var pomTemplate = `<project xmlns="http://maven.apache.org/POM/4.0.0"       
   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
   xsi:schemaLocation="http://maven.apache.org/POM/4.0.0  
   http://maven.apache.org/xsd/maven-4.0.0.xsd">          
   <modelVersion>4.0.0</modelVersion>

  <groupId>{{.GroupID}}</groupId>
  <artifactId>{{.ArtifactID}}</artifactId>
  <version>{{.Version}}</version>
 
  <name>{{.ProjectName}}</name>
  <description>{{.Description}}</description>
  <url>{{.URL}}</url>
  <licenses>
    <license>
      <name>{{.LicenseName}}</name>
      <url>{{.LicenseURL}}</url>
    </license>
  </licenses>
  <scm>
    <url>{{.ScmURL}}</url>
  </scm>
  <developers>
    <developer>
      <email>{{.DeveloperEmail}}</email>
      <name>{{.DeveloperName}}</name>
      <url>{{.DeveloperURL}}</url>
      <id>{{.DeveloperID}}</id>
    </developer>
  </developers>
  <dependencies>
%s
  </dependencies>
</project>`

func LoadConfig(fn string) ([]POMConfig, error) {
	cfg := []POMConfig{}
	file, err := os.Open(fn)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		file.Close()

		file, err = os.Open(fn)
		if err != nil {
			return cfg, err
		}
		defer file.Close()
		cfg2 := POMConfig{}
		decoder = json.NewDecoder(file)
		err = decoder.Decode(&cfg2)
		if err != nil {
			return cfg, err
		}
		cfg = append(cfg, cfg2)
	}
	return cfg, err
}

func GeneratePOM(fn string, cfg POMConfig) error {
	pom, err := os.Create(fn)
	if err != nil {
		return err
	}
	deps := genDependencies(cfg)
	templateWithDeps := fmt.Sprintf(pomTemplate, deps)
	tmplt := template.Must(template.New("pom").Parse(templateWithDeps))
	tmplt.Execute(pom, cfg)
	return nil
}

func genDependencies(cfg POMConfig) string {
	bs, err := xml.MarshalIndent(cfg.Dependencies, "    ", "  ")
	if err != nil {
		panic("can't happen, err marshaling xml")
	}
	return string(bs)
}
