package request

import(
	"io"
	"strings"
	"fmt"
	"errors"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(b []byte) ([]string, error){


	stringFinal := ""
	partes := []string{}
	for i:=0 ; i<len(b); i++{
		// parse := make([]byte, 8)
		char := string(b[i])
		fmt.Println("esto es desde aqui", char)
		if(char =="\r"){
			if(string(b[i])=="\n"){

				fmt.Println("lo encontre")
			}
			break
		}


		stringFinal += char



	}


		// fmt.Println("esto es desde alla", stringFinal)

	partes = strings.Split(stringFinal," ")

	if(strings.Compare(partes[0], strings.ToUpper(partes[0])) != 0 ){
		return []string{}, errors.New("algo ha pasado en el metodo")
	}

			fmt.Println("esto es desde alla2", partes[2])


	versionHttp := partes[2]
	numero := strings.Split(versionHttp,"/")
				fmt.Println("esto es desde alla2", numero)

	if(numero[1] != "1.1"){
		return []string{}, errors.New("Problema con la version")
	}
	partes[2] =numero[1]

			fmt.Println("esto es desde alla3")

	return partes, nil
	

}


func RequestFromReader(reader io.Reader) (*Request, error){

	header, err := io.ReadAll(reader)
	if err != nil{
		fmt.Println("algo ha pasado")
	}

	request, err := parseRequestLine(header)
	if err != nil{
		return nil, err
	}
	RequestLineF := RequestLine{
		HttpVersion: request[2],
		RequestTarget: request[1],
		Method: request[0],
	}

	return &Request{
		RequestLine: RequestLineF,
	}, nil


}