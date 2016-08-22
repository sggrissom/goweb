package main

import (
	"fmt"
)

func RegisterPages() {
	pages["test"] = TestPage
	pages["test2"] = Test2Page
	pages["settings"] = SettingsPage
}

func SettingsPage(page string) {
	fmt.Printf("settings\n")
}

func TestPage(page string) {
	fmt.Printf("test\n")
}

func Test2Page(page string) {
	fmt.Printf("test2\n")
}
