package main

import (
	"encoding/base64"
	"fmt"
	"compress/zlib"
	"bytes"
	"io/ioutil"
	"strings"
	"encoding/hex"
)

func main() {
	str := "eNpljkEOwzAIBH+TMwHDJgc/xqqxcmiVym7V7xdXuVUrIbQMC+9+z7YhhQgihRs3W0G2s9sOYocYhxf9cnip3kdOmzUIUuD0D1OlklaTkEboj+IZXOZEFdNnrsHA1PzKEYRrN24CdlFZo9LcWJ5ljM/Za8a8JiHEfV4e/jrOmuN31fQFA80zYA=="

	requestInfo := decode(string(uncompress(str)))

	fmt.Print(requestInfo)
}

func uncompress(str string) []byte{
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
	return buff
}
func decode(str string)(result map[string] string){
	result = make(map[string] string)

	for _, kv := range strings.Split(str, "&") {
		var pair = strings.Split(kv, "=")
		var value string = ""

		if len(pair) != 2 {
			value = ""
		}else {
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