package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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

//Gets the company details
func Company(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(comp)
}

//Adds a new company to db with a unique id which is not in db
func NewCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("mysql", "root:rootpassword@tcp(localhost:3306)/golang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var book Comp
	var cid int
	json.NewDecoder(r.Body).Decode(&book)
	comp = append(comp, book)
	insert, err := db.Query("INSERT INTO  golang.company(company_name) SELECT * FROM (SELECT (?)) AS tmp WHERE NOT EXISTS (SELECT company_name FROM golang.company WHERE company_name = (?)) LIMIT 1", book.Cname, book.Cname)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	rows, err := db.Query("SELECT comp_id FROM company WHERE company_name=(?)", book.Cname)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {

		if err := rows.Scan(&cid); err != nil {
			panic(err.Error())
		}
		defer rows.Close()
	}
	for v, _ := range book.Kra {
		e, f := book.Kra[v]
		if !f {
			e = Kras{}
		}
		insert, err := db.Query("INSERT INTO golang.kra(comp_id,kra_name,kra_title) SELECT (?),(?),(?) FROM DUAL WHERE NOT EXISTS (SELECT kra_name FROM golang.kra WHERE kra_name=(?))", cid, v, e.Title, v)
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
		var rid int
		rows, err := db.Query("SELECT kra_id FROM kra WHERE kra_name=(?)", v)
		if err != nil {
			panic(err.Error())
		}
		for rows.Next() {
			if err := rows.Scan(&rid); err != nil {
				panic(err.Error())
			}
			defer rows.Close()
		}
		for _, v := range book.Kra[v].Kpi {
			insert, err := db.Query("INSERT INTO  golang.kpi(kra_id,kpi_description) VALUES (?,?) ", rid, v)
			if err != nil {
				panic(err.Error())
			}
			defer insert.Close()
		}
		mrid := make([]int, 1)
		row, err := db.Query("SELECT kpi_id FROM golang.kpi ")
		if err != nil {
			fmt.Println("Failed to run query", err)
			return
		}
		for row.Next() {
			var k int
			if err := row.Scan(&k); err != nil {
				panic(err.Error())
			}
			mrid = append(mrid, k)
			defer row.Close()
		}
		for i := 1; i < len(mrid); i++ {
			res, err := db.Query("INSERT INTO golang.tracker(kpi_id,ass_period,goal_achieved) SELECT (?),(?),(?) FROM DUAL WHERE NOT EXISTS (SELECT kpi_id FROM golang.tracker WHERE kpi_id=(?))", mrid[i], time.Now(), "false", mrid[i])
			if err != nil {
				panic(err.Error())
			}
			defer res.Close()
		}

	}
	json.NewEncoder(w).Encode(book)
}

//Gets Kra of a specific company if it exists
func Getkra(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	for _, y := range comp {
		if y.Cname == vars["cmp"] {
			for kr, _ := range y.Kra {
				if kr == vars["kra"] {
					json.NewEncoder(w).Encode(y.Kra[kr])
					return
				}
			}
		} else {
			json.NewEncoder(w).Encode(&Comp{})
		}
	}
	json.NewEncoder(w).Encode(&Comp{})
}

//Adds new Kra to a Company
func NewKra(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(localhost:3306)/golang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var book Comp
	vars := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&book)

	for _, y := range comp {
		if y.Cname == vars["cmp"] {
			for a, b := range book.Kra {
				y.Kra[a] = b
			}
			json.NewEncoder(w).Encode(book.Kra)
		} else {
			json.NewEncoder(w).Encode(&Comp{})
		}
	}
	var x int
	rows, err := db.Query("SELECT comp_id FROM company WHERE company_name=(?)", vars["cmp"])
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		if err := rows.Scan(&x); err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		var krid int
		for _, y := range comp {
			for a, _ := range book.Kra {
				e, f := book.Kra[a]
				if !f {
					e = Kras{}
				}
				insert, err := db.Query("INSERT INTO golang.kra(comp_id,kra_name,kra_title) SELECT (?),(?),(?) FROM DUAL WHERE NOT EXISTS (SELECT kra_name FROM golang.kra WHERE kra_name=(?))", x, a, e.Title, a)
				//insert, err := db.Query("INSERT INTO kra(comp_id,kra_name,kra_title) VALUES (?,?,?)", cid, kr, e.Title)
				if err != nil {
					panic(err.Error())
				}
				defer insert.Close()
				rows, err := db.Query("SELECT kra_id FROM kra WHERE kra_name=(?)", a)
				if err != nil {
					panic(err.Error())
				}
				for rows.Next() {
					if err := rows.Scan(&krid); err != nil {
						panic(err.Error())
					}
					defer rows.Close()

					for _, v := range y.Kra[a].Kpi {
						insert, err := db.Query("INSERT INTO  golang.kpi(kra_id,kpi_description) vALUES (?,?) ", krid, v)
						if err != nil {
							panic(err.Error())
						}
						defer insert.Close()
						rid := make([]int, 1)
						rows, err := db.Query("SELECT kpi_id FROM golang.kpi ")
						if err != nil {
							fmt.Println("Failed to run query", err)
							return
						}
						for rows.Next() {
							var k int
							if err := rows.Scan(&k); err != nil {
								panic(err.Error())
							}
							rid = append(rid, k)
							defer rows.Close()
						}
						for i := 1; i < len(rid); i++ {
							res, err := db.Query("INSERT INTO golang.tracker(kpi_id,ass_period,goal_achieved) SELECT (?),(?),(?) FROM DUAL WHERE NOT EXISTS (SELECT kpi_id FROM golang.tracker WHERE kpi_id=(?))", rid[i], time.Now(), "true", rid[i])
							if err != nil {
								panic(err.Error())
							}
							defer res.Close()
						}
					}

				}
			}

		}
	}
}

//Updates a Kra
func UpdateKra(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(localhost:3306)/golang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var book Comp
	json.NewDecoder(r.Body).Decode(&book)
	var krid int

	for _, y := range comp {
		if y.Cname == vars["cmp"] {
			for kr, _ := range y.Kra {
				for a, b := range book.Kra {
					if kr == vars["kra"] {
						y.Kra[a] = b
						delete(y.Kra, kr)
						e, f := y.Kra[a]
						if !f {
							e = Kras{}
						}
						insert, err := db.Query("UPDATE kra SET kra_name=(?),kra_title=(?) WHERE kra_name=(?)", a, e.Title, kr)
						if err != nil {
							panic(err.Error())
						}
						defer insert.Close()
						rows, err := db.Query("SELECT kra_id FROM kra WHERE kra_name=(?)", a)
						if err != nil {
							panic(err.Error())
						}
						for rows.Next() {
							if err := rows.Scan(&krid); err != nil {
								panic(err.Error())
							}
							defer rows.Close()
						}
						rid := make([]int, 1)
						row, err := db.Query("SELECT kpi_id FROM golang.kpi WHERE kra_id=(?)", krid)
						if err != nil {
							fmt.Println("Failed to run query", err)
							return
						}
						for row.Next() {
							var k int
							if err := row.Scan(&k); err != nil {
								panic(err.Error())
							}
							rid = append(rid, k)
							defer row.Close()
						}
						for i := 1; i < len(rid); i++ {
							for _, v := range y.Kra[a].Kpi {
								//UPDATE kra SET kra_name=(?),kra_title=(?) WHERE kra_name=(?)
								insert, err := db.Query("UPDATE kpi SET kpi_description=(?) WHERE kpi_id=(?)", v, rid[i])
								if err != nil {
									panic(err.Error())
								}
								defer insert.Close()
							}
						}

					}

				}
			}
		}
	}
	json.NewEncoder(w).Encode(book.Kra)
}

//Delete a Kra
func DeleteKra(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(localhost:3306)/golang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var krid int
	for _, y := range comp {
		if y.Cname == vars["cmp"] {
			for kr, _ := range y.Kra {
				if kr == vars["kra"] {
					json.NewEncoder(w).Encode(y.Kra[kr])
					delete(y.Kra, kr)
					rows, err := db.Query("SELECT kra_id FROM kra WHERE kra_name=(?)", kr)
					if err != nil {
						panic(err.Error())
					}
					for rows.Next() {
						if err := rows.Scan(&krid); err != nil {
							panic(err.Error())
						}
						defer rows.Close()
					}
					rid := make([]int, 1)
					row, err := db.Query("SELECT kpi_id FROM golang.kpi WHERE kra_id=(?)", krid)
					if err != nil {
						fmt.Println("Failed to run query", err)
						return
					}
					for row.Next() {
						var k int
						if err := row.Scan(&k); err != nil {
							panic(err.Error())
						}
						rid = append(rid, k)
						defer row.Close()
					}
					for i := 1; i < len(rid); i++ {
						insert, err := db.Query("DELETE FROM tracker WHERE kpi_id=(?)", rid[i])
						if err != nil {
							panic(err.Error())
						}
						defer insert.Close()
					}
					inser, err := db.Query("DELETE FROM kpi WHERE kra_id=(?)", krid)
					if err != nil {
						panic(err.Error())
					}
					defer inser.Close()
					insert, err := db.Query("DELETE FROM kra WHERE kra_name=(?)", kr)
					if err != nil {
						panic(err.Error())
					}
					defer insert.Close()
				}
			}
		}
	}

}
func main() {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(localhost:3306)/golang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	r := mux.NewRouter()
	cm := new(Comp)
	cm.Cname = "Param"
	cm.Kra = make(map[string]Kras)
	e, f := cm.Kra["K1"]
	if !f {
		e = Kras{}
	}
	e.Title = "Title 1"
	e.Kpi = make(map[string]string)
	cm.Kra["K1"] = e
	cm.Kra["K1"].Kpi["KPI1"] = "Metrics One"
	cm.Kra["K1"].Kpi["KPI2"] = "Metrics Two"
	comp = append(comp, Comp{Cname: cm.Cname, Kra: cm.Kra})
	//Created a dummy data manually
	insert, err := db.Query("INSERT INTO  golang.company(company_name) SELECT * FROM (SELECT (?)) AS tmp WHERE NOT EXISTS (SELECT company_name FROM golang.company WHERE company_name = (?)) LIMIT 1", cm.Cname, cm.Cname)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	rows, err := db.Query("SELECT comp_id FROM company WHERE company_name=(?)", cm.Cname)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var cid int

		if err := rows.Scan(&cid); err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		for _, y := range comp {
			for kr, _ := range y.Kra {
				insert, err := db.Query("INSERT INTO golang.kra(comp_id,kra_name,kra_title) SELECT (?),(?),(?) FROM DUAL WHERE NOT EXISTS (SELECT kra_name FROM golang.kra WHERE kra_name=(?))", cid, kr, e.Title, kr)
				if err != nil {
					panic(err.Error())
				}
				defer insert.Close()
			}
		}
		var rid int
		for _, y := range comp {
			for kr, _ := range y.Kra {
				rows, err := db.Query("SELECT kra_id FROM kra WHERE kra_name=(?)", kr)
				if err != nil {
					panic(err.Error())
				}
				for rows.Next() {
					if err := rows.Scan(&rid); err != nil {
						panic(err.Error())
					}
					defer rows.Close()
				}
			}
			for _, y := range comp {
				for kr, _ := range y.Kra {
					for _, v := range cm.Kra[kr].Kpi {
						insert, err := db.Query("INSERT INTO  golang.kpi(kra_id,kpi_description) VALUES (?,?) ", rid, v)
						if err != nil {
							panic(err.Error())
						}
						defer insert.Close()
					}
				}
			}
			rid := make([]int, 1)
			rows, err := db.Query("SELECT kpi_id FROM golang.kpi ")
			if err != nil {
				fmt.Println("Failed to run query", err)
				return
			}
			for rows.Next() {
				var k int
				if err := rows.Scan(&k); err != nil {
					panic(err.Error())
				}
				rid = append(rid, k)
				defer rows.Close()
			}
			for i := 1; i < len(rid); i++ {
				res, err := db.Query("INSERT INTO golang.tracker(kpi_id,ass_period,goal_achieved) SELECT (?),(?),(?) FROM DUAL WHERE NOT EXISTS (SELECT kpi_id FROM golang.tracker WHERE kpi_id=(?))", rid[i], time.Now(), "true", rid[i])
				if err != nil {
					panic(err.Error())
				}
				defer res.Close()
			}
			r.HandleFunc("/cmp", Company).Methods("GET")
			r.HandleFunc("/addcmp", NewCompany).Methods("POST")
			r.HandleFunc("/kra/{cmp}/{kra}", Getkra).Methods("GET")
			r.HandleFunc("/newkra/{cmp}", NewKra).Methods("POST")
			r.HandleFunc("/updkra/{cmp}/{kra}", UpdateKra).Methods("PUT")
			r.HandleFunc("/delete/{cmp}/{kra}", DeleteKra).Methods("DELETE")
			http.Handle("/", r)
			http.ListenAndServe(":1000", nil)
		}
	}
}
