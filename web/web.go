package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/SpencerBrown-MongoDB/htmx/contact"
)

const port = "8080"

func Setup() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "assets", "favicon.ico"))
	})
	http.HandleFunc("/contacts/view", viewHandler)
	http.HandleFunc("/contacts/edit", editHandler)
	http.HandleFunc("/contacts/delete", deleteHandler)
	http.HandleFunc("/contacts", contactHandler)
	http.HandleFunc("/contacts/new", newContactHandler)
	fmt.Printf("Starting web server on port %s: use Ctrl+C to stop\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusFound)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query().Get("q")
	someContacts := contact.Search(queryString)
	type templateData = struct {
		Contacts []contact.Contact
		Query    string
	}
	runTemplate(w, &contacts_html, &templateData{
		Contacts: someContacts,
		Query:    queryString,
	})
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query().Get("id")
	if queryString == "" {
		userError(w, "View requires an id query like /contacts/view?id=NUMBER")
		return
	}
	theContact, err := contact.GetByID(queryString)
	if err != nil {
		userError(w, fmt.Sprintf("error: %v", err))
		return
	}
	runTemplate(w, &viewContact_html, theContact)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		queryString := r.URL.Query().Get("id")
		if queryString == "" {
			userError(w, "Edit requires an id query like /contacts/view?id=NUMBER")
			return
		}
		theContact, err := contact.GetByID(queryString)
		if err != nil {
			userError(w, fmt.Sprintf("error: %v", err))
			return
		}
		runTemplate(w, &editContact_html, theContact)
	case "POST":
		idString := r.URL.Query().Get("id")
		if idString == "" {
			userError(w, "Edit requires an id query like /contacts/edit?id=NUMBER")
			return
		}
		theContact, err := contact.GetByID(idString)
		if err != nil {
			userError(w, fmt.Sprintf("error: %v", err))
			return
		}
		theContact.First = r.FormValue("first_name")
		theContact.Last = r.FormValue("last_name")
		theContact.Phone = r.FormValue("phone")
		theContact.Email = r.FormValue("email")
		err = contact.Update(theContact)
		if err == nil {
			runTemplate(w, &contactEdited_html, theContact)
		} else {
			userError(w, fmt.Sprintf("Error: %v", err))
		}
	default:
		userError(w, "Must use GET or POST for edit contact")
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	if idString == "" {
		userError(w, "Delete requires an id query like /contacts/delete?id=NUMBER")
		return
	}
	theContact, err := contact.GetByID(idString)
	if err != nil {
		userError(w, fmt.Sprintf("Error: %v", err))
		return
	}
	err = contact.Delete(theContact)
	if err == nil {
		runTemplate(w, &contactDeleted_html, theContact)
	} else {
		userError(w, fmt.Sprintf("Error: %v", err))
	}
}

func newContactHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// user requesting new contact
		runTemplate(w, &newContact_html, &contact.Contact{
			Email: r.URL.Query().Get("email"),
			First: r.URL.Query().Get("first"),
			Last:  r.URL.Query().Get("last"),
			Phone: r.URL.Query().Get("phone"),
		})
	case "POST":
		// user saving new contact
		var newContact = contact.Contact{
			Email: r.FormValue("email"),
			First: r.FormValue("first_name"),
			Last:  r.FormValue("last_name"),
			Phone: r.FormValue("phone"),
		}
		err := contact.NewContact(&newContact)
		if err == nil {
			runTemplate(w, &contactCreated_html, &newContact)
		} else {
			userError(w, fmt.Sprintf("Error creating new contact: %v", err))
		}
	default:
		userError(w, "Must use GET or POST for new contact")
	}
}

func runTemplate(w io.Writer, templateString *string, data any) {
	htmlString := start_html + *templateString + end_html
	theTemplate, err := template.New("theTemplate").Parse(htmlString)
	webError(w, "Internal error parsing template", err)
	err = theTemplate.Execute(w, data)
	webError(w, "Internal error executing template", err)
}

func userError(w io.Writer, e string) {
	runTemplate(w, &error_html, e)
}

func webError(w io.Writer, e string, err error) {
	if err == nil {
		return
	}
	errString := fmt.Sprintf(e+": %v", err)
	_, _ = io.WriteString(w, errString)
	log.Fatalf(errString)
}
