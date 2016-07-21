package main

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func dumpCert(label string, data []byte) {
	offset := 0
	for len(data) > 0 {
		der, rest := pem.Decode(data)
		if der == nil {
			log.Debugf("Decode fail on %s", label)
			return

		}
		c, err := x509.ParseCertificate(der.Bytes)
		if err != nil {
			log.Debugf("Parse error on %s %s", label, err)
			return
		}
		n := time.Now()
		valid := "VALID"
		if n.Before(c.NotBefore) {
			valid = "NOT YET VALID"
		} else if n.After(c.NotAfter) {
			valid = "EXPIRED"
		}
		ca := ""
		if c.IsCA {
			ca = " (Cert is a CA)"
		}
		sans := c.DNSNames
		for _, ip := range c.IPAddresses {
			sans = append(sans, ip.String())
		}
		fmt.Printf(`%s:%d
  Subject:%s
    CN:%s
    ID:%s
    O:%s
    OU:%s
    Serial:%s
    Valid between %s and %s (%s)
    SANs: %s
  Issuer
    CN:%s
    ID:%s
    O:%s
    OU:%s
`,
			label, offset,
			ca,
			c.Subject.CommonName,
			hex.EncodeToString(c.SubjectKeyId),
			c.Subject.Organization,
			c.Subject.OrganizationalUnit,
			hex.EncodeToString(c.SerialNumber.Bytes()),
			c.NotBefore.String(), c.NotAfter.String(), valid,
			sans,
			c.Issuer.CommonName,
			hex.EncodeToString(c.AuthorityKeyId),
			c.Issuer.Organization,
			c.Issuer.OrganizationalUnit,
		)
		offset++
		data = rest
	}
	// Could dump much more...  https://golang.org/pkg/crypto/x509/#Certificate
}

func scan() {
	pem_candidates, _ := filepath.Glob("/*/*.pem")
	crt_candidates, _ := filepath.Glob("/*/*.crt")
	key_candidates, _ := filepath.Glob("/*/*.key")
	candidates := append(pem_candidates, crt_candidates...)
	candidates = append(candidates, key_candidates...)
	log.Infof("Candidates: %v", candidates)
	for _, filename := range candidates {
		log.Infof("Loading: %s", filename)
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			continue
		}
		dumpCert(filename, data)
	}
}

func main() {
	log.SetOutput(os.Stderr)
	app := cli.NewApp()
	app.Name = `Certificate Information Tool

By default the tool will read from stdin and dump the cert`
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "scan",
			Usage: "Scan /*/*.pem and display information for all readable certs found",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.Bool("scan") {
			scan()
		} else {
			data, _ := ioutil.ReadAll(os.Stdin)
			dumpCert("stdin", data)
		}
		return nil
	}
	app.Run(os.Args)
}
