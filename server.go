package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/yalp/jsonpath"
	"go.uber.org/zap"
)

func RequestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		targetMux.ServeHTTP(w, r)

		if string(r.RequestURI) != "/isAlive" {
			Logger.Infow("",
				zap.String("method", string(r.Method)),
				zap.String("uri", string(r.RequestURI)),
				zap.Duration("duration", time.Since(start)*1000),
			)
		}
	})
}

func startServer(port *int) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/createInvitation", createCredentialByQrPost).Methods("POST")
	router.HandleFunc("/sendInvitation", createCredentialByEmailPost).Methods("POST")
	router.HandleFunc("/createCredential", createCredentialByRegistrationPOST).Methods("POST")
	router.HandleFunc("/options", optionsGET).Methods("GET")
	router.HandleFunc("/isAlive", isAliveGet).Methods("GET")

	portString := ":" + strconv.Itoa(*port)
	log.Fatal(http.ListenAndServe(portString, RequestLogger(router)))
}

func createCredentialByQrPost(w http.ResponseWriter, r *http.Request) {
	// Get config
	config, _ := getConfig()

	token, err := GetToken(r, config.inviteIdentityProviderOidURL)
	if err != nil {
		Logger.Error(err)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(err.Error())

		return
	}

	// Check if admin
	adminRolePath := config.inviteAdminRolePath
	adminRoles := config.adminRoles

	tokenPayload, _ := json.Marshal(token.Claims.(jwt.MapClaims))
	var tokenData interface{}
	err = json.Unmarshal(tokenPayload, &tokenData)
	roles, err := jsonpath.Read(tokenData, adminRolePath)

	var adminRolesObjects []interface{}
	err = json.Unmarshal([]byte(adminRoles), &adminRolesObjects)

	isAdmin := false
	for _, adminRoleObject := range adminRolesObjects {
		for _, role := range roles.([]interface{}) {
			if role.(string) == adminRoleObject.(string) {
				isAdmin = true
			}
		}
	}

	if !isAdmin {
		http.Error(w, "Missing admin role in token payload.", http.StatusBadRequest)
		return
	}

	// Get body params
	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var bodyTemplateKeys map[string]interface{}
	bodyTemplateKeys = jsonBody["templateKeyMap"].(map[string]interface{})
	var selectedRoles ([]string)
	for _, value := range jsonBody["selectedRoleKeys"].(([]interface{})) {
		selectedRoles = append(selectedRoles, string(value.(string)))
	}

	// Create data
	var tokenDataMapping map[string]interface{}
	err = json.Unmarshal([]byte(config.credentialMapping), &tokenDataMapping)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		Logger.Info("Error unmarshalling Credential Mapping" + err.Error())
		Logger.Debug("Credential Mapping Config: " + config.credentialMapping)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(config.credentialDataTemplate), &data)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		Logger.Info("Error unmarshalling Credential Data Template" + err.Error())
		Logger.Debug("Credential Data Template Config: " + config.credentialDataTemplate)
		return
	}

	for key, value := range tokenDataMapping {
		if tokenValue, found := tokenData.(map[string]interface{})[value.(string)]; found {
			data[key] = tokenValue
		}

		if bodyTemplateKeys[key] != nil {
			data[key] = bodyTemplateKeys[key]
		}
	}

	data["Claims"] = strings.Join(selectedRoles, "|")
	payload := make(map[string]interface{})
	payload["userData"] = data

	Logger.Debug("Create Invitation")
	dataString, _ := json.Marshal(payload)

	// Create invitation
	credentialEndpoint := config.credentialEndpoint
	qrCodeData, err := createInvitation(credentialEndpoint, string(dataString), token.Raw)
	if err != nil {
		Logger.Error(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err.Error())

		return
	}

	qrCodeContent := qrCodeData["data"].(map[string]interface{})
	url := qrCodeContent["invitationUrl"].(string)
	content := make(map[string]interface{})
	content["invitationUrl"] = url
	var finalContent, _ = json.Marshal(content)
	// Generate QR Code
	if r.Header.Get("Content-Type") == "application/octect-stream" {
		qrCode, err := generateQR(url)
		if err != nil {
			Logger.Error(err)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(err.Error())

			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(qrCode)
	} else {
		if r.Header.Get("Content-Type") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			w.Write(finalContent)
		}
	}
	return
}

func createCredentialByEmailPost(w http.ResponseWriter, r *http.Request) {
	// Get config
	config, _ := getConfig()

	token, err := GetToken(r, config.inviteIdentityProviderOidURL)
	if err != nil {
		Logger.Error(err)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(err.Error())

		return
	}

	// Check if admin
	adminRolePath := config.inviteAdminRolePath
	adminRoles := config.adminRoles

	tokenPayload, _ := json.Marshal(token.Claims.(jwt.MapClaims))
	var tokenData interface{}
	err = json.Unmarshal(tokenPayload, &tokenData)
	roles, err := jsonpath.Read(tokenData, adminRolePath)

	var adminRolesObjects []interface{}
	err = json.Unmarshal([]byte(adminRoles), &adminRolesObjects)

	isAdmin := false
	for _, adminRoleObject := range adminRolesObjects {
		for _, role := range roles.([]interface{}) {
			if role.(string) == adminRoleObject.(string) {
				isAdmin = true
			}
		}
	}

	if !isAdmin {
		http.Error(w, "Missing admin role in token payload.", http.StatusBadRequest)
		return
	}

	// Get body params
	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var bodyTemplateKeys map[string]interface{}
	bodyTemplateKeys = jsonBody["templateKeyMap"].(map[string]interface{})
	var selectedRoles ([]string)
	for _, value := range jsonBody["selectedRoleKeys"].(([]interface{})) {
		selectedRoles = append(selectedRoles, string(value.(string)))
	}

	// Create data
	var tokenDataMapping map[string]interface{}
	err = json.Unmarshal([]byte(config.credentialMapping), &tokenDataMapping)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		Logger.Info("Error unmarshalling Credential Mapping" + err.Error())
		Logger.Debug("Credential Mapping Config: " + config.credentialMapping)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(config.credentialDataTemplate), &data)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		Logger.Info("Error unmarshalling Credential Data Template" + err.Error())
		Logger.Debug("Credential Data Template Config: " + config.credentialDataTemplate)
		return
	}

	for key, value := range tokenDataMapping {
		if tokenValue, found := tokenData.(map[string]interface{})[value.(string)]; found {
			data[key] = tokenValue
		}

		if bodyTemplateKeys[key] != nil {
			data[key] = bodyTemplateKeys[key]
		}
	}

	emailRecipient := bodyTemplateKeys["emailRecipient"]
	if emailRecipient.(string) == "" {
		http.Error(w, "Email address cannot be empty.", http.StatusBadRequest)
		return
	}

	data["Claims"] = strings.Join(selectedRoles, "|")
	payload := make(map[string]interface{})
	payload["userData"] = data

	Logger.Debug("Create Invitation")
	dataString, _ := json.Marshal(payload)

	// Create invitation
	credentialEndpoint := config.credentialEndpoint
	qrCodeData, err := createInvitation(credentialEndpoint, string(dataString), token.Raw)
	if err != nil {
		Logger.Error(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err.Error())

		return
	}
	Logger.Debug("Create QR")
	qrCodeContent := qrCodeData["data"].(map[string]interface{})
	url := qrCodeContent["invitationUrl"].(string)
	// Generate QR Code
	qrCode, err := generateQR(url)
	if err != nil {
		Logger.Error(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err.Error())

		return
	}

	Logger.Debug("Prepare Email")
	templateKeys := []string{}
	var emailKeysJson []interface{}
	err = json.Unmarshal([]byte(config.mailTemplateKeys), &emailKeysJson)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		Logger.Info("Error unmarshalling Mailtemplate Keys" + err.Error())
		Logger.Debug("Template Keys: " + config.mailTemplateKeys)
		return
	}

	for _, keyObject := range emailKeysJson {
		object := keyObject.(map[string]interface{})
		templateKeys = append(templateKeys, object["key"].(string))
	}
	emailTemplate := config.mailTemplate
	valueMap := make(map[string]string)
	for _, key := range templateKeys {
		if _, ok := bodyTemplateKeys[key]; ok {
			valueMap[key] = bodyTemplateKeys[key].(string)
		}
	}
	Logger.Debug("Send Email")
	err = sendEmail(emailTemplate, valueMap, base64.StdEncoding.EncodeToString(qrCode), emailRecipient.(string))

	w.WriteHeader(http.StatusOK)

	return
}

func createCredentialByRegistrationPOST(w http.ResponseWriter, r *http.Request) {
	// Get config
	config, _ := getConfig()

	token, err := GetToken(r, config.createIdentityProviderOidURL)
	if err != nil {
		Logger.Error(err)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(err.Error())

		return
	}

	// Get body params
	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var selectedRoles ([]string)
	for _, value := range jsonBody["selectedRoleKeys"].(([]interface{})) {
		selectedRoles = append(selectedRoles, string(value.(string)))
	}

	// Check roles
	adminRolePath := config.createAdminRolePath

	tokenPayload, _ := json.Marshal(token.Claims.(jwt.MapClaims))
	var tokenData interface{}
	err = json.Unmarshal(tokenPayload, &tokenData)
	roles, err := jsonpath.Read(tokenData, adminRolePath)
	var finalRoles []string
	var hasContent = false
	for _, role := range roles.([]interface{}) {
		for _, selectedRole := range selectedRoles {
			if role.(string) == selectedRole {
				finalRoles = append(finalRoles, role.(string))
				hasContent = true
			}
		}
	}

	if !hasContent {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create data
	var tokenDataMapping map[string]interface{}
	json.Unmarshal([]byte(config.credentialMapping), &tokenDataMapping)
	var data map[string]interface{}
	json.Unmarshal([]byte(config.credentialDataTemplate), &data)

	for key, value := range tokenDataMapping {
		if tokenValue, found := tokenData.(map[string]interface{})[value.(string)]; found {
			data[key] = tokenValue
		}
	}

	data["Claims"] = strings.Join(finalRoles, "|")

	payload := make(map[string]interface{})
	payload["userData"] = data

	dataString, _ := json.Marshal(payload)
	// Create invitation
	credentialEndpoint := config.credentialEndpoint
	qrCodeData, err := createInvitation(credentialEndpoint, string(dataString), token.Raw)
	if err != nil {
		Logger.Error(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err.Error())

		return
	}

	qrCodeContent := qrCodeData["data"].(map[string]interface{})
	url := qrCodeContent["invitationUrl"].(string)
	content := make(map[string]interface{})
	content["invitationUrl"] = url
	var finalContent, _ = json.Marshal(content)
	// Generate QR Code
	if r.Header.Get("Content-Type") == "application/octect-stream" {
		qrCode, err := generateQR(url)
		if err != nil {
			Logger.Error(err)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(err.Error())

			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(qrCode)
	} else {
		if r.Header.Get("Content-Type") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			w.Write(finalContent)
		}
	}

	return
}

func prepareQrOutput(w http.ResponseWriter, content string) {

}

func optionsGET(w http.ResponseWriter, r *http.Request) {
	// Get config
	config, _ := getConfig()

	w.Header().Set("Content-Type", "application/json")

	options := make(map[string]interface{})

	// Get roles
	var roles []string
	roleObjects, err := listRoles(config.claimMappingServiceURL)
	if err != nil {
		Logger.Error(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err.Error())

		return
	}
	for _, roleObject := range roleObjects {
		role := roleObject.(map[string]interface{})
		roles = append(roles, role["Role"].(string))
	}
	options["roleKeys"] = roles

	// Get template email keys
	var emailKeysJson []interface{}
	json.Unmarshal([]byte(config.mailTemplateKeys), &emailKeysJson)
	options["templateEmailKeys"] = emailKeysJson

	json.NewEncoder(w).Encode(options)

	return
}

func isAliveGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	return
}
