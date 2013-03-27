package gadget

type User interface {
	Authenticated() bool
}

type AnonymousUser struct {
}

func (u *AnonymousUser) Authenticated() bool { return false }

type UserIdentifier func(*Request) User

var identifyUser UserIdentifier

func IdentifyUsersWith(ui UserIdentifier) {
	identifyUser = ui
}

func clearUserIdentifier() {
	identifyUser = nil
}
