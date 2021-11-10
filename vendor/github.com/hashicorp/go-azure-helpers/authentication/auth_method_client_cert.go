package authentication

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/hashicorp/go-multierror"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"golang.org/x/crypto/pkcs12"
)

type servicePrincipalClientCertificateAuth struct {
	auxiliaryTenantIds []string
	clientId           string
	clientCertPath     string
	clientCertPassword string
	environment        string
	subscriptionId     string
	tenantId           string
	tenantOnly         bool
}

func (a servicePrincipalClientCertificateAuth) build(b Builder) (authMethod, error) {
	method := servicePrincipalClientCertificateAuth{
		clientId:           b.ClientID,
		clientCertPath:     b.ClientCertPath,
		clientCertPassword: b.ClientCertPassword,
		environment:        b.Environment,
		subscriptionId:     b.SubscriptionID,
		tenantId:           b.TenantID,
		tenantOnly:         b.TenantOnly,
		auxiliaryTenantIds: b.AuxiliaryTenantIDs,
	}
	return method, nil
}

func (a servicePrincipalClientCertificateAuth) isApplicable(b Builder) bool {
	return b.SupportsClientCertAuth && b.ClientCertPath != ""
}

func (a servicePrincipalClientCertificateAuth) name() string {
	return "Service Principal / Client Certificate"
}

func (a servicePrincipalClientCertificateAuth) getADALToken(_ context.Context, sender autorest.Sender, oauthConfig *OAuthConfig, endpoint string) (autorest.Authorizer, error) {
	if oauthConfig.OAuth == nil {
		return nil, fmt.Errorf("getting Authorization Token for client cert: an OAuth token wasn't configured correctly; please file a bug with more details")
	}

	// Get the certificate and private key from pfx file
	certificate, rsaPrivateKey, err := decodePkcs12File(a.clientCertPath, a.clientCertPassword)
	if err != nil {
		return nil, fmt.Errorf("decoding pkcs12 certificate: %v", err)
	}

	spt, err := adal.NewServicePrincipalTokenFromCertificate(*oauthConfig.OAuth, a.clientId, certificate, rsaPrivateKey, endpoint)
	if err != nil {
		return nil, err
	}

	spt.SetSender(sender)

	err = spt.Refresh()
	if err != nil {
		return nil, err
	}

	auth := autorest.NewBearerAuthorizer(spt)
	return auth, nil
}

func (a servicePrincipalClientCertificateAuth) getMSALToken(ctx context.Context, _ autorest.Sender, _ *OAuthConfig, endpoint string) (autorest.Authorizer, error) {
	certificate, rsaPrivateKey, err := decodePkcs12File(a.clientCertPath, a.clientCertPassword)
	if err != nil {
		return nil, fmt.Errorf("decoding pkcs12 certificate: %v", err)
	}

	environment, err := environments.EnvironmentFromString(a.environment)
	if err != nil {
		return nil, fmt.Errorf("environment config error: %v", err)
	}

	conf := auth.ClientCredentialsConfig{
		Environment:        environment,
		TenantID:           a.tenantId,
		AuxiliaryTenantIDs: a.auxiliaryTenantIds,
		ClientID:           a.clientId,
		PrivateKey:         x509.MarshalPKCS1PrivateKey(rsaPrivateKey),
		Certificate:        certificate.Raw,
		Scopes:             []string{fmt.Sprintf("%s/.default", strings.TrimRight(endpoint, "/"))},
		TokenVersion:       auth.TokenVersion2,
	}

	authorizer := conf.TokenSource(ctx, auth.ClientCredentialsAssertionType)
	if authTyped, ok := authorizer.(autorest.Authorizer); ok {
		return authTyped, nil
	}

	return nil, fmt.Errorf("returned auth.Authorizer does not implement autorest.Authorizer")
}

func (a servicePrincipalClientCertificateAuth) populateConfig(c *Config) error {
	c.AuthenticatedAsAServicePrincipal = true
	c.GetAuthenticatedObjectID = buildServicePrincipalObjectIDFunc(c)
	return nil
}

func (a servicePrincipalClientCertificateAuth) validate() error {
	var err *multierror.Error

	fmtErrorMessage := "a %s must be configured when authenticating as a Service Principal using a Client Certificate"

	if !a.tenantOnly && a.subscriptionId == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "Subscription ID"))
	}

	if a.clientId == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "Client ID"))
	}

	if a.clientCertPath == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "Client Certificate Path"))
	} else {
		// validate the certificate path is a valid pfx file
		_, _, derr := decodePkcs12File(a.clientCertPath, a.clientCertPassword)
		if derr != nil {
			err = multierror.Append(err, fmt.Errorf("the Client Certificate Path is not a valid pfx file: %v", derr))
		}
	}

	if a.tenantId == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "Tenant ID"))
	}

	return err.ErrorOrNil()
}

func decodePkcs12File(f string, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	certificateData, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, nil, fmt.Errorf("reading Client Certificate %q: %v", f, err)
	}

	privateKey, certificate, err := pkcs12.Decode(certificateData, password)
	if err != nil {
		return nil, nil, err
	}

	rsaPrivateKey, isRsaKey := privateKey.(*rsa.PrivateKey)
	if !isRsaKey {
		return nil, nil, fmt.Errorf("PKCS#12 certificate must contain an RSA private key")
	}

	return certificate, rsaPrivateKey, nil
}
