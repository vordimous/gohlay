package internal

// Run will scan a configured topic and deliver messages
func Run() {
	CheckForDeliveries()
	HandleDeliveries()
}
