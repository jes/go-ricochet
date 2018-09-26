package application

type AcceptAllContactHandler struct {
	Rai *ApplicationInstance
}

func (aach *AcceptAllContactHandler) ContactRequest(hostname string, name string, message string) string {
	return "Pending"
}
func (aach *AcceptAllContactHandler) ContactRequestRejected() {
}
func (aach *AcceptAllContactHandler) ContactRequestAccepted() {
}
func (aach *AcceptAllContactHandler) ContactRequestError() {
}
