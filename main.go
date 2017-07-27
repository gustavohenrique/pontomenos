package main

import (
        "encoding/json"
        "errors"
	"time"
        "fmt"
        "log"
	"os"

        "github.com/parnurzeal/gorequest"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Credential struct {
	Email string `json:email`
	Password string `json:password`
	Token string
	ClientID string
}

func getServer() *gin.Engine {
	server := gin.Default()
        server.Use(cors.New(cors.Config{
                AllowAllOrigins: true,
                AllowMethods:    []string{"GET", "POST"},
                AllowHeaders:    []string{"Location", "Accept", "Authorization", "Content-Type"},
                ExposeHeaders:   []string{"Link", "Location"},
                MaxAge:          1 * time.Hour,
        }))
        return server
}

func authenticateBy(credential Credential) (error, Credential) {
	url := "http://api.pontomaisweb.com.br/api/auth/sign_in"
	requestData := fmt.Sprintf(`{
		"email": "%s",
		"password": "%s"}`,
		credential.Email, credential.Password)
	request := gorequest.New()

	response, body, errs := request.Post(url).
				Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0.1; MotoG3 Build/MOB31K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/51.0.2704.106 Mobile Safari/537.36").
				Set("Content-Type", "application/json").
				Set("X-Requested-With", "br.com.pontomais.pontomais").
				Send(requestData).
				End()
	var err error
	if len(errs) > 0 {
		log.Println("[ERROR] Falhou autenticando na api do pontomais.")
		return err, credential
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(body), &result)

	if response.StatusCode > 201 || err != nil {
		msg := fmt.Sprintf("HttpStatus: %d. Email: %s. error: %s", response.StatusCode, credential.Email, err)
		err = errors.New(msg)
		return err, credential
	}

	credential.Token = result["token"].(string)
	credential.ClientID = result["client_id"].(string)
	return err, credential
}

func registerTimeclockBy(credential Credential) (error, time.Time) {
	url := "https://api.pontomaisweb.com.br/api/time_cards/register"

	requestData := fmt.Sprintf(`{
		"_path": "/meu_ponto/registro_de_ponto",
		"time_card": {
			"accuracy": 600,
			"accuracy_method": true,
			"address": "Av. das Nações Unidas, 11541 - Cidade Monções, São Paulo - SP, Brasil",
			"latitude": -23.6015797,
			"location_edited": false,
			"longitude": -46.694767,
			"original_address": "Av. das Nações Unidas, 11541 - Cidade Monções, São Paulo - SP, Brasil",
			"original_latitude": -23.6015797,
			"original_longitude": -46.694767,
			"reference_id": null
		}}`)
	request := gorequest.New()

	response, body, errs:= request.Post(url).
				Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0.1; MotoG3 Build/MOB31K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/51.0.2704.106 Mobile Safari/537.36").
				Set("Content-Type", "application/json").
				Set("X-Requested-With", "br.com.pontomais.pontomais").
				Set("token-type", "Bearer").
				Set("uid", credential.Email).
				Set("access-token", credential.Token).
				Set("client", credential.ClientID).
				Send(requestData).
				End()
	var err error
	if len(errs) > 0 {
	    log.Println("[ERROR] Falhou ao registrar ponto na api do pontomais.", errs[0])
	    return err, time.Now()
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(body), &result)

	if response.StatusCode > 201 || err != nil {
		msg := fmt.Sprintf("HttpStatus: %d. Email: %s. error: %s", response.StatusCode, credential.Email, err)
		err = errors.New(msg)
		return err, time.Now()
	}

	timecard := result["untreated_time_card"].(map[string]interface{})
	created_at := timecard["created_at"].(string)
	log.Println(created_at + ":" + credential.Email + ":" + credential.Password)
	var t time.Time
	t, err = time.Parse(time.RFC3339Nano, created_at)

	return err, t
}

func makeResponse(message string) interface{} {
	return map[string]string{"message": message}
}

func RegisterTimeclock(ctx *gin.Context) {
	credential := Credential{}
	ctx.BindJSON(&credential)

	var err error
	err, credential = authenticateBy(credential)

	status := 200
	message := "Success"

	if err != nil {
		status = 403
		message = fmt.Sprintf("Autenticacao falhou: %s", err)
		log.Println("[ERROR]", message)
	}

	err, t := registerTimeclockBy(credential)
	if err != nil {
		status = 500
		message = fmt.Sprintf("Falhou ao registrar o ponto. email: %s. token: %s. client_id: %s. error: %s", credential.Email, credential.Token, credential.ClientID, err)
		log.Println("[ERROR]", message)
	}

	message = fmt.Sprintf("%s", t.Format("20060102150405"))

	ctx.JSON(status, makeResponse(message))
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	server := getServer()
	server.POST("/", RegisterTimeclock)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "7000"
	}
	log.Println("Listening on", port)
	server.Run(":"+port)
}
