package headers


import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	// "fmt"
	
)


func TestHeaders(t *testing.T){
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	// data1 := []byte("Host: locaaaaaaaaat:11111\r\n\r\n")
	// headers.Parse(data1)
	// fmt.Println("esto es header", headers)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	// assert.Equal(t, "localhost:42069", headers["Hosttwo"])

	assert.Equal(t, 23, n)
	assert.False(t, done)
	
	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)



}
// Test: Valid single header