package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"testing"

	"github.com/cert-manager/cert-manager/test/acme/dns"
	"github.com/lost-woods/cert-manager-porkbun-webhook/porkbun"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var zapLogger, _ = zap.NewProduction()

var (
	domain    = os.Getenv("TEST_DOMAIN_NAME")
	apiKey    = os.Getenv("TEST_API_KEY")
	secretKey = os.Getenv("TEST_SECRET_KEY")

	configFile         = "_test/data/config.json"
	secretYamlFilePath = "_test/data/cert-manager-porkbun-webhook-secret.yaml"
	secretName         = "cert-manager-porkbun-webhook-secret"
	apiKeyRef          = "api-key"
	secretKeyRef       = "secret-key"
)

type SecretYaml struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind,omitempty" json:"kind,omitempty"`
	SecretType string `yaml:"type" json:"type"`
	Metadata   struct {
		Name string `yaml:"name"`
	}
	Data struct {
		ApiKey    string `yaml:"api-key"`
		SecretKey string `yaml:"secret-key"`
	}
}

func TestRunsSuite(t *testing.T) {
	slogger := zapLogger.Sugar()

	secretYaml := SecretYaml{}
	secretYaml.ApiVersion = "v1"
	secretYaml.Kind = "Secret"
	secretYaml.SecretType = "Opaque"
	secretYaml.Metadata.Name = secretName
	secretYaml.Data.ApiKey = apiKey
	secretYaml.Data.SecretKey = secretKey

	secretYamlFile, err := yaml.Marshal(&secretYaml)
	if err != nil {
		slogger.Errorf("Error: %v", err.Error())
	}
	_ = os.WriteFile(secretYamlFilePath, secretYamlFile, 0644)

	providerConfig := porkbun.PorkbunDNSProviderConfig{
		SecretNameRef:      secretName,
		ApiKeySecretRef:    apiKeyRef,
		SecretKeySecretRef: secretKeyRef,
	}
	file, _ := json.MarshalIndent(providerConfig, "", " ")
	_ = os.WriteFile(configFile, file, 0644)

	// resolvedFQDN must end with a '.'
	if domain[len(domain)-1:] != "." {
		domain = domain + "."
	}

	fixture := dns.NewFixture(&porkbun.PorkbunSolver{},
		dns.SetDNSName(domain),
		dns.SetResolvedZone(domain),
		dns.SetResolvedFQDN(GetRandomString(8)+"."+domain),
		dns.SetAllowAmbientCredentials(false),
		dns.SetManifestPath("_test/data"),
		dns.SetStrict(true),
	)

	fixture.RunConformance(t)

	_ = os.Remove(configFile)
	_ = os.Remove(secretYamlFilePath)
}

func GetRandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
