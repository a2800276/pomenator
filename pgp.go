package pomenator

import "golang.org/x/crypto/openpgp"

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type SecretsConfig struct {
	SecretKeyringFn string `json:"secring"`
	KeyID           string `json:"keyId"`
	SecretKeyPasswd string `json:"passwd"`
	RepoUser        string `json:"repo_user"`
	RepoPasswd      string `json:"repo_passwd"`
}

func LoadSecretsConfig(fn string) (c SecretsConfig, err error) {
	file, err := os.Open(fn)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	return
}

func Sign(fn string, secCfg SecretsConfig) error {
	ascFn := fmt.Sprintf("%s.asc", fn)
	asc, err := os.OpenFile(ascFn, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer asc.Close()

	entity, err := loadEntity(secCfg)
	if err != nil {
		return err
	}

	jar, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer jar.Close()

	//func ArmoredDetachSign(w io.Writer, signer *Entity, message io.Reader, config *packet.Config) (err error)
	return openpgp.ArmoredDetachSign(asc, entity, jar, nil)
}

// try to retrieve and decrypt the key identified by the PGPConfig
func loadEntity(cfg SecretsConfig) (signer *openpgp.Entity, err error) {
	keyRing, err := os.Open(cfg.SecretKeyringFn)
	if err != nil {
		return
	}
	defer keyRing.Close()

	el, err := openpgp.ReadKeyRing(keyRing)
	if err != nil {
		// retry, maybe it's armored.
		keyRing, err = os.Open(cfg.SecretKeyringFn)
		if err != nil {
			return
		}
		defer keyRing.Close()

		el, err = openpgp.ReadArmoredKeyRing(keyRing)
		if err != nil {
			return
		}
	}

	// pgp display keyid is usually hex of the last 32 bits of the 64bit keyid.
	keyId, err := strconv.ParseUint(cfg.KeyID, 16, 32)
	if err != nil {
		return
	}

	for _, signer = range el {
		//fmt.Printf("%v\n", signer)
		//fmt.Printf("%v\n\n", signer.PrimaryKey.KeyId)
		if keyId == (signer.PrimaryKey.KeyId & 0xffffffff) {
			err = signer.PrivateKey.Decrypt([]byte(cfg.SecretKeyPasswd))
			break
		}
	}
	return
}
