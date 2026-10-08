package sshlistener

import(
	"log"
	"errors"
	"strings"
	"golang.org/x/crypto/ssh"
	//"net/http"
	"net/url"
	//"encoding/json"
	//"io"
)

type RequestedResource struct{
	Uname string
	EnvName string
	EnvVersion string
}

func unameParser(input string) (RequestedResource, error) {
	println("Parsing uname:", input)
	var returnable RequestedResource
	var tmp []string
	tmp = strings.Split(input, "/")

	if len(tmp) != 3 {
		log.Default().Println("Invalid input string in auth cb")
		return returnable, errors.New("invalid input format")
	}
	
	returnable.Uname = tmp[0]
	returnable.EnvName = tmp[1]
	returnable.EnvVersion = tmp[2]

	if returnable.Uname == "" || returnable.EnvName == "" || returnable.EnvVersion == "" {
		log.Default().Println("Invalid input string in auth cb")
		return returnable, errors.New("invalid input format")
	}
	
	return returnable, nil
}




func KeycloacDirectAuth(username, password string) (bool, error) {

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", "{addLater}")
	data.Set("username", username)
	data.Set("password", password)


	//req, err := http.NewRequest("POST", "http://example.com/realms/{addLater}protocol/openid-connect/token", strings.NewReader(data.Encode()))
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Accept", "application/json")
	//if err != nil {
	//	log.Default().Println("keycloak unresponsive")
	//	return false, err
	//}

	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Default().Println("keycloak request failed")
	//	return false, err
	//}
	//defer resp.Body.Close()

	//if resp.StatusCode != http.StatusOK {
	//	log.Default().Println("keycloak authentication failed")
	//	return false, errors.New("authentication failed")
	//}

	//body, err := io.ReadAll(resp.Body)
	//var result map[string]interface{}
	//json.Unmarshal(body, &result)

	return true, nil
}



func AuthCb(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {

	requestedResource, err := unameParser(conn.User())
	if err != nil {
		if err.Error() == "invalid input format" {
			log.Default().Println("Authentication failed due to invalid input format")
		}
		return nil, err
	}

	authState, err := KeycloacDirectAuth(requestedResource.Uname, string(password))
	if err != nil {
		return nil, err
	}
	if !authState {
		return nil, errors.New("authentication failed")
	}

	// dont get what perms are for, guess when i run it it will just blow up on me
	return &ssh.Permissions{}, nil
}