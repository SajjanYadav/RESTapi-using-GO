package main

func main(){
	app := App{}
	app.Initialise(DbUser,DbPassword,DBName)
	app.Run("localhost:10000")
}