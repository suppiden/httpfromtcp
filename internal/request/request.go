package request

import (
	"errors"
	"fmt"

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
	Headers     headers.Headers
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func ParseHeaders(data []byte) (headers.Headers, int, error) {
	headersFinal := make(map[string]string)

	if !strings.Contains(string(data), "\r\n\r\n") {
		return nil, 0, nil
	}

	// Find where headers end (double CRLF)
	headerEndIndex := strings.Index(string(data), "\r\n\r\n")
	if headerEndIndex == -1 {
		return nil, 0, nil
	}

	// Extract just the headers portion (after request line)
	headersData := string(data[:headerEndIndex])
	headersTroceados := strings.Split(headersData, "\r\n")

	h := headers.NewHeaders()

	for i := 0; i < len(headersTroceados); i++ {
		line := headersTroceados[i]
		if line == "" {
			continue
		}
		_, isFinishedHeaders, err := h.Parse([]byte(strings.ToLower(line) + "\r\n\r\n"))
		if err != nil {
			return nil, 0, err
		}
		for k, v := range h {
			headersFinal[k] = v
		}
		if isFinishedHeaders {
			break
		}
	}

	// Return total bytes consumed (headers + CRLF CRLF)
	return headersFinal, headerEndIndex + 4, nil
}

func ParseBody(r Request, data []byte) ([]byte, error) {

	length, err := r.Headers.Get("Content-Length")
	fmt.Println("AAAAAAAAAA " , string(data))
	if err != nil {
		fmt.Println("a ver que a err ", err)
		r.state = 1
		return nil, nil
	}
	cuerpoRaw := strings.Split(string(data), "\r\n\r\n")
	cuerpoTrimed := strings.Trim(cuerpoRaw[1], "\x00")
	fmt.Println("a ver que a vcc ", cuerpoTrimed)

	if cuerpoTrimed == "" {
		return nil, nil
	}

	lengthTrim := strings.Trim(length, "\r\n")

	lengthInt, err := strconv.Atoi(lengthTrim)
	if err != nil {
		return nil, errors.New("Error convirtiendo la longitud")
	}

	// fmt.Println("a ver que a   " ,cuerpoRaw[1])

	// indexFinal := strings.Index(cuerpoRaw[1], "\n")
	// if indexFinal == -1 {
	// 	return nil, errors.New("Body incompleto")
	// }

	if lengthInt < len(cuerpoTrimed) {
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

	switch r.state {
	case 0:
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
		if strings.Contains(versionHttp, "/") {
			numero = strings.Split(versionHttp, "/")

		}
		// fmt.Println("esto es desde alla2", numero)

		// numeroValido := numero[1][:len(numero[1])-1]
		numeroValido := strings.Trim(numero[1], "\r\n")

		if numeroValido != "1.1" {
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

	case 2:
		fmt.Println("aaaaaaaaaaaaaaw")

		headers, consumed, err := ParseHeaders(data)
		if err != nil {
			fmt.Println("aaaaaaaaaaaaaa1")

			return 0, errors.New("hubo un error parseando los headers ")
		}

		r.Headers = headers
		if headers != nil {
			fmt.Println("aaaaaaaaaaaaaa2")

			r.state = 3

		}
		fmt.Println("aaaaaaaaaaaaaaas", r.Headers)

		// r.state = 3
		return consumed, nil

	case 3:

		fmt.Println("aaaaaaaaaaaaaa")

		body, err := ParseBody(*r, data)
		if err != nil {
			return 0, errors.New("hubo un error parseando el body ")
		}

		if body == nil {
			r.state = 1
			return 0, nil
		}
		r.Body = body
		r.state = 1

		// fmt.Println("a ver el r", r)

		return len(body), nil

	case 1:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")

	}

}

func parseRequestLine(b []byte) (int, error) {

	// stringFinal := ""
	// partes := []string{}
	// fmt.Println("BLeeeeeeeeeeee  ", len(b))

	registeredNurse := false
	numeroProcesado := 0
	for i := range b {
		// parse := make([]byte, 8)
		// char := string(b[i])
		numeroProcesado++
		fmt.Println("AAAAAAAAAAAA", "  " ,string(b[i]))
		if string(b[i]) == "\r" && i != len(b)-1 {
			// fmt.Println("dentro del bucle for1")

			if string(b[i+1]) == "\n" {
				registeredNurse = true

				// fmt.Println("dentro del bucle for2")

			}

		}

		if registeredNurse {
			// fmt.Println("sum0a1111  ", string(b), "  pres ",numeroProcesado)
			return numeroProcesado + 1, nil
		}
	}

	// stringFinal += char

	// fmt.Println("BLAAAAAAAAAAAAAAAa  ", numeroProcesado)

	return 0, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	fullData := make([]byte, 0)

	r := Request{
		state: 0,
	}

	for r.state != int(done) {
		// Read available data
		n, err := reader.Read(buf)
		if n > 0 {
			fullData = append(fullData, buf[:n]...)
		}

		// Try to parse what we have
		for {
			parsed, parseErr := r.parse(fullData)
			if parseErr != nil {
				return nil, parseErr
			}
			if parsed == 0 {
				// Need more data
				break
			}
			// Remove parsed data from buffer
			fullData = fullData[parsed:]

			if r.state == int(done) {
				return &r, nil
			}
		}

		if err != nil {
			if err == io.EOF {
				// No more data, check if we're done
				if r.state == int(done) || r.state == int(requestStateParsingBody) {
					r.state = int(done)
					return &r, nil
				}
			}
			return nil, err
		}
	}

	return &r, nil
}
