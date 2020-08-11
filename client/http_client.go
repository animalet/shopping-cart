package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	client   *http.Client
	endpoint string
}

func (c *Client) createCart() *string {
	var response, _ = http.PostForm(c.endpoint, nil)
	if c.isStatusOk(response) {
		bodyBytes := c.readBody(response)
		s := string(bodyBytes)
		return &s
	}
	panic(fmt.Sprintf("Couldn't create the cart: %s", response.Status))
}

func (c *Client) isStatusOk(response *http.Response) bool {
	return response.StatusCode >= 200 && response.StatusCode < 300
}

func (c *Client) removeCart(cartId string) *string {
	request, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", c.endpoint, cartId), nil)
	response, _ := c.client.Do(request)
	if c.isStatusOk(response) {
		bodyBytes := c.readBody(response)
		s := string(bodyBytes)
		return &s
	}

	panic(fmt.Sprintf("Couldn't remove the cart: %s", response.Status))
}

func (c *Client) getTotal(cartId *string) *string {
	var response, _ = http.Get(fmt.Sprintf("%s/%s", c.endpoint, *cartId))
	if c.isStatusOk(response) {
		bodyBytes := c.readBody(response)
		s := string(bodyBytes)
		return &s
	}

	panic(fmt.Sprintf("Couldn't retrieve the cart total price: %s", response.Status))
}

func (c *Client) addProduct(cartId *string, code ProductCode) *string {
	request, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s/%s", c.endpoint, *cartId, code), nil)
	response, _ := c.client.Do(request)
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		bodyBytes := c.readBody(response)
		s := string(bodyBytes)
		return &s
	}

	panic(fmt.Sprintf("Couldn't add product \"%s\" to the cart \"%s\" remove the cart: %s", *cartId, code, response.Status))
}

func (c *Client) readBody(response *http.Response) []byte {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return bodyBytes
}
