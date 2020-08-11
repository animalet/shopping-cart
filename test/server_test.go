package test

import (
	"flag"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/go-resty/resty/v2"
	"shopping-cart/server"
	"os"
	"regexp"
	"testing"
)

var opts = godog.Options{Output: colors.Colored(os.Stdout), Format: "progress", Tags: "not @ignore"}
var client = resty.New()
var cartId *string
var lastStatusCode = -1
var started = false

var endpoint = "http://localhost:8000/cart"

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
		cartId = nil
	})

	ctx.Step(`^an existing empty cart$`, anExistingEmptyCart)
	ctx.Step(`^an existing not empty cart$`, anExistingNotEmptyCart)
	ctx.Step(`^the amount (\d+)\.(\d+)€ matches$`, theAmountMatches)
	ctx.Step("^an online cart service$", anOnlineCartService)
	ctx.Step(`^purchasing the following items: ([A-Z ,]+)$`, purchasingTheFollowingItems)
	ctx.Step(`^the cart is removed$`, theCartIsRemoved)
	ctx.Step(`^purchasing any items$`, purchasingAnyItems)
	ctx.Step(`^the amount (\w+) matches$`, theAmountMatches)
	ctx.Step(`^a non existing cart$`, aNonExistingCart)
	ctx.Step(`^cart is not found$`, cartIsNotFound)
}

func cartIsNotFound() error {
	if lastStatusCode != 404 {
		return fmt.Errorf("expected code was 404. Found: %d", lastStatusCode)
	}
	return nil
}

func aNonExistingCart() error {
	s := "UNPROBABLEID"
	cartId = &s
	return nil
}

func anOnlineCartService() error {
	if !started {
		go server.StartServer()
		started = true
	}
	return nil
}

func purchasingAnyItems() error { return purchasingTheFollowingItems("MUG, PEN") }

func purchasingTheFollowingItems(items string) error {
	for _, code := range regexp.MustCompile("[ ,]+").Split(items, -1) {
		if response, err := client.R().Put(fmt.Sprintf("%s/%s/%s", endpoint, *cartId, code)); err == nil && isOkStatus(response) {
			lastStatusCode = response.StatusCode()
		} else {
			return fmt.Errorf("couldn't purchase a product: %s\n%v", err, response)
		}

	}
	return nil
}

func anExistingEmptyCart() error {
	if post, err := client.R().Post(endpoint); err == nil && isOkStatus(post) {
		ptr := post.String()
		fmt.Printf("Got cart id %s\n", ptr)
		cartId = &ptr
	} else {
		return fmt.Errorf("couldn't create a cart: %s \n %v", err, post)
	}
	return nil
}

func theCartIsRemoved() error {
	if response, err := client.R().Delete(fmt.Sprintf("%s/%s", endpoint, *cartId)); err != nil || !isOkStatus(response) {
		return fmt.Errorf("couldn't remove a cart: %s\n%v", err, response)
	}
	return nil
}

func isOkStatus(response *resty.Response) bool {
	statusCode := response.StatusCode()
	fmt.Printf("Got status: %d\n", statusCode)
	return response.StatusCode() > 199 || statusCode < 300
}

func anExistingNotEmptyCart() error {
	if err := anExistingEmptyCart(); err == nil {
		if err := purchasingTheFollowingItems("MUG, TSHIRT"); err != nil {
			return err
		}
		return err
	}
	return nil
}

func theAmountMatches(amount1 int, amount2 int) error {
	if post, err := client.R().Get(endpoint + "/" + *cartId); err == nil && isOkStatus(post) {
		stepAmount := fmt.Sprintf("%d.%02d", amount1, amount2)
		incomingAmount := []rune(post.String())
		s := string(incomingAmount[0 : len(incomingAmount)-1])
		if s != stepAmount {
			return fmt.Errorf("amounts do not match: step=%s,response=%s", stepAmount, s)
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("couldn't get a total: %s\n%v", err, post)
	}
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(
		func() {

		},
	)
}
