package main

import "pomenator"
import (
	"flag"
	"fmt"
	"os"
)

var pomCfg = flag.String("json", "", "name of the main config")
var secring = flag.String("keyring", "", "secure keyring file")
var passwd = flag.String("passwd", "", "secure keyring password")
var keyID = flag.String("key-id", "", "pgp key id")
var cfgFn = flag.String("pgp-config", "./.pgp_config.json", "json file containing pgp keyring, id and password")

func main() {

	flag.Parse()

	if *pomCfg == "" {
		usage("missing main json config")
	}

	// if any of the required flags aren't set, try to load from a
	// config file.
	var cfg pomenator.PGPConfig

	if *secring == "" || *keyID == "" {
		var err error
		if cfg, err = pomenator.LoadPGPConfig(*cfgFn); err != nil {
			usage("could not load config: %s (%v)\n", *cfgFn, err)
		}
	} else {
		cfg.SecretKeyringFn = *secring
		cfg.KeyID = *keyID
		cfg.SecretKeyPasswd = *passwd
	}

	if err := pomenator.GenerateAllArtifacts(*pomCfg, cfg); err != nil {
		usage("error: %v\n", err)
	}
}

func usage(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	flag.Usage()
	os.Exit(1)
}
