package request

import (
	"errors"
	//"fmt"
	"io"
	"strings"
)

type state int

const (
	initialized state = iota
	done
	requestStateParsingHeaders
)

const bufferSize int = 8

type Request struct {
	RequestLine RequestLine
	state       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}



func (r *Request) parse(data []byte) (int, error) {
	// if r.state == 0 {
	// 	return 0, nil
	// }

	// fmt.Println("dentro e perte", string(data),"    ", len(data))

	if r.state == 0 { 
		parseado, err := parseRequestLine(data)
		// fmt.Println("esto es en parse stringfinal1", string(data[:parseado]))
		if err != nil {cd
			return 0, err
		}

		if parseado == 0 && err == nil {
		// fmt.Println("esto es en parse stringfinal2")

			return 0, nil
		}

		
		// fmt.Println("esto es en parse stringfinal3")
		if len(data) == parseado {
			
			r.state = 2
			
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
		numeroValido := strings.Trim(numero[1], "\r")


		if numeroValido!= "1.1" {
			return 0, errors.New("Problema con la version")
		}
		
		partes[2] = numeroValido
		version := partes[2]

		// fmt.Println("esto es desde alla3")
		r.RequestLine = RequestLine{
			HttpVersion:   version,
			RequestTarget: direccion,
			Method:        metodo,
		}
		r.state = 2

		return parseado, nil

	} else if r.state == 1 {
		return 0, errors.New("error: trying to read data in a done state")
	} else if r.state == 2 {

	headers = NewHeaders()
	n, done, err = headers.Parse(data[parseado:])
	if err != nil{
		return 0, errors.New("error: unknown state")

	}

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
				return numeroProcesado, nil
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
			// fmt.Println("Una vez ", len(buf), "   ", string(buf), "  ", readToIndex)
			consumed, err1 := io.ReadFull(reader, buf[readToIndex:])
			r.parse(buf)
			// fmt.Println("Una SEGUNDA VEZ consumed ", consumed, "  parsed " , parsed, " indez ", readToIndex, " la string", string(buf), "lo  longitud", len(buf))
			readToIndex =  consumed + readToIndex
			// fmt.Println("Una TERECRa VEZ consumed "," indez ", readToIndex)
			
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
				// fmt.Println("BEEEEEEEEEEEEEEEEEEe", parsed,"   ", readToIndex, "   ", consumed)
				
				if consumed == len(buf){
				// // fmt.Println("esto es el segundo parse pero primero", "  ", string(buf),"   ")


				// fmt.Println("esto es el segundo parse", string(newBud), "  ", string(buf),"   ", len(newBud))
				// readToIndex =0

				}
			}
			// fmt.Println("esto va mal", err, r)

			
			if err1 != nil{
			// // fmt.Println("acabó",readToIndex, err1, r)
			r.state = 1
			break
				
			}

		}
	
	// fmt.Println(" a ver que devuelve el struct", r)
	return &r, nil
}