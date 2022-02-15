package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/sessions"
	"github.com/mateors/mcb"
)

var store = sessions.NewCookieStore([]byte("super-secret"))

type UserDetails struct {
	FullName string `json:"fullname"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var db *mcb.DB
var tpl *template.Template

func init() {

	db = mcb.Connect("130.185.118.116", "sazid", "SaZ!d2022", false)

	res, err := db.Ping()
	if err != nil {

		fmt.Println(res)
		os.Exit(1)
	}
	fmt.Println(res, err)

}

func index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	user_name, ok := session.Values["userID"]
	fmt.Println("ok:", ok)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
		return
	}
	tpl.ExecuteTemplate(w, "index.html", user_name)
}

func login(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.html", nil)

}

func loginAuth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userName := r.FormValue("username")
	password := r.FormValue("password") //map[password:ross123]
	sql := fmt.Sprintf("SELECT * FROM chaldal_erp.sazid.signup_data WHERE `username` = '%s' AND  `password`='%s'", userName, password)
	res := db.Query(sql)

	if len(res.Result) == 0 {
		fmt.Println("Check Username And Password")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	} else {

		session, _ := store.Get(r, "session")

		session.Values["userID"] = userName
		session.Save(r, w)
		str := fmt.Sprintf("Hello %s!! Welcome Here!!", userName)
		tpl.ExecuteTemplate(w, "index.html", str)
		return

	}

}
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****logoutHandler running*****")
	session, _ := store.Get(r, "session")
	delete(session.Values, "userID")
	session.Save(r, w)
	tpl.ExecuteTemplate(w, "login.html", "Logged Out")
}
func registration(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "registration.html", nil)

}

func registrationAuth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fullname := r.FormValue("fullname")
	username := r.FormValue("username")

	// record := userDetails{}
	// record.fullName = fullname
	// record.userName = username
	// record.email = email
	// record.password = password

	var myData UserDetails

	form := make(url.Values, 0)
	form.Add("bucket", "chaldal_erp.sazid.signup_data") //bucket and collection-> namespace:bucket.scope.collection
	//document ID
	form.Add("aid", r.FormValue("username"))
	form.Add("fullname", r.FormValue("fullname"))
	form.Add("username", r.FormValue("username"))
	form.Add("email", r.FormValue("email"))
	form.Add("password", r.FormValue("password"))
	p := db.Insert(form, &myData) //pass by reference (&myData)
	fmt.Println("Status:", p.Errors)

	session, _ := store.Get(r, "session")

	session.Values["userID"] = username
	session.Save(r, w)
	str := fmt.Sprintf("Thanks %s, for your registration!", fullname)
	tpl.ExecuteTemplate(w, "index.html", str)

}

func main() {
	tpl, _ = template.ParseGlob("*.html")

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/registration", registration)
	http.HandleFunc("/registrationauth", registrationAuth)
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/loginAuth", loginAuth)
	http.HandleFunc("/logout", logoutHandler)

	fmt.Printf("Starting server got testing\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
