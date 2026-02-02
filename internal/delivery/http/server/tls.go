package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"sellers-accounts-backend/config"

	"go.uber.org/zap"
)

func ensureTLSFiles(tlsCfg config.TLSConfig, log *zap.SugaredLogger) error {
	if tlsCfg.CertFile == "" || tlsCfg.KeyFile == "" {
		return errors.New("tls.certFile and tls.keyFile are required when TLS is enabled")
	}

	certExists := fileExists(tlsCfg.CertFile)
	keyExists := fileExists(tlsCfg.KeyFile)
	if certExists && keyExists {
		return nil
	}

	if !tlsCfg.SelfSigned.Enabled {
		return fmt.Errorf("tls enabled but cert/key missing; either place files at %s/%s or enable tls.selfSigned", tlsCfg.CertFile, tlsCfg.KeyFile)
	}

	return generateSelfSignedCert(tlsCfg, log)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func generateSelfSignedCert(tlsCfg config.TLSConfig, log *zap.SugaredLogger) error {
	if err := os.MkdirAll(filepath.Dir(tlsCfg.CertFile), 0o755); err != nil {
		return fmt.Errorf("create cert directory: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(tlsCfg.KeyFile), 0o755); err != nil {
		return fmt.Errorf("create key directory: %w", err)
	}

	commonName := tlsCfg.SelfSigned.CommonName
	if commonName == "" {
		commonName = "localhost"
	}

	validFor := tlsCfg.SelfSigned.ValidForDays
	if validFor <= 0 {
		validFor = 365
	}

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("generate key: %w", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("generate serial: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(validFor) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range tlsCfg.SelfSigned.Hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}
	if len(template.DNSNames) == 0 && len(template.IPAddresses) == 0 {
		template.DNSNames = append(template.DNSNames, commonName)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("create certificate: %w", err)
	}

	certOut, err := os.Create(tlsCfg.CertFile)
	if err != nil {
		return fmt.Errorf("open cert file: %w", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		certOut.Close()
		return fmt.Errorf("write cert: %w", err)
	}
	if err := certOut.Close(); err != nil {
		return fmt.Errorf("close cert file: %w", err)
	}

	keyOut, err := os.Create(tlsCfg.KeyFile)
	if err != nil {
		return fmt.Errorf("open key file: %w", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		keyOut.Close()
		return fmt.Errorf("write key: %w", err)
	}
	if err := keyOut.Close(); err != nil {
		return fmt.Errorf("close key file: %w", err)
	}

	log.Infow("generated self-signed TLS certificate", "cert", tlsCfg.CertFile, "key", tlsCfg.KeyFile, "validForDays", validFor)

	return nil
}
