package main

import (
	"fmt"

	"github.com/google/uuid"
)

type Student struct {
	ID uuid.UUID
	StudentData
}

type StudentData struct {
	Name    string
	Age     int
	Address string
}

func (s *Student) New() *Student {
	return &Student{
		ID: uuid.New(),
		StudentData: StudentData{
			Name:    "Vlad",
			Age:     20,
			Address: "Jakarta",
		},
	}
}

func (s *Student) GetID() uuid.UUID {
	return s.ID
}

func main() {
	s := new(Student).New()
	fmt.Printf("Student: %v\n", s)
}
