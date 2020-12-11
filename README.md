# microsoft-security-center-exporter
Prometheus exporter for various Microsoft Defender ATP metrics, taken from REST API

## Usage

1. Clone this repository
2. `go build ./`
3. Set required env var ([ref](https://github.com/Azure/azure-sdk-for-go#more-authentication-details))
    * Your application must be able to read security center score
```
1. **Client Credentials**: Azure AD Application ID and Secret.

    - `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
    - `AZURE_CLIENT_ID`: Specifies the app client ID to use.
    - `AZURE_CLIENT_SECRET`: Specifies the app secret to use.

2. **Client Certificate**: Azure AD Application ID and X.509 Certificate.

    - `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
    - `AZURE_CLIENT_ID`: Specifies the app client ID to use.
    - `AZURE_CERTIFICATE_PATH`: Specifies the certificate Path to use.
    - `AZURE_CERTIFICATE_PASSWORD`: Specifies the certificate password to use.

3. **Resource Owner Password**: Azure AD User and Password. This grant type is *not
   recommended*, use device login instead if you need interactive login.

    - `AZURE_TENANT_ID`: Specifies the Tenant to which to authenticate.
    - `AZURE_CLIENT_ID`: Specifies the app client ID to use.
    - `AZURE_USERNAME`: Specifies the username to use.
    - `AZURE_PASSWORD`: Specifies the password to use.
```
4. `./azure-security-center-exporter`
5. `curl localhost:8080/metrics`


## Metrics outputs

This exporter will read all available subscription and export the secure score percentage and current value.

```
# HELP azure_security_center_secure_score_percentage Azure Security Center Secure Score as percentage
# TYPE azure_security_center_secure_score_percentage gauge
azure_security_center_secure_score_percentage{subscription_id="Azure Subscription A"} 0.819
azure_security_center_secure_score_percentage{subscription_id="Azure Subscription B"} 0.7561
# HELP azure_security_center_secure_score_point Azure Security Center Secure Score as point
# TYPE azure_security_center_secure_score_point gauge
azure_security_center_secure_score_point{subscription_id="Azure Subscription A"} 47.5
azure_security_center_secure_score_point{subscription_id="Azure Subscription B"} 31
...
(various go instrumentation)
```