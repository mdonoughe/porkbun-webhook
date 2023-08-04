<p align="center">
  <img src="https://raw.githubusercontent.com/cert-manager/cert-manager/d53c0b9270f8cd90d908460d69502694e1838f5f/logo/logo-small.png" height="256" width="256" alt="cert-manager project logo" />
</p>

# Porkbun Webhook solver for Cert Manager

This is an unofficial webhook solver for [Cert Manager](https://cert-manager.io/) and [Porkbun](https://porkbun.com/).

## Usage

1. Deploy the webhook:

    ```
    helm install porkbun-webhook ./deploy/porkbun-webhook \
        --set groupName=<your group>
    ```

2. Create a secret containing your [API key](https://porkbun.com/account/api):

    ```
    kubectl create secret generic porkbun-key \
        --from-literal=api-key=<your key> \
        --from-literal=secret-key=<your key>
    ```

3. Create a role and role binding:

    ```
    kubectl apply -f rbac.yaml
    ```

4. Configure a certificate issuer:

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

## Running the test suite

All DNS providers **must** run the DNS01 provider conformance testing suite,
else they will have undetermined behaviour when used with cert-manager.

**It is essential that you configure and run the test suite when creating a
DNS01 webhook.**

To run the tests, first put your api key into testdata/porkbun-solver/porkbun-secret.yaml.

Then you can run the test suite with:

```bash
TEST_ZONE_NAME=<your domain>. make test
```
