package main

import "pomenator"
import (
	"flag"
	"fmt"
	"os"
)

var input = flag.String("input", "", "name of the file to sign")
var secring = flag.String("keyring", "", "secure keyring file")
var passwd = flag.String("passwd", "", "secure keyring password")
var keyID = flag.String("key-id", "", "pgp key id")
var cfgFn = flag.String("pgp-config", "./.pgp_config.json", "json file containing pgp keyring, id and password")

func main() {

	flag.Parse()
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

	if *input == "" {
		usage("did not provide a file to sign")
	}

	if err := pomenator.Sign(*input, cfg); err != nil {
		usage("error signing %s with (cfg: %v) (%v)\n", *input, cfg, err)
	}
}

func usage(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	flag.Usage()
	os.Exit(1)
}
