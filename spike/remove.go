package main

import "log"

func main() {
	var slice []int

	slice = append(slice, 1)
	slice = append(slice, 2)
	slice = append(slice, 3)
	slice = append(slice, 4)
	slice = append(slice, 5)

	for i, v := range slice {
		if i == 2 {
			l := len(slice)
			slice[i] = slice[l-1]
			slice = slice[:l-1]
			continue
		}
		log.Printf("%d: %d\n", i, v)
	}
	log.Println()

	for i, v := range slice {
		log.Printf("%d: %d\n", i, v)
	}
}
