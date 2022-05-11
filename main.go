package main

func main() {
	starter := getStarter()
	err := starter()
	if err != nil {
		panic(err)
	}
}
