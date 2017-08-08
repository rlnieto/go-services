package password

import(
  "encoding/base64"
  "math/rand"
  "crypto/sha256"
  "time"
)

const randomLength = 16

/*------------------------------------------------------------------------------
 Genera un string con caracteres aleatorios entre los códigos ASCII 32 y 126

------------------------------------------------------------------------------*/
func GenerateSalt(length int) string{
  var salt []byte
  var asciiPad int64

  if length == 0{
    length = randomLength
  }

  asciiPad = 32

  for i:=0; i<length; i++ {
    salt = append(salt, byte(rand.Int63n(94) + asciiPad))
  }

  return string(salt)
}

/*------------------------------------------------------------------------------
 Calcula el hash de password + salt

------------------------------------------------------------------------------*/
func GenerateHash(salt string, password string) string{
  var hash string
  fullString := salt + password

  sha := sha256.New()
  sha.Write([]byte(fullString))

  hash = base64.URLEncoding.EncodeToString(sha.Sum(nil))

  return hash
}

/*------------------------------------------------------------------------------
 Calcula salt y hash para una password

------------------------------------------------------------------------------*/
func ReturnPassword(password string) (string, string){
  rand.Seed(time.Now().UTC().UnixNano())

  salt := GenerateSalt(0)
  hash := GenerateHash(salt, password)

  return salt, hash
}


func GenerateSessionID(length int) string {
	var salt []byte
	var asciiPad int64

	if length == 0 {
		length = randomLength
	}

	asciiPad = 60

	for i := 0; i < length; i++ {
		salt = append(salt, byte(rand.Int63n(62)+asciiPad))
	}

	return string(salt)
}
