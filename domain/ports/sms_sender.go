package ports

type SmsSender interface {
	SendToAuthor(to string, content string) error
}
