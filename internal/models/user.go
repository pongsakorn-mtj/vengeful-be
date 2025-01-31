package models

import "time"

type RegisterRequest struct {
	FirstName             string `json:"firstName" binding:"required" validate:"required,min=2,max=50"`
	LastName              string `json:"lastName" binding:"required" validate:"required,min=2,max=50"`
	PhoneNo               string `json:"phoneNo" binding:"required" validate:"required,min=10,max=15"`
	Email                 string `json:"email" binding:"required" validate:"required,email"`
	IsAcceptTnc           bool   `json:"isAcceptTnc" binding:"required"`
	IsAcceptPrivacyPolicy bool   `json:"isAcceptPrivacyPolicy" binding:"required"`
}

type RegisterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserResponse struct {
	FirstName             string    `json:"firstName"`
	LastName              string    `json:"lastName"`
	PhoneNo               string    `json:"phoneNo"`
	Email                 string    `json:"email"`
	IsAcceptTnc           bool      `json:"isAcceptTnc"`
	IsAcceptPrivacyPolicy bool      `json:"isAcceptPrivacyPolicy"`
	CreatedAt             time.Time `json:"createdAt"`
}

type GetAllUsersResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    *GetUsersData `json:"data,omitempty"`
}

type GetUsersData struct {
	Users      []UserResponse `json:"users"`
	Pagination *Pagination    `json:"pagination"`
}

type Pagination struct {
	CurrentPage  int64 `json:"currentPage"`
	TotalPages   int64 `json:"totalPages"`
	TotalRecords int64 `json:"totalRecords"`
	Limit        int64 `json:"limit"`
}
