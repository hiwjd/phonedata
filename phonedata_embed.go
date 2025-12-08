//go:build embed

package phonedata

import (
	"bytes"
	_ "embed"
	"fmt"
)

//go:embed phone.dat
var data []byte

var embedClient *QueryClient

func init() {
	var err error
	embedClient, err = New(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
}

func Find(phoneNum string) (pr *PhoneRecord, err error) {
	return embedClient.Find(phoneNum)
}

func Debug() {
	fmt.Println(embedClient.version())
	fmt.Println(embedClient.totalRecord())
	fmt.Println(embedClient.firstRecordOffset())
}
