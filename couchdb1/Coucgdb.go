package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/leesper/couchdb-golang"
)

type Comp struct {
	Cname string          `json:"cname"`
	Kra   map[string]Kras `json:"kra"`
}
type Kras struct {
	Title string            `json:"title"`
	Kpi   map[string]string `json:"kpi"`
}

var comp []Comp
var cm Comp

func Company(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(comp)
}
func NewCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Comp
	vars := mux.Vars(r)
	db, err := couchdb.NewDatabase("http://prasanna:prasanna@127.0.0.1:5984/company")
	if err != nil {
		panic(err)
	}
	json.NewDecoder(r.Body).Decode(&book)
	//verify if company is in db
	er := db.Contains(vars["id"])
	if er != nil {
		km := make(map[string]interface{})
		km["company"] = book
		err = db.Set(vars["id"], km)
		if err != nil {
			panic(err)
		}
		comp = append(comp, book)
		json.NewEncoder(w).Encode(book)
	} else {
		fmt.Fprintf(w, "Doc already exists")
	}
}

//Gets kra using company id
func Getkra(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	db, err := couchdb.NewDatabase("http://prasanna:prasanna@127.0.0.1:5984/company")
	if err != nil {
		panic(err)
	}
	var temp Comp
	v := url.Values{}
	ab, err := db.Get(vars["id"], v)
	if err != nil {
		panic(err)
	}
	//reads the db of a specific document
	for _, value := range ab {
		switch value.(type) {
		case string:
		case map[string]interface{}:
			bolB, _ := json.Marshal(value)
			if err := json.Unmarshal(bolB, &temp); err != nil {
				panic(err)
			}
		}
	}
	for kr, _ := range temp.Kra {
		if kr == vars["kra"] {
			json.NewEncoder(w).Encode(temp.Kra[kr])
			return
		}
	}

}

//Adds new kra to a company using company id
func NewKra(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Comp
	vars := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&book)
	db, err := couchdb.NewDatabase("http://prasanna:prasanna@127.0.0.1:5984/company")
	if err != nil {
		panic(err)
	}
	var temp Comp
	v := url.Values{}
	ab, err := db.Get(vars["id"], v)
	if err != nil {
		panic(err)
	}
	for _, value := range ab {
		switch value.(type) {
		case string:
		case map[string]interface{}:
			bolB, _ := json.Marshal(value)
			if err := json.Unmarshal(bolB, &temp); err != nil {
				panic(err)
			}
		}
	}
	for _, y := range comp {
		if y.Cname == temp.Cname {
			for a, b := range book.Kra {
				y.Kra[a] = b
				temp.Kra[a] = b
			}
		}
	}
	var n int
	for k, _ := range temp.Kra {
		for a, _ := range book.Kra {
			if k == a {
				n = 1
			}
		}
	}
	for a, b := range book.Kra {
		if n == 1 {
			fmt.Fprintf(w, "kra already exists")
		} else {
			temp.Kra[a] = b
			fmt.Println(temp)
			ere := db.Delete(vars["id"])
			if ere != nil {
				panic(err)
			}
			km := make(map[string]interface{})
			fmt.Println(temp)
			km["company"] = temp
			err = db.Set(vars["id"], km)
			if err != nil {
				panic(err)
			}
			json.NewEncoder(w).Encode(book.Kra)
		}
	}
}

//Updates a Kra using company id and kra name
func UpdateKra(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var book Comp
	json.NewDecoder(r.Body).Decode(&book)
	db, err := couchdb.NewDatabase("http://prasanna:prasanna@127.0.0.1:5984/company")
	if err != nil {
		panic(err)
	}
	var temp Comp
	v := url.Values{}
	ab, err := db.Get(vars["id"], v)
	if err != nil {
		panic(err)
	}
	for _, value := range ab {
		switch value.(type) {
		case string:
			//fmt.Println(value)
		case map[string]interface{}:
			bolB, _ := json.Marshal(value)
			if err := json.Unmarshal(bolB, &temp); err != nil {
				panic(err)
			}
		}
	}
	for k, _ := range temp.Kra {
		for a, b := range book.Kra {
			if k == vars["kra"] {
				temp.Kra[a] = b
				delete(temp.Kra, vars["kra"])
				fmt.Println(temp)
				ere := db.Delete(vars["id"])
				if ere != nil {
					panic(err)
				}
				km := make(map[string]interface{})
				fmt.Println(temp)
				km["company"] = temp
				err = db.Set(vars["id"], km)
				if err != nil {
					panic(err)
				}
				json.NewEncoder(w).Encode(book.Kra)
			}
		}
	}
}

//Deletes Kra
func DeleteKra(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	db, err := couchdb.NewDatabase("http://prasanna:prasanna@127.0.0.1:5984/company")
	if err != nil {
		panic(err)
	}
	var temp Comp
	v := url.Values{}
	ab, err := db.Get(vars["id"], v)
	if err != nil {
		panic(err)
	}
	for _, value := range ab {
		switch value.(type) {
		case string:
			//fmt.Println(value)
		case map[string]interface{}:
			bolB, _ := json.Marshal(value)
			if err := json.Unmarshal(bolB, &temp); err != nil {
				panic(err)
			}
		}
	}
	for k, _ := range temp.Kra {
		if k == vars["kra"] {
			json.NewEncoder(w).Encode(temp.Kra[vars["kra"]])
			delete(temp.Kra, vars["kra"])
			ere := db.Delete(vars["id"])
			if ere != nil {
				panic(err)
			}
			km := make(map[string]interface{})
			fmt.Println(temp)
			km["company"] = temp
			err = db.Set(vars["id"], km)
			if err != nil {
				panic(err)
			}
		}

	}
}
func main() {
	r := mux.NewRouter()
	cm.Cname = "Param"
	cm.Kra = make(map[string]Kras)
	e, f := cm.Kra["K1"]
	if !f {
		e = Kras{}
	}
	e.Title = "Title "
	e.Kpi = make(map[string]string)
	cm.Kra["K1"] = e
	cm.Kra["K1"].Kpi["KPI1"] = "Metrics One"
	cm.Kra["K1"].Kpi["KPI2"] = "Metrics Two"
	db, err := couchdb.NewDatabase("http://prasanna:prasanna@127.0.0.1:5984/company")
	if err != nil {
		panic(err)
	}
	/*km := make(map[string]interface{})
	km["company"] = cm
	km["_rev"] = "13-3786a8428adf7d24ec337eeff07803e3"
	err = db.Set("1", km)
	if err != nil {
		panic(err)
	}*/
	//dummy db
	er := db.Contains("100")
	if er != nil {
		fmt.Println("doc not in db")
		km := make(map[string]interface{})
		km["company"] = cm
		err = db.Set("1", km)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("doc already in db")
	}
	comp = append(comp, Comp{Cname: cm.Cname, Kra: cm.Kra})
	r.HandleFunc("/addcmp/{id}", NewCompany).Methods("POST")
	r.HandleFunc("/cmp", Company).Methods("GET")
	r.HandleFunc("/kra/{id}/{kra}", Getkra).Methods("GET")
	r.HandleFunc("/newkra/{id}", NewKra).Methods("POST")
	r.HandleFunc("/updkra/{id}/{kra}", UpdateKra).Methods("PUT")
	r.HandleFunc("/delete/{id}/{kra}", DeleteKra).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(":1000", nil)
}
