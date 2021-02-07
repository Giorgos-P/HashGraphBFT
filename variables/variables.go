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
)

// Initialize - Variables initializer method
func Initialize(id int, n int, t int, c int) {
	ID = id
	N = n
	F = (N - 1) / 3
	T = t
	Clients = c
	Remote = false
}
