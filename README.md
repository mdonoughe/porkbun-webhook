[![Build](https://github.com/dmahmalat/cert-manager-porkbun-webhook/actions/workflows/publish.yml/badge.svg)](https://github.com/dmahmalat/cert-manager-porkbun-webhook/actions/workflows/publish.yml)

# ACME webhook for porkbun DNS API
Usage:
```bash
helm install my-release oci://ghcr.io/dmahmalat/charts/cert-manager-porkbun-webhook
```

To test:
```bash
TEST_DOMAIN_NAME=<domain> TEST_API_KEY=$(echo -n '<API Key>' | base64) TEST_SECRET_KEY=$(echo -n '<SECRET Key>' | base64) make test
```

# Example Issuer
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: <your e-mail>
    privateKeySecretRef:
      name: letsencrypt-key
    solvers:
      - selector:
          dnsZones:
            - <your domain>
        dns01:
          webhook:
            groupName: <your group>
            solverName: porkbun
            config:
              apiKeySecretRef:
                name: porkbun-key
                key: api-key
              secretKeySecretRef:
                name: porkbun-key
                key: secret-key
```

# Credits
This is based on the projects [mdonoughe/porkbun-webhook](https://github.com/mdonoughe/porkbun-webhook) and [cert-manager/webhook-example](https://github.com/cert-manager/webhook-example)
Additional credits to project [nblxa/cert-manager-webhook-google-domains](https://github.com/nblxa/cert-manager-webhook-google-domains) for various fixes, updates and automation.