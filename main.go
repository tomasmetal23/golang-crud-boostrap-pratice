package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Funcion para la conexion y autenticacion con la Base de datos
func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Password := "nEis4bGZe"
	Nombre := "crud"

	conexion, err := sql.Open(Driver, Usuario+":"+Password+"@tcp(10.89.2.3)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion

}

// plantilla
var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	//Rutas
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)
	//Para integrar un nuevo directorio local en el proyecto
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Running Service...")

	http.ListenAndServe(":8080", nil)
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()

	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string
		var fired_at sql.NullString
		err = registros.Scan(&id, &nombre, &correo, &fired_at)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)

	}
	//para mostrar el arreglo empleado en la consola ##fmt.Println(arregloEmpleado)

	//fmt.Fprintf(w, "Hello Saiyans!!!")
	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)

}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)

}

// funccion que inserta datos con el metodo "POST" desde un formulario a la BD
func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()

		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES(?,?) ")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)

	}
}

// funccion que actualiza los datos desde un form con el metodo "POST" desde un formulario a la BD
func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()

		modificarRegistros, err := conexionEstablecida.Prepare(" UPDATE empleados SET nombre=?,correo=? WHERE id=? ")

		if err != nil {
			panic(err.Error())
		}
		modificarRegistros.Exec(nombre, correo, id)

		http.Redirect(w, r, "/", 301)

	}
}

// Funccion editar
func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	fmt.Println(idEmpleado)
	conexionEstablecida := conexionBD()

	registros, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)
	empleado := Empleado{}
	for registros.Next() {
		var id int
		var nombre, correo string
		var fired_at sql.NullString
		err = registros.Scan(&id, &nombre, &correo, &fired_at)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}

	fmt.Println(empleado)
	plantillas.ExecuteTemplate(w, "editar", empleado)

}

// Funcion borrar
func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	fmt.Println(idEmpleado)

	conexionEstablecida := conexionBD()
	//Instruccion SQL para borrar
	borrarRegistros, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistros.Exec(idEmpleado)

	http.Redirect(w, r, "/", 301)

}

type Empleado struct {
	Id       int
	Nombre   string
	Correo   string
	fired_at string
}
