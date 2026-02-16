package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/velocity-trinity/core/pkg/logger"
)

var initCACmd = &cobra.Command{
	Use:   "init-ca",
	Short: "Initialize a Certificate Authority for mTLS",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log.Info("Generating Root CA...")

		// Generate CA Key
		caPriv, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			logger.Log.Fatal(err.Error())
		}

		// Generate CA Certificate
		caTemplate := x509.Certificate{
			SerialNumber: big.NewInt(2026),
			Subject: pkix.Name{
				Organization: []string{"Velocity Trinity CA"},
				CommonName:   "Velocity Trinity Root CA",
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().Add(365 * 24 * 10 * time.Hour), // 10 years
			IsCA:                  true,
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			BasicConstraintsValid: true,
		}

		caBytes, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caPriv.PublicKey, caPriv)
		if err != nil {
			logger.Log.Fatal(err.Error())
		}

		// Write Files
		pemToFile("ca.crt", "CERTIFICATE", caBytes)
		pemToFile("ca.key", "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(caPriv))
		
		fmt.Println("✅ Generated ca.crt and ca.key")
		fmt.Println("   Next: Generate server/client certs signed by this CA (Not implemented in this CLI yet)")
	},
}

func pemToFile(filename, typeName string, bytes []byte) {
	out, _ := os.Create(filename)
	defer out.Close()
	pem.Encode(out, &pem.Block{Type: typeName, Bytes: bytes})
}

func init() {
	rootCmd.AddCommand(initCACmd)
}
