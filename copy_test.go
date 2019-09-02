package goutils

import (
	"fmt"
	"testing"
)

type People struct {
	PeopleName string
	Age        int
	Gender     int
}

func (p *People) DoubleAge() int {
	return p.Age * 2
}

type User struct {
	UserName  string
	Age       int
	DoubleAge int
	Gender    int
}

func (u *User) PeopleName(peopleName string) {
	u.UserName = peopleName
}

func (u *User) AgeCopy(age int) {
	u.Age = age + 10
}

func TestCopy(t *testing.T) {
	var (
		people  = People{PeopleName: "Tom", Age: 18, Gender: 1}
		peoples = []People{{PeopleName: "Tom", Age: 18, Gender: 1}, {PeopleName: "Jerry", Age: 20, Gender: 1}}
		user    = User{}
		users   = make([]User, 0)
	)

	_ = Copy(&user, &people)

	fmt.Printf("%#v \n", user)
	// User{
	//    UserName: "Tom",       // Copy from method PeopleName
	//    Age: 28,               // Copy from method AgeCopy
	//    DoubleAge: 36,         // Copy from method DoubleAge
	//    Gender: 1,             // Copy from field
	// }

	// Copy struct to slice
	_ = Copy(&users, &people)

	fmt.Printf("%#v \n", users)
	// []User{
	//   {UserName: "Tom", Age: 28, DoubleAge: 36, Gender: 1}
	// }

	// Copy slice to slice
	users = make([]User, 0)
	_ = Copy(&users, &peoples)

	fmt.Printf("%#v \n", users)
	// []User{
	//   {Name: "Tom", Age: 28, DoubleAge: 36, Gender: 1},
	//   {Name: "Jerry", Age: 30, DoubleAge: 40, Gender: 1},
	// }

}
