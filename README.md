
# Movie ticket booking App

This is the server side code based on movie ticket booking application.

### How to get started 
 1. Clone the repository and move to the specified folder.
 2. Move to the directory by `cd cmd/main`
 3. Build the application using `go build`
 4. Run the application usinf `go run ./main.go `
 5. The application will run at `localhost:9010`
 6. Create a database and modify the path  in `app.go` file .
 `d, err  := gorm.Open("mysql","root:<password>@/<database_name>?charset=utf8&parseTime=True&loc=Local")`
 7. Modify the code with your email and password
 `senderEmail  := os.Getenv("SENDER_EMAIL")
senderPassword  := os.Getenv("SENDER_PASSWORD")`
 Yay your client side is set up and you are good to go 💪.

