//go:build !embed

package phonedata

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
)

const PHONE_DAT = "phone.dat"

var dftClient *QueryClient

func init() {
	dir := os.Getenv("PHONE_DATA_DIR")
	if dir == "" {
		_, fulleFilename, _, _ := runtime.Caller(0)
		dir = path.Dir(fulleFilename)
	}
	fmt.Println(dir)
	data, err := os.ReadFile(path.Join(dir, PHONE_DAT))
	if err != nil {
		panic(err)
	}

	dftClient, err = New(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
}

func Find(phoneNum string) (pr *PhoneRecord, err error) {
	return dftClient.Find(phoneNum)
}

func Debug() {
	fmt.Println(dftClient.version())
	fmt.Println(dftClient.totalRecord())
	fmt.Println(dftClient.firstRecordOffset())
}
