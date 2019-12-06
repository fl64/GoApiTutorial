// main.go

package main

func main() {
	a := App{}
	a.Initialize("root", "password", "gm_licenses")

	a.Run(":8080")
}

