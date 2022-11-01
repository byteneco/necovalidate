package necovalidate

// @Name
// @FF
func validate() {

}

type Customer struct {
	// @Email
	// @Max 10
	// @Min 10
	// @NotEmpty
	Name string
	Age  int
}

type F struct {
	// @Email
	Email string
}
