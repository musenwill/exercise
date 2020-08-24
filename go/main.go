package main

func main() {
	foo := make(map[string]struct{})
	bar := make(map[string]interface{})

	foo["foo"] = struct{}{}
	bar["bar"] = struct{}{}

}
