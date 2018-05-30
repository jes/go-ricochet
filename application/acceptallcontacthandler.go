package application

type AcceptAllContactHandler struct{}

func (aach *AcceptAllContactHandler) ContactRequest(name string, message string) string {
	return "Pending"
}
func (aach *AcceptAllContactHandler) ContactRequestRejected() {
}
func (aach *AcceptAllContactHandler) ContactRequestAccepted() {
}
func (aach *AcceptAllContactHandler) ContactRequestError() {
}
