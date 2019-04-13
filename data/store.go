package data

import (
	"encoding/json"
	"fmt"
	"github/kerrai1990/phonebook-rest-api/models"
	"strconv"

	"github.com/go-redis/redis"
)

//Client -
type Client struct {
	client *redis.Client
}

// New -
func New() *Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	return &Client{client: client}
}

//CheckStatus -
func (c *Client) CheckStatus() {
	ping := c.client.Ping()
	fmt.Println(ping)
}

//Index -
func (c *Client) Index() {
}

//Get -
func (c *Client) Get(id int) (*models.Person, error) {

	p := c.client.Get("people-" + strconv.Itoa(id))

	b, err := p.Bytes()
	if err != nil {
		return nil, err
	}

	var person models.Person
	if err := json.Unmarshal(b, &person); err != nil {
		return nil, err
	}

	return &person, nil

}

//Create -
func (c *Client) Create(person models.Person) (*models.Person, error) {

	id := strconv.Itoa(person.ID)

	obj, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}

	err = c.client.Set(fmt.Sprintf("people-%s", id), obj, 0).Err()
	if err != nil {
		panic(err)
	}

	p := c.client.Get("people-" + strconv.Itoa(person.ID))

	b, err := p.Bytes()
	if err != nil {
		return nil, err
	}

	var newPerson models.Person
	if err := json.Unmarshal(b, &newPerson); err != nil {
		return nil, err
	}

	return &newPerson, nil
}

//Update -
func (c *Client) Update(person models.Person) (*models.Person, error) {

	id := strconv.Itoa(person.ID)

	p := c.client.Get("people-" + id)

	b, err := p.Bytes()
	if err != nil {
		return nil, err
	}

	var result models.Person
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}

	return &result, nil

}

//Delete -
func (c *Client) Delete() {

	fmt.Println("DELETE ME")
}
