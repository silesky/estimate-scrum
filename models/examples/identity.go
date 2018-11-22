package models

import (
	"fmt"
)

type identity interface {
	GetID() string
	GetName() string
	GetAdminStatus() bool
}

type User struct {
	ID      string
	Name    string
	IsAdmin bool
}

// this is actually a method being declared on a struct... kind of like if i did
/*
class User implements Identity {
	private string Id
	private string Name
	private string IsAdmin
	GetId() {
		return this.Id
	}
	GetName() {
		return this.Name
	}
 GetName() {
		return this.Name
	}
	GetAdminStatus() {
		return this.IsAdmin
	}
*/

// GetID implements GetID interface on User Struct
func (u User) GetID() string {
	return u.ID
}

// GetName implement GetName interface on User Struct
func (u User) GetName() string {
	return u.Name
}

// implement isAdmin interface on User Struct (the argument name is the item that's it's getting implemented on)
func (u User) getAdminStatus() bool {
	return u.IsAdmin
} // implement

// PrintUserDetails takes a user
func PrintUserDetails(i identity) {
	fmt.Println(i.GetID())
	fmt.Println(i.GetName())
	fmt.Println(i.GetAdminStatus())
}
