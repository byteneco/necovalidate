package necovalidate

// @Name
// @FF
func validate() {

}

type Ele interface {
	int
}

type Customer[T Ele, L Ele] struct {
	// @Email
	// @Max 10
	Name string
	Age  int
}
