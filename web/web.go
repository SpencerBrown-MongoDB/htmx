package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/SpencerBrown-MongoDB/htmx/contact"
)

const port = "8080"

func Setup() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "static", "favicon.ico"))
	})
	http.HandleFunc("/assets/htmx.min.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "assets", "htmx.min.jx"))
	})
	http.HandleFunc("/contacts/view/", viewHandler)
	http.HandleFunc("/contacts/edit/", editHandler)
	http.HandleFunc("/contacts/delete/", deleteHandler)
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
	runTemplate(w, "contacts", &templateData{
		Contacts: someContacts,
		Query:    queryString,
	})
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/contacts/view/")
	if idString == "" {
		userError(w, "View requires an id like /contacts/view/ID_NUMBER")
		return
	}
	theContact, err := contact.GetByID(idString)
	if err != nil {
		userError(w, fmt.Sprintf("contact ID error: %v", err))
		return
	}
	runTemplate(w, "viewContact", theContact)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		idString := strings.TrimPrefix(r.URL.Path, "/contacts/edit/")
		if idString == "" {
			userError(w, "Edit requires an ID like /contacts/edit/CONTACT_ID")
			return
		}
		theContact, err := contact.GetByID(idString)
		if err != nil {
			userError(w, fmt.Sprintf("Error with contact ID: %v", err))
			return
		}
		runTemplate(w, "editContact", theContact)
	case "POST":
		idString := strings.TrimPrefix(r.URL.Path, "/contacts/edit/")
		if idString == "" {
			userError(w, "Edit requires an ID like /contacts/edit/CONTACT_ID")
			return
		}
		theContact, err := contact.GetByID(idString)
		if err != nil {
			userError(w, fmt.Sprintf("Error with contact ID: %v", err))
			return
		}
		theContact.First = r.FormValue("first_name")
		theContact.Last = r.FormValue("last_name")
		theContact.Phone = r.FormValue("phone")
		theContact.Email = r.FormValue("email")
		err = contact.Update(theContact)
		if err == nil {
			runTemplate(w, "contactEdited", theContact)
		} else {
			userError(w, fmt.Sprintf("Error: %v", err))
		}
	default:
		userError(w, "Must use GET or POST for edit contact")
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/contacts/delete/")
	if idString == "" {
		userError(w, "Delete requires a contact ID like /contacts/delete/CONTACT_ID")
		return
	}
	theContact, err := contact.GetByID(idString)
	if err != nil {
		userError(w, fmt.Sprintf("Error in contact ID: %v", err))
		return
	}
	err = contact.Delete(theContact)
	if err == nil {
		runTemplate(w, "contactDeleted", theContact)
	} else {
		userError(w, fmt.Sprintf("Error deleting contact: %v", err))
	}
}

func newContactHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// user requesting new contact
		runTemplate(w, "newContact", &contact.Contact{
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
			runTemplate(w, "contactCreated", &newContact)
		} else {
			userError(w, fmt.Sprintf("Error creating new contact: %v", err))
		}
	default:
		userError(w, "Must use GET or POST for new contact")
	}
}

func runTemplate(w io.Writer, templateFile string, data any) {
	htmlString := readEmbed(w, templateFile)
	theTemplate, err := template.New("theTemplate").Parse(htmlString)
	webError(w, "Internal error parsing template", err)
	err = theTemplate.Execute(w, data)
	webError(w, "Internal error executing template", err)
}

func readEmbed(w io.Writer, templateFile string) string {
	starter, err := asset.ReadFile(filepath.Join("template", "start.html"))
	if err != nil {
		webError(w, "Internal error reading start.html", err)
	}
	ender, err := asset.ReadFile(filepath.Join("template", "end.html"))
	if err != nil {
		webError(w, "Internal error reading end.html", err)
	}
	content, err := asset.ReadFile(filepath.Join("template", templateFile+".html"))
	if err != nil {
		webError(w, "Internal error reading "+templateFile, err)
	}
	full := append(starter, content...)
	full = append(full, ender...)
	return string(full)
}

func userError(w io.Writer, e string) {
	runTemplate(w, "error", e)
}

func webError(w io.Writer, e string, err error) {
	if err == nil {
		return
	}
	errString := fmt.Sprintf(e+": %v", err)
	_, _ = io.WriteString(w, errString)
	log.Fatalf(errString)
}
