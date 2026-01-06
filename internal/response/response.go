package response

import (
	"fmt"
	"io"
	"strconv"
	// "strings"
	"tcp/internal/headers"
)

type StatusCode int


const (
	Status_200 StatusCode = 200 + iota
	status_400
	status_500
	
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case Status_200:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			fmt.Println ("ha pasado un error escribriendo el status", err)
		}
	case status_400:
		w.Write([]byte("HTTP/1.1 400 Bad Request"))
	case status_500:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error"))
	}
	
	return nil
}


func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(map[string]string)
	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	// fmt.Println("esto es en headers ", h)
	return h

}



func WriteHeaders(w io.Writer, headers headers.Headers) error{

	var menesajeHeaders string; 
	for k,v := range headers {
		menesajeHeaders += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	headersFinal := menesajeHeaders + "\r\n"

	_, err := w.Write([]byte(headersFinal))
	if err != nil {
		fmt.Println("ha habido un error escribiendo en los headers", err)
	}

	// fmt.Println("Estos son los nuevos headers ", headersFinal)

	
	return nil
}