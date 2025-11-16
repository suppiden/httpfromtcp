package headers

import(
	"errors"
	// "strings"
	"regexp"
	"fmt"
)

func validateFields(s string) bool {
	// re := regexp.MustCompile("([A-Za-z!#$%&'*+-.^_`|~]{1,})\\d")
	re := regexp.MustCompile("^[A-Za-z0-9!#$%&'*+-.^_`|~]")


	return re.MatchString(s)
}

type Headers map[string]string

func (h *Headers) Parse(data []byte) (n int, done bool, err error) {

	var head string
	var pair string
	// var acabado bool
	var consumedHead int
	var consumedPair int
	var encontrado bool
	// var headerFinal Headers

	if string(data[0]) == "\r" && string(data[1]) == "\n" {
		return 0, true, nil
	}

	for i := 0; len(data) > i; i++ {
		if string(data[i]) == "\r" && string(data[i+1]) == "\n" {
			encontrado = true
		}
	}

	if !encontrado {
		return 0, false, nil
	}

	for i := 0; len(data) > i; i++ {

		// fmt.Println("a ver cuando sale out ", i, "y esto es el len ", len(data))
		if string(data[i]) != " " {
			if string(data[i]) == ":" {
				if string(data[i+2]) == " " {
					return 0, false, errors.New("Mal formayo")
				}
				break
			}
			consumedHead++
			// fmt.Println("esto es la letra ", string(data[i]), " y esto la suma", consumedHead)
			head += string(data[i])

		}
	}

	// fmt.Printf("esta es el head %s", head)

	data = data[consumedHead+1:]

	if string(data[1]) == " " {
		return 0, false, errors.New("Mal formayo")

	}

	for i := 0; len(data) > i; i++ {

		if string(data[i]) != " " {
			if string(data[i]) == "\r" && string(data[i+1]) == "\n" {
				break
			}
			consumedPair++
			// fmt.Println("esto es la letra ", string(data[i]), " y esto la suma", consumedPair)

			pair += string(data[i])

		}
	}

	// fmt.Printf(" y este es el par %s", pair)

	// headerFinal := Headers{
	// 	head: pair
	// }
	isHead := validateFields(head)
	if !isHead {
		fmt.Println("esto es el head", head)
		return 0, false, errors.New("invalid characters head")
	}


	isPair := validateFields(pair)
	if !isPair {
		return 0, false, errors.New("invalid characters body")
	}

	
	// head = strings.ToLower(head)
	// if strings.Contains(pair, ":") {
	// 	index := strings.Index(pair, ":")
	// 	pair = strings.ToLower(pair[:index]) + pair[index:]
	// 	}else{
	// 		pair = strings.ToLower(pair)
	// 	}
	// fmt.Println("esto es el pair ", pair
	
	
	if _ , ok := (*h)[head]; !ok {
		headerFinal := Headers{
			head: pair,
		}
	
		// fmt.Println("\n a ver que hay", h)
	
		*h = headerFinal

	}else{
		resultado := (*h)[head]
		(*h)[head] = resultado + " " + pair
	}

	fmt.Println("\n a ver que hay", h)

	return consumedHead + consumedPair +4, false, nil
}

func NewHeaders() Headers{
	var header Headers 

	return header
}