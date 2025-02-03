package main 

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/go-redis/redis/v8"
    "context"
)

// Global variables 
var ctx = context.Background()
var redisClient *redis.Client

// Init function to start the redis server 
func init() {
    redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB: 0,
    })
}

// Pokemon -> characteristics 
type Pokemon struct {
    Name  string `json: "name"`
    Type  string `json: "type"`
    XP    int    `json: "xp"`
    Power string `json: "power"`
    Level int    `json: "level"`
}

func getPokemonByType(pokemonType string) ([]Pokemon, error)  {
    keys, err := redisClient.Keys(ctx, fmt.Sprintf("pokemon:%s:*", pokemonType)).Result()
    if err != nil {
        return nil, err
    }

    var pokemons []Pokemon 
    for _, key := range keys {
        data, err := redisClient.Get(ctx, key).Result()
        if err != nil {
            return nil, err
        }
        var p Pokemon 
        if err := json.Unmarshal([]byte(data), &p); err != nil {
            return nil, err
        }
        pokemons = append(pokemons, p)
    }

    return pokemons, nil 
}

func handlePokemonType(w http.ResponseWriter, r *http.Request, pokemonType string) {
    pokemons, err := getPokemonByType(pokemonType) 
    if err != nil {
        http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
        return 
    }
    jsonResponse(w, pokemons)
}

// Json Response Encoder 
func jsonResponse(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, "Failed to encode data", http.StatusInternalServerError)
    }
}

func main() {
    http.HandleFunc("/water", func(w http.ResponseWriter, r *http.Request) {
        handlePokemonType(w, r, "water")
    })
    http.HandleFunc("/electric", func(w http.ResponseWriter, r *http.Request) {
        handlePokemonType(w, r, "electric")
    })
    http.HandleFunc("/grass", func(w http.ResponseWriter, r *http.Request) {
        handlePokemonType(w, r, "grass")
    })
    http.HandleFunc("/legendary", func(w http.ResponseWriter, r *http.Request) {
        handlePokemonType(w, r, "legendary")
    })
    http.HandleFunc("/fire", func(w http.ResponseWriter, r *http.Request) {
        handlePokemonType(w, r, "fire")
    })

    fmt.Println("Starting Pokemon API server on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
