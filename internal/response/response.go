package response

import (
	"fmt"
	"net"
	"strconv"

	// "strings"
	"tcp/internal/headers"
)

type StatusCode int

type Writer struct {
	Con net.Conn
}

const (
	Status_200 StatusCode = 200
	Status_400 StatusCode = 400
	Status_500 StatusCode = 500
)

const (
	Html_error_400 = "<html> <head> <title>400 Bad Request</title> </head> <body> <h1>Bad Request</h1> <p>Your request honestly kinda sucked.</p> </body> </html>"
	Html_error_500 = "<html> <head> <title>500 Internal Server Error</title> </head> <body> <h1>Internal Server Error</h1> <p>Okay, you know what? This one is on me.</p> </body> </html>"
	Html_succes    = "<html> <head> <title>200 OK</title> </head> <body> <h1>Success!</h1> <p>Your request was an absolute banger.</p> </body> </html>"
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	switch statusCode {
	case Status_200:
		_, err := w.Con.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			fmt.Println("ha pasado un error escribriendo el status", err)
		}
	case Status_400:
		w.Con.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case Status_500:
		w.Con.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(map[string]string)
	if contentLen != 0 {
		h["Content-Length"] = strconv.Itoa(contentLen)
	} else {
		h["Transfer-Encoding"] = "chunked"
	}

	h["Connection"] = "close"
	h["Content-Type"] = "text/html"

	// fmt.Println("esto es en headers ", h)
	return h

}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	var menesajeHeaders string
	for k, v := range headers {
		menesajeHeaders += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	headersFinal := menesajeHeaders + "\r\n"

	_, err := w.Con.Write([]byte(headersFinal))
	if err != nil {
		fmt.Println("ha habido un error escribiendo en los headers", err)
	}

	// fmt.Println("Estos son los nuevos headers ", headersFinal)

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	b, err := w.Con.Write([]byte(p))
	if err != nil {
		return 0, err
	}

	return b, nil

}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	body := fmt.Sprintf("%X\r\n%s\r\n", len(p), string(p))
	fmt.Println("eso es el body", body)

	b, err := w.Con.Write([]byte(body))
	if err != nil {
		return 0, err
	}


	return b, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	b, err := w.Con.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return 0, err
	}

	return b, nil

}
