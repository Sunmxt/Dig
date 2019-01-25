package discover

type Registry interface {
    Service(name string) (Service, error)
    Poll() (bool, error)
    Close()
}


