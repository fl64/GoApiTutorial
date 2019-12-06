// main.go

package main

func main() {
	a := App{}
	a.Initialize("root", "password", "rest_api_example")

	a.Run(":8080")
}
