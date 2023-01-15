package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "github.com/pborman/uuid"
    // "github.com/gorilla/mux"

    "appstore/service"
    "appstore/model"
    jwt "github.com/form3tech-oss/jwt-go"

)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    
    // http.ResponseWriter is an interface
    // http.Request is a struct
    
    // Parse from body of request to get a json object.
    fmt.Println("Received one upload request")

    // Read from token
    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

    // Define app
    app := model.App{
        Id:      uuid.New(),
        User:    username.(string),
        Title:   r.FormValue("title"),
        Description: r.FormValue("description"),
    }
    
    // Read the price
    price, err := strconv.ParseFloat(r.FormValue("price"), 64)
    fmt.Printf("%v,%T", price, price)
    if err != nil {
        fmt.Println(err)
    }
    app.Price = int(price * 100.0)
    
    // Read Media file
    file, _, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }
    
    // Save app to ES
    err = service.SaveApp(&app, file)
    if err != nil {
        http.Error(w, "Failed to save app to backend", http.StatusInternalServerError)
        fmt.Printf("Failed to save app to backend %v\n", err)
        return
    }
    fmt.Println("App is saved successfully.")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
   fmt.Println("Received one search request")
   w.Header().Set("Content-Type", "application/json")
   title := r.URL.Query().Get("title")
   description := r.URL.Query().Get("description")
   username := r.URL.Query().Get("username")


   var apps []model.App
   var err error
   apps, err = service.SearchApps(title, description, username)
   if err != nil {
       http.Error(w, "Failed to read Apps from backend", http.StatusInternalServerError)
       return
   }


   js, err := json.Marshal(apps)
   if err != nil {
       http.Error(w, "Failed to parse Apps into JSON format", http.StatusInternalServerError)
       return
   }
   w.Write(js)
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one checkout request")
    w.Header().Set("Content-Type", "text/plain")
    if r.Method == "OPTIONS" {
        return
    }
 
    appID := r.FormValue("appID")
    s, err := service.CheckoutApp(r.Header.Get("Origin"), appID)
    if err != nil {
        fmt.Println("Checkout failed.")
        w.Write([]byte(err.Error()))
        return
    }
 
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(s.URL))
 
    fmt.Println("Checkout process started!")
 }
 
 
func deleteHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for delete")

    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"].(string)
    // id := mux.Vars(r)["id"]
    id := r.URL.Query().Get("item_id")

    if err := service.DeleteApp(id, username); err != nil {
        http.Error(w, "Failed to delete app from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to delete app from backend %v\n", err)
        return
    }
    fmt.Println("App is deleted successfully")
}

