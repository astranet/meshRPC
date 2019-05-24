package greeter

import (
	"errors"
	"log"
)

type Postcard struct {
	PictureURL string
	Address    string
	Recipient  string
	Message    string
}

//go:generate meshRPC expose -P greeter -y

type Service interface {
	Greet(name string) (message string, err error)
	SendPostcard(card *Postcard) (err error)
}

func NewService() Service {
	return &service{}
}

type service struct{}

func (s *service) Greet(name string) (string, error) {
	message := "Hello, " + name
	return message, nil
}

func (s *service) SendPostcard(card *Postcard) error {
	if card == nil {
		return errors.New("no postcard")
	} else if card.Recipient == "" {
		return errors.New("no recipient")
	}
	log.Printf("sending %#v to %s", *card, card.Recipient)
	return nil
}
