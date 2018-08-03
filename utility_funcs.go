package main

import (
	"time"
	"math"
	"bytes"

	"fmt"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"log"
	"encoding/json"

	"os"
	"net/url"
)

//
////Takes the date added string, and returns the days left before expiration
func converttime(dateAdded string)  float64{
	t, _ := time.Parse(timeLayout, dateAdded)
	future := t.AddDate(0, 0, 60)
	days := (future.Sub(time.Now()).Hours() / 24)
	days = toFixed(days, 2)
	return days
}

//Used for date calculation
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num * output)) / output
}


//Function to get the token used in the CloudBolt API calls
func getToken(urlStr string,  creds []byte)string{

	var urlBuffer bytes.Buffer

	proxyStr := os.Getenv("HTTPS_PROXY")
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Authenticating with CloudBolt API....")
	urlBuffer.WriteString(urlStr)
	fmt.Println(urlStr)
	urlBuffer.WriteString("/api/v2/api-token-auth/")
	req, err := http.NewRequest("POST", urlBuffer.String(), bytes.NewBuffer(creds))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatal(err)
	}


	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	token := new(Token)
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println(err)
	}

	return token.Token

}

//Wrapper for http requests - specifically the API calls
func httpGetSecure(urlStr string, token string)([]byte){

	var buffer bytes.Buffer

	proxyStr := os.Getenv("HTTPS_PROXY")
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		fmt.Println(err)
	}
	buffer.WriteString("Bearer ")
	buffer.WriteString(token)
	req, err := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("Authorization", buffer.String())
	if err != nil {
		log.Fatal(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy: http.ProxyURL(proxyURL),

	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body

}