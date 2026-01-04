package request

import (
	"errors"
	// "fmt"
	"io"
	"strconv"
	"strings"
	"tcp/internal/headers"
)

type state int

const (
	initialized state = iota
	done
	requestStateParsingHeaders
	requestStateParsingBody
)

const bufferSize int = 8

type Request struct {
	RequestLine RequestLine
	state       int
	Headers 	headers.Headers
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
	
}

func ParseHeaders(data []byte)( headers.Headers, int, error){

	var longitudHeadersParseados int


	 headersFinal := make(map[string]string)
	if(!strings.Contains(string(data), "\r\n\r\n")){
		return nil, 0, nil
	}

	// fmt.Println("debug2   ", strings.Split(string(data), "\r\n"))

	headersTroceados := strings.Split(string(data), "\r\n")
	h := headers.NewHeaders()
	// parseado :=0
	for i:= 0; i<= len(headersTroceados); i++ {
		longitudHeadersParseados +=  len(headersTroceados)
		_, isFishedHeaders, err := h.Parse([]byte(strings.ToLower(headersTroceados[i]) + "\r\n\r\n"))
		for k, v := range h {
			headersFinal[k] = v
		}
		// fmt.Println("algo hay aqui   HEADERS ",h)
		if err != nil{
			// fmt.Println("esto es el error", err)
			return h, 0, errors.New("error: unknown state")

		}
		if isFishedHeaders {			
			return headersFinal, longitudHeadersParseados, nil
		}

	}

	return nil, 0, nil



}

func ParseBody(r Request, data []byte) ([]byte, error){


	length, err := r.Headers.Get("Content-Length")
	if err != nil {
		r.state = 1
		return nil, nil
	}
	cuerpoRaw := strings.Split(string(data), "\r\n\r\n")
	cuerpoTrimed := strings.Trim(cuerpoRaw[1],"\x00")


	lengthTrim  := strings.Trim(length, "\r\n")


	lengthInt, err := strconv.Atoi(lengthTrim)
		if err != nil {
		return nil, errors.New("Error convirtiendo la longitud")
	}


	// 	fmt.Println("a ver que a" ,cuerpoRaw[1])

	// indexFinal := strings.Index(cuerpoRaw[1], "\n")
	// if indexFinal == -1 {
	// 	return nil, errors.New("Body incompleto")
	// }

	if(lengthInt < len(cuerpoTrimed) ) {
			// fmt.Println("a ver que a" ,len(cuerpoRaw[1]), lengthInt)
			// for _, v := range strings.Trim(cuerpoRaw[1],"\x00") {
			// 	fmt.Println("valor", v)
			// }
		return nil, errors.New("La longitud no coincide")
	}

	// if(lengthInt == len(data) ) {

	return []byte(cuerpoTrimed), nil
	// }





}

func (r *Request) parse(data []byte) (int, error) {
	// if r.state == 0 {
	// 	return 0, nil
	// }

	if r.state == 0 || r.state == 2 || r.state == 3 { 
		parseado, err := parseRequestLine(data)
		// fmt.Println("esto es en parse stringfinal1", string(data[:parseado]))
		if err != nil {
			return 0, err
		}

		if parseado == 0 && err == nil {
		// fmt.Println("esto es en parse stringfinal2")

			return 0, nil
		}

		
		// fmt.Println("esto es en parse stringfinal3")
		if len(data) == parseado {
			
			r.state = 1
			
		}

		if r.state == 2 {


			headers, consumed, err := ParseHeaders(data[parseado:])
			if err != nil {
				return 0, errors.New("hubo un error parseando los headers ")
			}

			r.Headers = headers
			if headers != nil {
				r.state = 3


			}


			// r.state = 3
			return parseado + consumed, nil

			// if err != nil {
			// 	return 0, errors.New("Problema PARSEANDO LAS CABECERAS")
			// }



			// h := headers.NewHeaders()
			// conmsumidosHeaders, isFishedHeaders, err := h.Parse(data[parseado:])
			// fmt.Println("algo hay aqui   HEADERS", string(data[parseado:]))
			// r.RequestLine.Headers = h
			// if err != nil{
			// 	return 0, errors.New("error: unknown state")

			// }
		}

		if r.state == 3 {

			body, err := ParseBody(*r, data[parseado:])
			if err != nil {
				return 0, errors.New("hubo un error parseando el body ")
			}
			r.Body = body
			r.state = 1

			// fmt.Println("a ver el r", r)


			return parseado, nil

		}

		
		stringFinal := string(data[:parseado])
		// fmt.Println("esto es en parse stringfinal", stringFinal)

		partes := strings.Split(stringFinal, " ")

		if strings.Compare(partes[0], strings.ToUpper(partes[0])) != 0 {
			return 0, errors.New("algo ha pasado en el metodo")
		}

		metodo := partes[0]
		direccion := partes[1]

		// fmt.Println("esto es desde alla2", partes[2])

		versionHttp := partes[2]
		numero := []string{}
		if(strings.Contains(versionHttp, "/")){
			numero = strings.Split(versionHttp, "/")

		}
		// fmt.Println("esto es desde alla2", numero)

		// numeroValido := numero[1][:len(numero[1])-1]
		numeroValido := strings.Trim(numero[1], "\r\n")


		if numeroValido!= "1.1" {
			return 0, errors.New("Problema con la version " + numeroValido)
		}
		
		partes[2] = numeroValido
		version := partes[2]

		// fmt.Println("esto es desde alla3")
		r.RequestLine = RequestLine{
			HttpVersion:   version,
			RequestTarget: direccion,
			Method:        metodo,
		}
		// fmt.Sprintf("Request line: \r\n - Method: %s")
		r.state = 2

		// fmt.Println("esto es en el final ", parseado)

		return parseado, nil

	} else if r.state == 1 {
		return 0, errors.New("error: trying to read data in a done state")
	}else {
		return 0, errors.New("error: unknown state")

	}

}



func parseRequestLine(b []byte) (int, error) {

	// stringFinal := ""
	// partes := []string{}
	// fmt.Println("BLeeeeeeeeeeee  ", len(b))

	registeredNurse := false
	numeroProcesado := 0
	for i := 0; i < len(b); i++ {
		// parse := make([]byte, 8)
		// char := string(b[i])
		numeroProcesado++
		// fmt.Println("AAAAAAAAAAAA", "  " ,string(b[i]))
			if string(b[i]) == "\r" {
				// fmt.Println("dentro del bucle for1")

				if string(b[i+1]) == "\n" {
					registeredNurse = true

				// fmt.Println("dentro del bucle for2")

				}

				}

			if registeredNurse {
				// fmt.Println("sum0a1111  ", string(b), "  pres ",numeroProcesado)
				return numeroProcesado+1, nil
			}
		}

		// stringFinal += char
		
		// fmt.Println("BLAAAAAAAAAAAAAAAa  ", numeroProcesado)
	
		return 0, nil
	}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, bufferSize, bufferSize)
	
	readToIndex := 0
	// unread :=0
	
	r := Request{
		state: 0,
	}
	
	
	// header, err := io.ReadAll(reader)
	// if err != nil{
		// 	// fmt.Println("algo ha pasado")
		// }
		
		i:=0
		for {
			consumed, err1 := io.ReadFull(reader, buf[readToIndex:])
			_, erro := r.parse(buf)
			if erro != nil {

				return nil, errors.New("Hubo un problema parseando la request")
			}
			// fmt.Println(" a ver que devuelve el struct2", r)
			readToIndex =  consumed + readToIndex

			
			// if (consumed == len(buf[readToIndex:])) && (parsed != 0) {
			// 	readToIndex = 0
			// }
			i++
			// if(i ==3){
			// 	break
			// }

			if readToIndex >= len(buf){
				newBud := make([]byte, len(buf)*2,)
				copy(newBud, buf)
				buf = newBud
				// unread=  len(buf[readToIndex:]) - consumed  
				// readToIndex = unread
				// parsed =0
				
				if consumed == len(buf){
				// // fmt.Println("esto es el segundo parse pero primero", "  ", string(buf),"   ")


				// fmt.Println("esto es el segundo parse", string(newBud), "  ", string(buf),"   ", len(newBud))
				// readToIndex =0

				}
			}
			// fmt.Println("esto va mal", err, r)

			
			if err1 != nil && r.state == 1{
			break
				
			}

		}
	
	// fmt.Println(" a ver que devuelve el struct", r)
	return &r, nil
}