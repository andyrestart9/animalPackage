package main

import "fmt"

type Human struct{ Name string }
func (h Human) Speak() string { return "Hi, I'm " + h.Name }

type Robot struct{ ID int }
func (r Robot) Speak() string { return fmt.Sprintf("Beep! I am robot #%d", r.ID) }

// 沒有介面，只好各寫一支
func AnnounceHuman(h Human) { fmt.Println(h.Speak()) }
func AnnounceRobot(r Robot) { fmt.Println(r.Speak()) }

func main() {
	AnnounceHuman(Human{"Andy"})
	AnnounceRobot(Robot{42})
}