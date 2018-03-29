package main

import (
	"encoding/base64"
	"fmt"
	"compress/zlib"
	"bytes"
	"io/ioutil"
	"strings"
	"encoding/hex"
	"net/http"
)

func main() {
	//fmt.Println(testuncompressRequest(""))
	//testRoundTrip()
	testServer()
}
func testServer(){
	req,err := http.NewRequest("POST", "http://wanwang.endaosi.cn/gae_proxy.php", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Cookie", "eJwUxjEOQiEMBuDbMPe1pb8MPQwJbRg0GNB4ffO277ufbg8oFASRzslpF8gahzUQB8SY83aZ0Ufs4+Xdz/mtPRxkFwQCWILLKz5zDVdorfoPAAD//3FVGTY=")
	resp,err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	buff,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buff))
}
func testuncompressRequest(str string)map[string]string{
	str = "eJwUxjEOQiEMBuDbMPe1pb8MPQwJbRg0GNB4ffO277ufbg8oFASRzslpF8gahzUQB8SY83aZ0Ufs4+Xdz/mtPRxkFwQCWILLKz5zDVdorfoPAAD//3FVGTY="
	requestInfo := decode(uncompressRequest(str))
	return requestInfo
}
func testcompressRequest()string{
	req, err := http.NewRequest("GET", "http://ip.sb", nil)
	if err != nil {
		panic(err)
	}
	return compressRequest(encode(req))
}
func testRoundTrip(){
	req, err := http.NewRequest("GET", "https://api.ip.sb/ip", nil)
	if err != nil {
		panic(err)
	}
	resp, err := RoundTrip(req)
	if err != nil {
		panic(err)
	}
	buff, err := ioutil.ReadAll(resp.Body) // todo panic: unexpected EOF
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Printf("%s", buff)
}

func RoundTrip(req *http.Request) (*http.Response, error) {
	gaeRequest, err := http.NewRequest("POST", "http://wanwang.endaosi.cn/gae_proxy.php", req.Body)
	if err != nil {
		return nil,err
	}

	gaeRequest.Header.Add("Cookie", compressRequest(encode(req)))
	response, err := http.DefaultClient.Do(gaeRequest)

	return response, err
}

func encode(req *http.Request) map[string]string {
	requestInfo := make(map[string]string)
	requestInfo["url"] = req.URL.String()

	var headersWirer bytes.Buffer
	req.Header.Write(&headersWirer)
	requestInfo["headers"] = headersWirer.String()
	requestInfo["password"] = "passwor"
	requestInfo["method"] = req.Method
	return requestInfo
}

func compressRequest(requestInfo map[string]string) string {
	var kvs []string
	for k, v := range requestInfo {
		if v != "" {
			v = hex.EncodeToString([]byte(v))
		}
		kvs = append(kvs, fmt.Sprintf("%s=%s", k, v))
	}

	str := strings.Join(kvs, "&")

	var b bytes.Buffer

	w := zlib.NewWriter(&b)
	_, err := w.Write([]byte(str))
	if err != nil {
		panic(err)
	}
	w.Close()
	return base64.StdEncoding.EncodeToString(b.Bytes())
}
func uncompressRequest(str string) string {
	buff, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(err)
	}

	b := bytes.NewReader(buff)

	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	buff, err = ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(buff)
}
func decode(str string) (result map[string]string) {
	result = make(map[string]string)

	for _, kv := range strings.Split(str, "&") {
		var pair = strings.Split(kv, "=")
		var value string = ""

		if len(pair) != 2 {
			value = ""
		} else {
			if pair[1] != "" {
				valueBytes, err := hex.DecodeString(pair[1])
				if err != nil {
					panic(err)
				}
				value = string(valueBytes)
			}
		}
		result[pair[0]] = value
	}
	return
}
