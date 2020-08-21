//go:generate go run ../../../../cmd/generator $GOFILE ../../repository/$GOPACKAGE/repository.go
//go:generate go generate ../../repository/$GOPACKAGE/repository.go
package user

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
