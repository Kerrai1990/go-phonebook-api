package models

import "errors"

//Person -
type Person struct {
	ID        int      `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

//ValidatePerson -
func (p Person) ValidatePerson() error {
	if p.Firstname == "" || p.Lastname == "" {
		return errors.New("Firstname, and Lastname must be defined")
	}
	return nil
}
