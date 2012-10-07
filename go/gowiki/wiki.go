package main

import (
  	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
	"errors"
	"path/filepath"
	"os"
)

const lenPath=len("/view/")
var templates= template.Must(template.ParseFiles("edit.html","view.html"))
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")
var TEMPLATE_PATH="tmpl/"
var DATA_PATH="data/"

type Page struct{
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title +".txt"
	return ioutil.WriteFile(filename, p.Body,0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string ){
 
 
	p, err :=loadPage(title)
	if err!=nil{
		http.Redirect(w,r,"/edit/"+title,http.StatusFound)
		return
	}
  	renderTemplate(w,"view",p)
}

func editHandler(w http.ResponseWriter, r *http.Request,title string){
 
 	p, err := loadPage(title)
	if err !=nil {
		p =&Page{Title: title}
	}
  	renderTemplate(w,"edit",p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	err  :=templates.ExecuteTemplate(w,tmpl+".html",p)
	if err !=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	 
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string ){
  	body  := r.FormValue("body")
	p := &Page{Title :title,Body :[]byte(body)}
	err:= p.save()
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r,"/view/"+title,http.StatusFound)
}

func getTitle(w http.ResponseWriter, r *http.Request) (title string, err error) {
    title = r.URL.Path[lenPath:]
    if !titleValidator.MatchString(title) {
        http.NotFound(w, r)
        err = errors.New("Invalid Page Title")
    }
    return
}
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        title := r.URL.Path[lenPath:]
        if !titleValidator.MatchString(title) {
            http.NotFound(w, r)
            return
        }
        fn(w, r, title)
    }
}
 
 

func main() {
	walk :=func(path string, info os.FileInfo, err error)error{
		if err==nil && !info.IsDir() && filepath.Ext(path)== ".html"{
				http.HandleFunc("/view/",makeHandler(viewHandler))
		 	    http.HandleFunc("/edit/", makeHandler(editHandler))
				http.HandleFunc("/save/",makeHandler(saveHandler))
		}
		return err
	}
	root :=TEMPLATE_PATH+"edit.html"
	err :=filepath.Walk(root,walk)
	if err !=nil {
		return
	}
 

		http.ListenAndServe(":8080",nil)
 }