package main

/*lol comment*/
import (
	"fmt"
	"te4-introprojekt/metastore"
)

func main() {
	fmt.Println("Hello World!")
	var i = metastore.CreateInformation()
	fmt.Printf(i.GetCreated().String())
}
