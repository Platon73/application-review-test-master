package interfaces

type EmailSender interface {
	SendConfirmation(email, content string) error
}
