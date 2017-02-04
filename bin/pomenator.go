package main

import "pomenator"
import (
	"flag"
	"fmt"
	"os"
)

var pomCfg = flag.String("config", "", "name of the main config")
var secring = flag.String("keyring", "", "secure keyring file")
var passwd = flag.String("pgp-passwd", "", "secure keyring password")
var keyID = flag.String("key-id", "", "pgp key id")
var repoUser = flag.String("repo-user", "", "sonatype username")
var repoPasswd = flag.String("repo-passwd", "", "sonatype passwd")
var cfgFn = flag.String("secrets", "./.secrets.json", "json file containing pgp keyring, id and password, and your sonatype id, passwd")

func main() {

	flag.Parse()

	if *pomCfg == "" {
		usage("missing main json config\n")
	}

	// if any of the required flags aren't set, try to load from a
	// config file.
	var cfg pomenator.SecretsConfig

	/* pgp password could conceivable be "" */
	if *secring == "" || *keyID == "" || *repoUser == "" || *repoPasswd == "" {
		var err error
		if cfg, err = pomenator.LoadSecretsConfig(*cfgFn); err != nil {
			usage("could not load config: %s (%v)\n", *cfgFn, err)
		}
	}

	overwriteIfSet(&cfg.SecretKeyringFn, *secring)
	overwriteIfSet(&cfg.KeyID, *keyID)
	overwriteIfSet(&cfg.SecretKeyPasswd, *passwd)
	overwriteIfSet(&cfg.RepoUser, *repoUser)
	overwriteIfSet(&cfg.RepoPasswd, *repoPasswd)

	if cfg.SecretKeyringFn == "" || cfg.KeyID == "" || cfg.RepoUser == "" || cfg.RepoPasswd == "" {
		usage("missing mandatory config\n")
	}

	if err := pomenator.GenerateAllArtifacts(*pomCfg, cfg); err != nil {
		usage("error: %v\n", err)
	}
}

func overwriteIfSet(val *string, possibleOverwrite string) {
	if possibleOverwrite != "" {
		*val = possibleOverwrite
	}
}

func usage(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	flag.Usage()
	os.Exit(1)
}
