package main

import (
	"os"
	"testing"

	"github.com/cert-manager/cert-manager/test/acme/dns"

	"github.com/mdonoughe/porkbun-webhook/porkbun"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRunsSuite(t *testing.T) {
	// The manifest path should contain a file named config.json that is a
	// snippet of valid configuration that should be included on the
	// ChallengeRequest passed as part of the test cases.

	fixture := dns.NewFixture(porkbun.New(),
		dns.SetStrict(true),
		dns.SetResolvedZone(zone),
		dns.SetAllowAmbientCredentials(false),
		dns.SetManifestPath("testdata/porkbun-solver"),
		dns.SetBinariesPath("_test/kubebuilder/bin"),
	)
	//need to uncomment and  RunConformance delete runBasic and runExtended once https://github.com/cert-manager/cert-manager/pull/4835 is merged
	//fixture.RunConformance(t)
	fixture.RunBasic(t)
	fixture.RunExtended(t)

}
