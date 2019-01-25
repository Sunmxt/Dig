package discover

type Service interface {
    Endpoints() []string
    Metadata(endpoint string) map[string]string
    Watch() error
    Publish(node *Node) error
}
