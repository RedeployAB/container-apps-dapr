package server

// routes setups registers routes and handlers for the server.
func (s server) routes() {
	s.router.Handle("/reports", authenticate(s.security.Keys, s.reportHandler()))
}
