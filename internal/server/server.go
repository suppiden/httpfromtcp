package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"

	// "io"
	"net"
	"strconv"
	"tcp/internal/request"
	"tcp/internal/response"
)

type Server struct {
	state  net.Listener
	closed atomic.Bool
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	code    response.StatusCode
	message string
}

func HandlerFunc(w *response.Writer, req *request.Request) {

	// fmt.Println("Eso es reques line ", req.RequestLine.RequestTarget )

	if req.RequestLine.RequestTarget == "/yourproblem" {
			w.WriteStatusLine(response.Status_400)
			headersDefault := response.GetDefaultHeaders(len(response.Html_error_400))
			w.WriteHeaders(headersDefault)
			w.WriteBody([]byte(response.Html_error_400))
			return
		
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
			w.WriteStatusLine(response.Status_500)
			headersDefault := response.GetDefaultHeaders(len(response.Html_error_500))
			w.WriteHeaders(headersDefault)
			w.WriteBody([]byte(response.Html_error_500))
			return
		
	}



	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {

			buf := make([]byte, 1024)
			fullData := make([]byte, 0)


			resp, err1 := http.Get("https://httpbin.org/stream/100")
			
			
			
			
			
			
			w.WriteStatusLine(response.Status_200)
			headersDefault := response.GetDefaultHeaders(0)
			w.WriteHeaders(headersDefault)
			for {
				n, err := resp.Body.Read(buf)
				fmt.Println(string(buf))
				if err1 != nil {
					fmt.Println("Esta vez este es el error", err)
				}
				
				if n > 0 {
					fullData = append(fullData, buf[:n] ...)
				}
				n, errCh :=	w.WriteChunkedBody(fullData)
				if n < 0 || errCh != nil {
					break
				}
				
				if err != nil {
					fmt.Println("no he sacado info ", err, io.EOF)
					if(err == io.EOF){
						_, err := w.WriteChunkedBodyDone()

						if err != nil {
							fmt.Println("Some errors writing the end of chunk", err)
						}
						break
					}
				}
			}
			
		return
	}

	w.WriteStatusLine(response.Status_200)
	headersDefault := response.GetDefaultHeaders(len(response.Html_succes))
	w.WriteHeaders(headersDefault)
	_, err := w.WriteBody([]byte(response.Html_succes))
// _, err := w.Write([]byte("All good, frfr\n"))


	if err != nil {
	
		w.WriteStatusLine(response.Status_500)
		headersDefault := response.GetDefaultHeaders(len(response.Html_error_500))
		w.WriteHeaders(headersDefault)
		w.WriteBody([]byte(response.Html_error_500))
		
	}
}

// func conect(port int) (net.Listener, error){
// 	tcp, err := net.Listen("tcp", ":"+ string(port))
// 		if err != nil {
// 			fmt.Println("No se ha podido conectar")
// 			return nil, errors.New("No se ha podido conectar")
// 		}

// }

func Serve(port int, handler Handler) (*Server, error) {

	// server := Server{}
	ch := make(chan Server)
	server, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	s := Server{state: server}

	go func() {
		if err != nil {
			fmt.Println("No se ha podido conectar ", err)
			// return nil, errors.New("No se ha podido conectar")
		}
		s.listen(handler)
		ch <- Server{state: server}
	}()

	return &s, nil

}

// func Serve(port int) (*Server, error){

// 	// server := Server{}
// 	ch := make(chan Server)
// 	server, err := net.Listen("tcp", ":"+ strconv.Itoa(port))
// 	s := Server{ state: server}

// 	go func(){
// 		if err != nil {
// 			fmt.Println("No se ha podido conectar ", err)
// 			// return nil, errors.New("No se ha podido conectar")
// 		}
// 		s.listen()
// 		ch <- Server{state: server}
// 		}()

// 		return &s, nil

// }

func (s *Server) Close() error {
	s.closed.Store(true)

	if s.state != nil {
		return s.state.Close()
	}
	return nil

}

func (s *Server) listen(handler Handler) {
	for {
		fmt.Println("estoes aqui22")
		conec, err := s.state.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Println("error accepting connection:", err)
			continue
		}

		go func() {
			fmt.Println("estoes aqui33")
			s.handle(conec, handler)
		}()
	}

}




func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()
	wr := response.Writer{
		Con : conn,
	}

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println("este es el error ", err)
		return
	}

	// Create a buffer for the handler to write to
	// buf := &bytes.Buffer{}

	// Call the handler
	 handler(&wr, req)

	// if handlerErr != nil {
	// 	// Handler returned an error - write error response
	// 	errS := wr.WriteStatusLine(handlerErr.code)
	// 	if errS != nil {
	// 		fmt.Println("ha habido un error escribiendo el status", errS)
	// 		return
	// 	}
	// 	headersResponse := response.GetDefaultHeaders(len(handlerErr.message))
	// 	errH := wr.WriteHeaders(headersResponse)
	// 	if errH != nil {
	// 		fmt.Println("ha habido un error escribiendo los headers", errH)
	// 		return
	// 	}
	// 	conn.Write([]byte(handlerErr.message))
	// 	return
	// }

	// Handler succeeded - write success response
	// body := buf.Bytes()
	// errS := wr.WriteStatusLine(response.Status_200)
	// if errS != nil {
	// 	fmt.Println("ha habido un error escribiendo el status", errS)
	// 	return
	// }

	// headersResponse := response.GetDefaultHeaders(len(body))
	// errH := wr.WriteHeaders(headersResponse)
	// if errH != nil {
	// 	fmt.Println("ha habido un error escribiendo los headers", errH)
	// 	return
	// }

	// // Write the response body
	// conn.Write(body)
}

// func (s *Server) handle(conn net.Conn){
// 	defer conn.Close()

// 	errS := response.WriteStatusLine(conn, response.Status_200 )
// 	if errS != nil {
// 		fmt.Println("ha habido un error escribiendo el status", errS)
// 		return
// 	}

// 	headersResponse := response.GetDefaultHeaders(0)
// 	errH := response.WriteHeaders(conn, headersResponse)
// 	if errH != nil {
// 		fmt.Println("ha habido un error escribiendo los headers", errH)
// 		return
// 	}

// }

// func WriteHeaders(conn net.Conn, headersResponse any) any {
// 	panic("unimplemented")
// }

// func GetDefaultHeaders(i int) any {
// 	panic("unimplemented")
// }

// func WriteStatusLine(conn net.Conn, StatusCode int) any {
// 	panic("unimplemented")
// }
// func (s *Server) handle(conn net.Conn){
// 	mensaje := ("HTTP/1.1 200 OK\nContent-Type: text/plain\nContent-Length: 13 \nHello World!")

// 	i, errr := conn.Write([]byte(mensaje))
// 	if errr != nil {
// 		fmt.Println("No se ha impreso i", i, errr)
// 	}
// 	err := s.Close()
// 	if err != nil {
// 		fmt.Println("algo ha pasado")
// 	}

// }

// func WriteStatusLine(w io.Writer) {
// 	mensaje := ("HTTP/1.1 200 OK\nContent-Type: text/plain\n\n Hello World!")
// 	 w.Write([]byte(mensaje))

// }
