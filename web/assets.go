package web

// Static content as string variables

import _ "embed"

//go:embed assets/start.html
var start_html string

//go:embed assets/end.html
var end_html string

//go:embed assets/contacts.html
var contacts_html string

//go:embed assets/newContact.html
var newContact_html string

//go:embed assets/contactCreated.html
var contactCreated_html string

//go:embed assets/viewContact.html
var viewContact_html string

//go:embed assets/editContact.html
var editContact_html string

//go:embed assets/contactEdited.html
var contactEdited_html string

//go:embed assets/contactDeleted.html
var contactDeleted_html string

//go:embed assets/error.html
var error_html string
