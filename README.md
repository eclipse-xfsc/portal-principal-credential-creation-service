
## GoLand
1. Install dependencies with command `sudo go install`
2. Make sure that current user is owner of /usr/local/go folder: `sudo chown :username /usr/local/go/*/*`
3. Set Environment variables:
```
    PORT=8080
    ADMIN_ROLES=[\"gaia-x-visitor\"]
    CREATE_ADMIN_ROLE_PATH=$.roles
    INVITE_ADMIN_ROLE_PATH=$.roles
    CREDENTIAL_MAPPING": "{\"FirstName\":\"given_name\",\"LastName\":\"family_name\",\"FederationId\":\"fedId\",\"MiddleName\":\"middle_name\",\"PreferredUsername\":\"preferred_username\",\"Gender\":\"gender\",\"Email\":\"email\",\"Birthdate\":\"birthdate\"}
    CREDENTIAL_DATA_TEMPLATE={"FirstName":"given_name","LastName":"family_name","FederationId":"fedId","MiddleName":"middle_name","PreferredUsername":"preferred_username","Gender":"gender","Email":"email","Birthdate":"birthdate"}
    CREDENTIAL_ENDPOINT=http://localhost:3008/v1/map-user-info
    CLAIM_MAPPING_SERVICE_URL=http://localhost:8080
    INVITE_IDENTITY_PROVIDER_OID_URL=https://aas-integration.gxfs.dev
    CREATE_IDENTITY_PROVIDER_OID_URL=https://sso-integration.gxfs.dev/realms/intranet
    MAIL_SUPPORT_ADDRESS=portal@gxfs.dev
    MAIL_SMTP_PORT=587
    MAIL_SMTP_HOST=in-v3.mailjet.com
    MAIL_SMTP_USERNAME=USER
    MAIL_SMTP_PASSWORD=PASSWORD
    MAIL_TEMPLATE_KEYS=[{\"key\":\"qr\"},{\"key\":\"message\"}]
    MAIL_TEMPLATE=<html><head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" /><title>GXFS Integration Invitation</title></head><body><p> {{index . \"message\" }} <p><img src=\"cid:invitation.png\" alt=\"Go gopher\" /></p></p></body>
```

## Visual Studio Code
1. Copy launch.json to folder .vscode
2. Install delve - a go debugger
    `brew install delve

Content of launch.json:

```
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}",
            "env": {
                "PORT": "8080",
                "ADMIN_ROLES": "[\"gaia-x-visitor\"]",
                "CREATE_ADMIN_ROLE_PATH": "$.roles",
                "INVITE_ADMIN_ROLE_PATH": "$.roles",
                "CREDENTIAL_MAPPING": "{\"FirstName\":\"given_name\",\"LastName\":\"family_name\",\"FederationId\":\"fedId\",\"MiddleName\":\"middle_name\",\"PreferredUsername\":\"preferred_username\",\"Gender\":\"gender\",\"Email\":\"email\",\"Birthdate\":\"birthdate\"}",
                "CREDENTIAL_DATA_TEMPLATE": "{\"FirstName\":\"given_name\",\"LastName\":\"family_name\",\"FederationId\":\"fedId\",\"MiddleName\":\"middle_name\",\"PreferredUsername\":\"preferred_username\",\"Gender\":\"gender\",\"Email\":\"email\",\"Birthdate\":\"birthdate\"}",
                "CREDENTIAL_ENDPOINT": "http://localhost:3008/v1/map-user-info",
                "CLAIM_MAPPING_SERVICE_URL": "http://localhost:8080",
                "INVITE_IDENTITY_PROVIDER_OID_URL": "https://aas-integration.gxfs.dev",
                "CREATE_IDENTITY_PROVIDER_OID_URL": "https://sso-integration.gxfs.dev/realms/intranet",
                "MAIL_SUPPORT_ADDRESS": "portal@gxfs.dev",
                "MAIL_SMTP_PORT": "587",
                "MAIL_SMTP_HOST": "in-v3.mailjet.com",
                "MAIL_SMTP_USERNAME": "USER",
                "MAIL_SMTP_PASSWORD": "PASSWORD",
                "MAIL_TEMPLATE_KEYS": "[{\"key\":\"qr\"},{\"key\":\"message\"}]",
                "MAIL_TEMPLATE": "<html><head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" /><title>GXFS Integration Invitation</title></head><body><p> {{index . \"message\" }} <p><img src=\"cid:invitation.png\" alt=\"Go gopher\" /></p></p></body>"
            }
        }
    ]
}

```

