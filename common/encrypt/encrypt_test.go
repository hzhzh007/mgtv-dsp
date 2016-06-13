package encrypt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TEST_KEY = "e860399e0a424905"
)

func Test_EncryData(t *testing.T) {
	assert := assert.New(t)
	str := "happy every day"
	encrypted := AesCBCEncrypte([]byte(str), TEST_KEY)
	decrypted, err := AesCBCDecrypte(encrypted, TEST_KEY)
	assert.Nil(err)
	assert.Equal(str, decrypted)
}

func BenchmarkEncryptePrice(b *testing.B) {
	str := "2481"
	for i := 0; i < b.N; i++ {
		AesCBCEncrypte([]byte(str), TEST_KEY)
	}
}

func BenchmarkDecryptePrice(b *testing.B) {
	str := "2481"
	encrypted := AesCBCEncrypte([]byte(str), TEST_KEY)
	for i := 0; i < b.N; i++ {
		AesCBCDecrypte(encrypted, TEST_KEY)
	}
}
func Test_Decryqq(t *testing.T) {
	decrypted, err := AesCBCDecrypte("560330bd3aebfc34bb0715af3a6212c4", "975dfad4e0b94c38")
	t.Log("enc")
	t.Log(AesCBCEncrypte([]byte("456"), TEST_KEY))
	t.Log(decrypted, err)
	fmt.Println("test  decrypted:", decrypted)
}
