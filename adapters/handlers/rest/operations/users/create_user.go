//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2024 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/weaviate/weaviate/entities/models"
)

// CreateUserHandlerFunc turns a function with the right signature into a create user handler
type CreateUserHandlerFunc func(CreateUserParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateUserHandlerFunc) Handle(params CreateUserParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// CreateUserHandler interface for that can handle valid create user params
type CreateUserHandler interface {
	Handle(CreateUserParams, *models.Principal) middleware.Responder
}

// NewCreateUser creates a new http.Handler for the create user operation
func NewCreateUser(ctx *middleware.Context, handler CreateUserHandler) *CreateUser {
	return &CreateUser{Context: ctx, Handler: handler}
}

/*
	CreateUser swagger:route POST /users/{user_id} users createUser

create new user
*/
type CreateUser struct {
	Context *middleware.Context
	Handler CreateUserHandler
}

func (o *CreateUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	Params := NewCreateUserParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		*r = *aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)
}
