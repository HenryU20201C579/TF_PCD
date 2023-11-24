package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
)

// GameData representa los datos del juego
type GameData struct {
	NumPlayers int
	GameMap    [40]int
	NumTurno   int
}

var puertoRemoto string

func main() {
	// Configurar el manejo de archivos estáticos
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Configurar las rutas de manejo
	http.HandleFunc("/", home)
	http.HandleFunc("/start_game", startGame)
	http.HandleFunc("/show_game", showGame)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func startGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	maxObstaculosStr := r.FormValue("maxObstaculos")
	maxObstaculos, err := strconv.Atoi(maxObstaculosStr)
	if err != nil {
		http.Error(w, "Invalid number of obstacles", http.StatusBadRequest)
		return
	}

	color := r.FormValue("opcion")
	setPuertoRemoto(color)

	direccionRemota := fmt.Sprintf("localhost:%s", puertoRemoto)

	var gameMap [40]int
	invalidPositions := []int{0, 39}
	initializeGameMap(&gameMap, invalidPositions, maxObstaculos)

	gameData := GameData{
		NumPlayers: 4,
		GameMap:    gameMap,
		NumTurno:   0,
	}

	jsonBytes, err := json.Marshal(gameData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonStr := string(jsonBytes)
	fmt.Println(jsonStr)
	fmt.Println(direccionRemota)

	con, _ := net.Dial("tcp", direccionRemota)
	defer con.Close()
	fmt.Fprintln(con, jsonStr)

	http.Redirect(w, r, "/show_game", http.StatusSeeOther)
}

func showGame(w http.ResponseWriter, r *http.Request) {
	contentRojo, err := ioutil.ReadFile("archivo_ROJO.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo ROJO:", err)
		http.Error(w, "Error al leer el archivo ROJO", http.StatusInternalServerError)
		return
	}

	fileContentRojo := string(contentRojo)

	// Leer el contenido del archivo de texto AZUL
	contentAzul, err := ioutil.ReadFile("archivo_AZUL.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo AZUL:", err)
		http.Error(w, "Error al leer el archivo AZUL", http.StatusInternalServerError)
		return
	}

	fileContentAzul := string(contentAzul)

	// Leer el contenido del archivo de texto VERDE
	contentVerde, err := ioutil.ReadFile("archivo_VERDE.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo VERDE:", err)
		http.Error(w, "Error al leer el archivo VERDE", http.StatusInternalServerError)
		return
	}

	fileContentVerde := string(contentVerde)

	contentAmarillo, err := ioutil.ReadFile("archivo_AMARILLO.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo AMARILLO:", err)
		http.Error(w, "Error al leer el archivo", http.StatusInternalServerError)
		return
	}
	fileContentAmarillo := string(contentAmarillo)

	data := struct {
		FileContentRojo     template.HTML
		FileContentAzul     template.HTML
		FileContentVerde    template.HTML
		FileContentAmarillo template.HTML
	}{
		FileContentRojo:     template.HTML(fileContentRojo),
		FileContentAzul:     template.HTML(fileContentAzul),
		FileContentVerde:    template.HTML(fileContentVerde),
		FileContentAmarillo: template.HTML(fileContentAmarillo),
	}

	tmpl, err := template.ParseFiles("templates/show_game.html")
	if err != nil {
		fmt.Println("Error al parsear el template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("Error al ejecutar el template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func initializeGameMap(tabla *[40]int, invalidPositions []int, maxObstaculos int) {
	var contador int

	for contador < maxObstaculos {
		number := rand.Intn(40)
		found := false
		for _, v := range invalidPositions {
			if number == v {
				found = true
				break
			}
		}
		if !found {
			contador++
			(*tabla)[number] = -1
		}
	}
}

func setPuertoRemoto(color string) {
	switch color {
	case "rojo":
		puertoRemoto = "8000"
	case "azul":
		puertoRemoto = "8001"
	case "verde":
		puertoRemoto = "8002"
	default:
		puertoRemoto = "8003"
	}
}
