package main


import(
	"fmt"
	"net"
	"errors"
	"tcp/internal/request"
	// "io"
	// "strings"
)
// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	lines := make(chan string)
// 	go func() {
// 		defer f.Close()
// 		defer close(lines)
// 		currentLineContents := ""
// 		for {
// 			b := make([]byte, 8, 8)
// 			n, err := f.Read(b)
// 			if err != nil {
// 				if currentLineContents != "" {
// 					lines <- currentLineContents
// 				}
// 				if errors.Is(err, io.EOF) {
// 					break
// 				}
// 				fmt.Printf("error: %s\n", err.Error())
// 				return
// 			}
// 			str := string(b[:n])
// 			parts := strings.Split(str, "\n")
// 			for i := 0; i < len(parts)-1; i++ {
// 				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
// 				currentLineContents = ""
// 			}
// 			currentLineContents += parts[len(parts)-1]
// 		}
// 	}()
// 	return lines
// }

func main(){
	tcp, err := net.Listen("tcp", ":42069")

	if err !=nil {
		errors.New("lgo ha pasado")
	}




	for{



		conec, err := tcp.Accept()
			if err != nil{
		errors.New("lgo ha pasado 2")
	}

		// fmt.Println("Se ha conectado!!")


	 line, _ := request.RequestFromReader(conec)
	//  fmt.Println("esto es la liea", *line)

	//  for {
	// 		mensaje, ok := <-chanMes
	// 			if(!ok){
	// 				fmt.Println("Coneccion terminada")
	// 				return
	// 		}

// 	Request line:
// - Method: GET
// - Target: /
// - Version: 1.1
	mensaje := fmt.Sprintf("Request line: \n - Method: %s\n - Target: %s \n - Version: %s", line.RequestLine.Method, line.RequestLine.RequestTarget, line.RequestLine.HttpVersion )

	menesajeHeaders := "\nHeaders:" 
	for k,v := range line.Headers {
		menesajeHeaders += fmt.Sprintf("\n - %s: %s", k, v)
	}

	
	//  }
	
			fmt.Println(mensaje + menesajeHeaders)



	}

		fmt.Println("termino!")




}
