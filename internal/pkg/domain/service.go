package domain

type Service struct {
	User	int `json:"user,omitempty"`
	Forum  	int `json:"forum,omitempty"`
	Thread 	int	`json:"thread,omitempty"`
	Post 	int `json:"post,omitempty"`
}

type ServiceRepository interface {
	Clear()  error
	Status() (Service, error)
}

type ServiceUsecase interface {
	Clear()  error
	Status() (Service, error)
}
