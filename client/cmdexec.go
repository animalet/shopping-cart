package client

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

import (
	"os"
	"strconv"
)

type ProductCode string

const (
	PEN    ProductCode = "PEN"
	TSHIRT ProductCode = "TSHIRT"
	MUG    ProductCode = "MUG"
)

// Manages the client logic
type commandExecutor struct {
	selectionRegexp *regexp.Regexp
	client          *Client
	cartId        *string
	codes           []ProductCode
	purchases       []string
}

type CommandExecutor interface {
	StringCodesForConsole() *string
	Execute(text string)
}

func NewCommandExecutor(serverAddress string) *commandExecutor {
	// Regex accepts integers separated by blank spaces, "q", "r" and remaining blanks (to be trimmed)
	validationRegexp := regexp.MustCompile(`^\s*(([0-9]|\s*)+|[qr])\s*$`)
	return &commandExecutor{
		selectionRegexp: validationRegexp,
		client: &Client{
			client:   &http.Client{},
			endpoint: fmt.Sprintf("%s/%s", serverAddress, "cart"),
		},
		codes: []ProductCode{PEN, TSHIRT, MUG},
	}
}

// Convenience method to provide all the product codes to show them in the console
func (cart *commandExecutor) StringCodesForConsole() *string {
	codeStrings := make([]string, 0)
	for i, code := range cart.codes {
		codeStrings = append(codeStrings, fmt.Sprintf("%d: %s", i+1, code))
	}
	parsed := strings.Join(codeStrings, ", ")
	return &parsed
}

// Parses the input
func (cart *commandExecutor) Execute(text string) {
	if selection := strings.ToLower(strings.TrimSpace(text)); cart.selectionRegexp.MatchString(selection) {
		switch selection {
		case "q":
			fmt.Println("Bye!")
			os.Exit(1)
		case "r":
			if cart.cartId == nil {
				fmt.Println("No cart is active")
			} else {
				cart.removeCart()
			}
		default:
			cart.addProductToCart(selection)
			cart.getCartTotal()
		}
	} else {
		fmt.Printf("Invalid input: \"%s\"\n", selection)
	}
}

// Remove cart and clear cart id
func (cart *commandExecutor) removeCart() {
	defer cart.reset()
	cart.client.removeCart(*cart.cartId)
	fmt.Printf("Cart %s removed\n", *cart.cartId)
}

// Adds a products to a cart from the command string. If there is no cart it calls the creation service and stores
// the id for further operations
func (cart *commandExecutor) addProductToCart(commandString string) {
	if cart.cartId == nil {
		cart.cartId = cart.client.createCart()
		fmt.Printf("Cart id is %s\n", *cart.cartId)
	}
	for _, s := range strings.Split(commandString, " ") {
		if i, err := strconv.Atoi(s); err == nil {
			code := cart.codes[i-1]
			cart.client.addProduct(cart.cartId, code)
			productCode := code
			cart.purchases = append(cart.purchases, string(productCode))
		}
	}
}

func (cart *commandExecutor) getCartTotal() {
	total := cart.client.getTotal(cart.cartId)
	fmt.Println()
	fmt.Println(strings.Join(cart.purchases, ", "))
	fmt.Println(*total)
}

func (cart *commandExecutor) reset() {
	cart.cartId = nil
	cart.purchases = make([]string, 0)
}
