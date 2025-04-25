package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
)

type config struct {
	port                                                       int
	createAdminRolePath                                        string
	inviteAdminRolePath                                        string
	adminRoles                                                 string
	credentialMapping                                          string
	credentialDataTemplate                                     string
	credentialEndpoint                                         string
	claimMappingServiceURL                                     string
	inviteIdentityProviderOidURL, createIdentityProviderOidURL string
	mailSupportAddress, mailSmtpHost, mailSmtpPort             string
	mailSmtpUsername, mailSmtpPassword                         string
	mailTemplate, mailTemplateKeys                             string
}

func getConfig() (config, error) {
	port, found := os.LookupEnv("PORT")
	if !found {
		err := fmt.Errorf("Environemnt variable \"PORT\" not found")
		return config{}, err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return config{}, err
	}
	Logger.Debug(portInt)
	adminRoles, found := os.LookupEnv("ADMIN_ROLES")
	if !found {
		err := fmt.Errorf("Environemnt variable \"ADMIN_ROLES\" not found")
		return config{}, err
	}
	Logger.Debug(adminRoles)
	createAdminRolePath, found := os.LookupEnv("CREATE_ADMIN_ROLE_PATH")
	if !found {
		err := fmt.Errorf("Environemnt variable \"CREATE_ADMIN_ROLE_PATH\" for Create Principal not found")
		return config{}, err
	}
	Logger.Debug(createAdminRolePath)
	inviteAdminRolePath, found := os.LookupEnv("INVITE_ADMIN_ROLE_PATH")
	if !found {
		err := fmt.Errorf("Environemnt variable \"INVITE_ADMIN_ROLE_PATH\" for Create Invitation not found")
		return config{}, err
	}
	Logger.Debug(inviteAdminRolePath)
	credentialMapping, found := os.LookupEnv("CREDENTIAL_MAPPING")
	if !found {
		err := fmt.Errorf("Environemnt variable \"CREDENTIAL_MAPPING\" not found")
		return config{}, err
	}
	Logger.Debug(credentialMapping)
	credentialDataTemplate, found := os.LookupEnv("CREDENTIAL_DATA_TEMPLATE")
	if !found {
		err := fmt.Errorf("Environemnt variable \"CREDENTIAL_DATA_TEMPLATE\" not found")
		return config{}, err
	}

	Logger.Debug(credentialDataTemplate)
	credentialEndpoint, found := os.LookupEnv("CREDENTIAL_ENDPOINT")
	if !found {
		err := fmt.Errorf("Environemnt variable \"CREDENTIAL_ENDPOINT\" not found")
		return config{}, err
	}
	Logger.Debug(credentialEndpoint)

	claimMappingServiceURL, found := os.LookupEnv("CLAIM_MAPPING_SERVICE_URL")
	if !found {
		err := fmt.Errorf("Environemnt variable \"CLAIM_MAPPING_SERVICE_URL\" not found")
		return config{}, err
	}

	Logger.Debug(claimMappingServiceURL)
	inviteIdentityProviderOidURL, found := os.LookupEnv("INVITE_IDENTITY_PROVIDER_OID_URL")
	if !found {
		err := fmt.Errorf("Environemnt variable \"INVITE_IDENTITY_PROVIDER_OID_URL\" not found")
		return config{}, err
	}
	Logger.Debug(inviteIdentityProviderOidURL)
	createIdentityProviderOidURL, found := os.LookupEnv("CREATE_IDENTITY_PROVIDER_OID_URL")
	if !found {
		err := fmt.Errorf("Environemnt variable \"CREATE_IDENTITY_PROVIDER_OID_URL\" not found")
		return config{}, err
	}
	Logger.Debug(createIdentityProviderOidURL)
	mailSupportAddress, found := os.LookupEnv("MAIL_SUPPORT_ADDRESS")
	if !found {
		err := fmt.Errorf("Environemnt variable \"MAIL_SUPPORT_ADDRESS\" not found")
		return config{}, err
	}
	Logger.Debug(mailSupportAddress)
	mailSmtpHost, found := os.LookupEnv("MAIL_SMTP_HOST")
	if !found {
		err := fmt.Errorf("Environemnt variable \"MAIL_SMTP_HOST\" not found")
		return config{}, err
	}
	Logger.Debug(mailSmtpHost)
	mailSmtpPort, found := os.LookupEnv("MAIL_SMTP_PORT")
	if !found {
		err := fmt.Errorf("Environemnt variable \"MAIL_SMTP_PORT\" not found")
		return config{}, err
	}
	Logger.Debug(mailSmtpPort)
	mailSmtpUsername, found := os.LookupEnv("MAIL_SMTP_USERNAME")
	if !found {
		err := fmt.Errorf("Environemnt variable \"MAIL_SMTP_USERNAME\" not found")
		return config{}, err
	}
	Logger.Debug(mailSmtpUsername)

	mailSmtpPassword, found := os.LookupEnv("MAIL_SMTP_PASSWORD")
	if !found {
		err := fmt.Errorf("Environemnt variable \"MAIL_SMTP_PASSWORD\" not found")
		return config{}, err
	}
	Logger.Debug(sha256.New().Write([]byte(mailSmtpPassword)))
	mailTemplateKeys, found := os.LookupEnv("MAIL_TEMPLATE_KEYS")
	if !found {
		err := fmt.Errorf("Environemnt variable \"MAIL_TEMPLATE_KEYS\" not found")
		return config{}, err
	}
	Logger.Debug(mailTemplateKeys)
	mailTemplate, found := os.LookupEnv("MAIL_TEMPLATE")
	if !found {
		err := fmt.Errorf("Environemnt variable \"MAIL_TEMPLATE\" not found")
		return config{}, err
	}
	Logger.Debug(mailTemplate)

	config := config{
		port:                         portInt,
		createAdminRolePath:          createAdminRolePath,
		inviteAdminRolePath:          inviteAdminRolePath,
		adminRoles:                   adminRoles,
		credentialMapping:            credentialMapping,
		credentialDataTemplate:       credentialDataTemplate,
		credentialEndpoint:           credentialEndpoint,
		claimMappingServiceURL:       claimMappingServiceURL,
		inviteIdentityProviderOidURL: inviteIdentityProviderOidURL, createIdentityProviderOidURL: createIdentityProviderOidURL,
		mailSupportAddress: mailSupportAddress, mailSmtpHost: mailSmtpHost, mailSmtpPort: mailSmtpPort,
		mailSmtpUsername: mailSmtpUsername, mailSmtpPassword: mailSmtpPassword,
		mailTemplate: mailTemplate, mailTemplateKeys: mailTemplateKeys,
	}

	return config, nil
}
