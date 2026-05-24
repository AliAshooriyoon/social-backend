package main

func main() {
	app := application{
		config: config{
			addr: ":8080",
		},
	}
	app.run()
}
