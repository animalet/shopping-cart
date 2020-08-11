package domain

type (
	Cart struct {
		Items map[string]int //keys are item codes (MUG, PEN, TSHIRT). values are the amount of each item in the cart
		Id    *string        //cart identification (should be a UUID)
	}
)
