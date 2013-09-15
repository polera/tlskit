package tlskit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type TLSRequest []struct {
	Server string `json:"server"`
	Port   int32  `json:"port"`
}

type TLSResponse []struct {
	Response
}

type Response struct {
	Server           string `json:"server"`
	ValidSince       string `json:"valid_since"`
	AtAlertThreshold bool   `json:"at_alert_threshold"`
	ExpirationDate   string `json:"expiration_date"`
	DaysValid        int32  `json:"days_valid"`
	Expired          bool   `json:"expired"`
	DaysToExpiration int32  `json:"days_to_expiration"`
	Bits             int32  `json:"bits"`
	Port             int32  `json:"port"`
}

func (r Response) String() string {
	if r.Expired {
		return fmt.Sprintf("%s has expired.", r.Server)
	}
	return fmt.Sprintf("%s is still valid.", r.Server)
}

const (
	VERSION  = "1.0.0"
	API_HOST = "tlskit.com"
	API_PATH = "/api"
)

func Lookup(request TLSRequest) ([]Response, error) {

	//	m := TLSRequest{{"uncryptic.com", 443}, {"amazon.com", 443}}
	jsonRequest, _ := json.Marshal(request)
	url := fmt.Sprintf("http://%s%s", API_HOST, API_PATH)
	req, err := http.Post(url,
		"application/json",
		bytes.NewReader(jsonRequest))

	if err != nil {
		return nil, err
	}

	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return nil, errors.New(req.Status)
	}

	r := new(TLSResponse)

	err = json.NewDecoder(req.Body).Decode(r)

	if err != nil {
		return nil, err
	}
	responses := make([]Response, len(*r))
	for i, child := range *r {
		responses[i] = child.Response
	}
    return responses, nil
}
