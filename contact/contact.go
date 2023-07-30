package contact

import (
	"fmt"
	"strconv"
)

type Contact struct {
	ID    int
	First string
	Last  string
	Phone string
	Email string
}

var theContacts []Contact = make([]Contact, 0)
var theLastContactID int

func GetContacts() error {
	someContacts := []Contact{
		{1, "joe", "smith", "512-111-2222", "joe.smith@somewhere.com"},
		{2, "mary", "jones", "512-111-3333", "mary.jones@somewhere.com"},
	}
	theContacts = append(theContacts, someContacts...)
	theLastContactID = 2
	return nil
}

func NewContact(contact *Contact) error {
	theLastContactID += 1
	contact.ID = theLastContactID
	theContacts = append(theContacts, *contact)
	return nil
}

func Search(query string) []Contact {
	if query == "" {
		return theContacts
	} else {
		someContacts := make([]Contact, 0)
		for _, aContact := range theContacts {
			if query == aContact.First ||
				query == aContact.Last ||
				query == aContact.Phone ||
				query == aContact.Email {
				someContacts = append(someContacts, aContact)
			}
		}
		return someContacts
	}
}

func GetByID(idString string) (*Contact, error) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}
	for _, theContact := range theContacts {
		if theContact.ID == id {
			return &theContact, nil
		}
	}
	return nil, fmt.Errorf("Contact with ID %d not found", id)
}

func Delete(deletedContact *Contact) error {
	var newTheContacts = make([]Contact, len(theContacts)-1)
	var deleted bool
	for _, theContact := range theContacts {
		var j int
		if theContact.ID == deletedContact.ID {
			deleted = true
		} else {
			newTheContacts[j] = theContact
			j += 1
		}
	}
	if deleted {
		theContacts = newTheContacts
	} else {
		return fmt.Errorf("Contact with ID %d not found", deletedContact.ID)
	}
	return nil
}

func Update(editedContact *Contact) error {
	var edited bool
	for i, theContact := range theContacts {
		if theContact.ID == editedContact.ID {
			theContacts[i] = *editedContact
			edited = true
			break
		}
	}
	if edited {
		return nil
	} else {
		return fmt.Errorf("Contact with ID %d not found", editedContact.ID)
	}
}
