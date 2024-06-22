package main

func main() {
	doctor := &Doctor{
		name: "doctor",
	}

	doctor.init()

	doctor.run()
}
