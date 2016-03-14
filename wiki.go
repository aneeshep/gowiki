package main

import ( 
//	"fmt"
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
	"errors"
    )

type Page struct {
	Title string
	Body []byte
}

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")


func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
 	return ioutil.WriteFile(filename, p.Body, 0600)
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
		
	}
}

func getTitle( w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w,r)
		return "", errors.New("Invalid Page Title")

	}
	return m[2], nil
}

func loadPage(title string) (*Page,  error){
	filename    := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
        if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html",p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {

	p, err  := loadPage(title)
	if err != nil {
		http.Redirect(w,r,"/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title }
	}

	renderTemplate(w, "edit", p)

}

func  saveHandler(w http.ResponseWriter, r *http.Request, title string) {

	body  := r.FormValue("body")
	p := &Page{ Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w,r,"/view/"+title, http.StatusFound)
}

func main() {
 //    p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page")}
 //    p1.save()
  
 //    p2, _ := loadPage("TestPage")
 //    fmt.Println(string(p2.Body))

	//fs := http.FileServer(http.Dir("tmpl/style"))
	//http.Handle("/style/", fs)

	fs := http.FileServer(http.Dir("tmpl"))
	http.Handle("/style/", fs)

    http.HandleFunc("/view/", MakeHandler(viewHandler))
    http.HandleFunc("/edit/", MakeHandler(editHandler))
    http.HandleFunc("/save/", MakeHandler(saveHandler))
    http.ListenAndServe(":8080",nil)

}


