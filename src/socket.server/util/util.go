package util

import(
	"fmt"
	"os"
)

func CheckError(err error){
	if err != nil {
		fmt.Println("Error = ",err)
		os.Exit(1)
	}
}

