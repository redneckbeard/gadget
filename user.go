package gadget

// User is the interface that wraps Gadget's pluggable authentication
// mechanism. Any type that provides an Authenticated method satisfies it.
type User interface {
	Authenticated() bool
}

// AnonymousUser is the default User type if the UserIdentifier registered by
// the application returns nil or if there is no UserIdentifier.
type AnonymousUser struct {
}

// Authenticated for AnonymousUser will always return false.
func (u *AnonymousUser) Authenticated() bool { return false }

// UserIdentifier is a function type that returns a User type or nil based on a
// Request.
type UserIdentifier func(*Request) User

var identifyUser UserIdentifier

// IdentifyUsersWith allows Gadget applications to register a UserIdentifier
// function to be called when processing each incoming request. The return
// value of the UserIdentifier will be set on the Request object if not nil;
// AnonymousUser will be used otherwise.
func IdentifyUsersWith(ui UserIdentifier) {
	identifyUser = ui
}

func clearUserIdentifier() {
	identifyUser = nil
}
