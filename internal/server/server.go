package server

import (
	"fmt"
	"sync/atomic"
	// "io"
	"net"
	"strconv"
	"tcp/internal/response"
)


type Server struct {
	state net.Listener
	closed   atomic.Bool
}

// func conect(port int) (net.Listener, error){
// 	tcp, err := net.Listen("tcp", ":"+ string(port))
// 		if err != nil {
// 			fmt.Println("No se ha podido conectar")
// 			return nil, errors.New("No se ha podido conectar")
// 		}

// }

func Serve(port int) (*Server, error){

	// server := Server{}
	ch := make(chan Server)
	server, err := net.Listen("tcp", ":"+ strconv.Itoa(port))
	s := Server{ state: server}

	go func(){
		if err != nil {
			fmt.Println("No se ha podido conectar ", err)
			// return nil, errors.New("No se ha podido conectar")
		}
		s.listen()
		ch <- Server{state: server}
		}()

		
		
		return &s, nil

}


func (s *Server) Close() error{
	s.closed.Store(true)
	
    if s.state != nil {
        return s.state.Close()
    }
    return nil


}

func (s *Server) listen(){
	// ch := make(chan Server)
		for{
			fmt.Println("estoes aqui22")
		conec, err := s.state.Accept()
			 if err != nil {
            if s.closed.Load() {
                return
            }
            fmt.Println("error accepting connection:", err)
            continue
        }

		go func(){
			fmt.Println("estoes aqui33")
			 s.handle(conec)
			 }()
			}

}







func (s *Server) handle(conn net.Conn){
	// mensaje := ("HTTP/1.1 200 OK\nContent-Type: text/plain\n\n Hello World!")
	defer conn.Close()


	errS := response.WriteStatusLine(conn, response.Status_200 )
	if errS != nil {
		fmt.Println("ha habido un error escribiendo el status", errS)
		return
	}

	headersResponse := response.GetDefaultHeaders(0)
	errH := response.WriteHeaders(conn, headersResponse)
	// fmt.Println("headersResponse ", headersResponse)
	if errH != nil {
		fmt.Println("ha habido un error escribiendo los headers", errH)
		return
	}

}

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




