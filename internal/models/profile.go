package models

type Profile struct {
	*User
	UserOrders []Order
}

func NewProfile() *Profile {
	return &Profile{
		User:       &User{},
		UserOrders: make([]Order, 0, 1000),
	}
}
