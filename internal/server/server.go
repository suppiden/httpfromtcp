package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"tcp/internal/headers"
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



	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream/100") {

			buf := make([]byte, 1024)
			fullData := make([]byte, 0)
			trailerTrailer := make(map[string]string) 


			resp, err1 := http.Get("https://httpbin.org/stream/100")
			
			
			
			
			
			
			w.WriteStatusLine(response.Status_200)
			headersDefault := response.GetDefaultHeaders(0)
			headersDefault["Trailer"] = "X-Content-Sha256, X-Content-Length"
			w.WriteHeaders(headersDefault)
			// w.WriteHeaders(trailerHeader)
			lengthBytes := make([]byte,0)
			for {

				n, err := resp.Body.Read(buf)
				if err1 != nil {
					fmt.Println("Esta vez este es el error", err)
				}
				
				if n > 0 {
					fullData = append(fullData, buf[:n] ...)
					n, errCh :=	w.WriteChunkedBody(buf[:n])
					if n < 0 || errCh != nil {
						break
					}
				}
				lengthBytes = append(lengthBytes, buf[:n]...)
				
				if err != nil {
					fmt.Println("no he sacado info ", err, io.EOF)
					if(err == io.EOF){
						_, err := w.WriteChunkedBodyDone()
						if err != nil {
							fmt.Println("Some errors writing the end of chunk", err)
						}
						
						hash := sha256.Sum256((lengthBytes))
						hashFinal := fmt.Sprintf("%x", hash)
						trailerTrailer["X-Content-Sha256"] = hashFinal

						lengthFinal := fmt.Sprintf("%d", len(lengthBytes))
						trailerTrailer["X-Content-Length"] = lengthFinal
						fmt.Println("no he sacado info ")
						// errTrailer := w.WriteTrailers(trailerTrailer)
						errTrailer := w.WriteTrailers(headers.Headers{
						"X-Content-Sha256":  hashFinal,
						"X-Content-Length":  strconv.Itoa(len(lengthBytes)),
						})

						fmt.Println("no he sacado info3 ")
						if errTrailer != nil {
							fmt.Println("an error ocurred writing the trailers ", errTrailer)
						}
						break
					}
				}
			}
			
		return
	}


		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/html") {

			buf := make([]byte, 1024)
			fullData := make([]byte, 0)
			trailerTrailer := make(map[string]string) 


			resp, err1 := http.Get("https://httpbin.org/html")
			
			
			
			
			
			
			w.WriteStatusLine(response.Status_200)
			headersDefault := response.GetDefaultHeaders(0)
			headersDefault["Trailer"] = "X-Content-Sha256, X-Content-Length"
			w.WriteHeaders(headersDefault)
			// w.WriteHeaders(trailerHeader)
			lengthBytes := make([]byte,0)
			for {

				n, err := resp.Body.Read(buf)
				if err1 != nil {
					fmt.Println("Esta vez este es el error", err)
				}
				
				if n > 0 {
					fullData = append(fullData, buf[:n] ...)
					n, errCh :=	w.WriteChunkedBody(buf[:n])
					if n < 0 || errCh != nil {
						break
					}
				}
				lengthBytes = append(lengthBytes, fullData...)
				
				if err != nil {
					fmt.Println("no he sacado info ", err, io.EOF)
					if(err == io.EOF){
						_, err := w.WriteChunkedBodyDone()
						if err != nil {
							fmt.Println("Some errors writing the end of chunk", err)
						}
						
						hash := sha256.Sum256((fullData))
						hashFinal := fmt.Sprintf("%x", hash)
						trailerTrailer["X-Content-Sha256"] = hashFinal

						lengthFinal := fmt.Sprintf("%d", len(fullData))
						trailerTrailer["X-Content-Length"] = lengthFinal
						fmt.Println("no he sacado info1 ")
						// errTrailer := w.WriteTrailers(trailerTrailer)
						errTrailer := w.WriteTrailers(headers.Headers{
						"X-Content-Sha256":  hashFinal,
						"X-Content-Length":  strconv.Itoa(len(fullData)),
						})

						fmt.Println("no he sacado info3 ")
						if errTrailer != nil {
							fmt.Println("an error ocurred writing the trailers ", errTrailer)
						}
						break
					}
				}
			}
			
		return
	}


	if req.RequestLine.RequestTarget == "/video" {
			w.WriteStatusLine(response.Status_200)
			
			
			by, err := os.ReadFile("../../cmd/httpserver/assets/vim.mp4")
			headersDefault := response.GetDefaultHeaders(len(by))
			headersDefault["Content-Type"] = "video/mp4"
			w.WriteHeaders(headersDefault)
			if err != nil {
				fmt.Println(" No se ha podido extraer el fichero ", err)
			}
			n, errCh :=	w.WriteBody(by)
				if n < 0 || errCh != nil {
					fmt.Println(" No se ha podido enviar el video ", errCh)
				}

			// _, errChFinal := w.WriteChunkedBodyDone()
			// 	if errChFinal != nil {
			// 		fmt.Println("Some errors writing the end of chunk", errChFinal)
			// 	}

				return

			
		
	}


	// 	if req.RequestLine.RequestTarget == "/video" {
	// 		w.WriteStatusLine(response.Status_200)
			
			
	// 		by, err := os.ReadFile("../../cmd/httpserver/assets/vim.mp4")
	// 		headersDefault := response.GetDefaultHeaders(0)
	// 		headersDefault["Content-Type"] = "video/mp4"
	// 		w.WriteHeaders(headersDefault)
	// 		if err != nil {
	// 			fmt.Println(" No se ha podido extraer el fichero ", err)
	// 		}
	// 		n, errCh :=	w.WriteChunkedBody(by)
	// 			if n < 0 || errCh != nil {
	// 				fmt.Println(" No se ha podido enviar el video ", errCh)
	// 			}

	// 		_, errChFinal := w.WriteChunkedBodyDone()
	// 			if errChFinal != nil {
	// 				fmt.Println("Some errors writing the end of chunk", errChFinal)
	// 			}

	// 			return

			
		
	// }

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
