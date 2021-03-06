package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type ValidateFacesService struct {
}

type item struct {
	Photo string `json:"photo,omitempty"`
}

type DetectFaceRequestBody struct {
	Faces []item `json:"faces"`
}

type ResponseItem struct {
	Source int `json:"source"`
	Error  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type DetectFaceResponseBody []ResponseItem

// TODO: model object and use this func for method?
func buildDTO(value string) DetectFaceRequestBody {
	i := item{value}
	faces := []item{i}

	return DetectFaceRequestBody{Faces: faces}
}

func getBase64ImageValue(filePath string) (string, error) {
	buf := new(bytes.Buffer)

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err // TODO: proc
	}

	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	_, err = encoder.Write(fileBytes)
	if err != nil {
		return "", err // TODO: proc
	}
	err = encoder.Close()
	if err != nil {
		return "", err // TODO: proc
	}

	return buf.String(), nil
}

func (service *ValidateFacesService) Validate(filePath string) (bool, error) {
	//n1 := time.Now().UnixNano()
	//fmt.Println("filePath: ", filePath)
	//fmt.Println("start base64", n1)
	value, err := getBase64ImageValue(filePath)
	//n2 := time.Now().UnixNano()
	//fmt.Println("end base64", n2, n2-n1, float64(n2-n1)/1000)
	if err != nil {
		return false, nil
	}

	detectFaceRequestBody := buildDTO(value)
	requestBody, err := json.Marshal(detectFaceRequestBody)

	kekBody := bytes.NewBuffer(requestBody)

	//fmt.Println("parsiv: ", filePath)

	resp, err := sendRequest(kekBody)
	//fmt.Printf("Resp: %+v\n", resp)
	if err != nil {
		return false, fmt.Errorf(" ||| send request to PARSIV: %+v", err)
	}

	respbody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var responseBb DetectFaceResponseBody
	err = json.Unmarshal(respbody, &responseBb)

	fmt.Printf("body: %+v\n", responseBb)

	if resp.StatusCode != http.StatusOK {
		// TODO: log error
		return false, nil
	}

	if responseBb[0].Error.Code == -3000 {
		return true, nil
	}
	fmt.Println("base64: ", value)
	return false, nil
}

func sendRequest(body io.Reader) (*http.Response, error) {
	url := os.Getenv("PARSIV_FACE_DETECT_URL")
	username := os.Getenv("PARSIV_AUTH_USERNAME")
	password := os.Getenv("PARSIV_AUTH_PASSWORD")

	//fmt.Println("username: ", username, "password: ", password)
	//fmt.Println("url: ", url)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	//fmt.Println("Authorization: ", req.Header.Get("Authorization"))

	client := &http.Client{}

	return client.Do(req)
}
