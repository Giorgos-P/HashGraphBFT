package variables

var (
	// N - Number of processors
	N int

	// F - Number of faulty processors
	F int

	// ID - This processor's id.
	ID int

	// View - This processor's view.
	View int

	// T - Threshold (??)
	T int

	// Clients - Size of Clients Set
	Clients int

	// Remote - If we are running locally or remotely
	Remote bool

	Stop bool
)

// Initialize - Variables initializer method
func Initialize(id int, n int, c int, rem int) {
	ID = id
	N = n
	F = (N - 1) / 3
	T = N - F
	//T = int(math.Ceil((2.0 / 3.0) * float64(N)))
	//F = N - T
	Clients = c

	if rem == 0 {
		Remote = false
	} else {
		Remote = true
	}
	Stop = false

	//fmt.Printf("Threshold=%d\n", T)
	//fmt.Printf("faulty=%d\n", F)

}
