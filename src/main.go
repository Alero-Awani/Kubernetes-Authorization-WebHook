package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	authorizationv1 "k8s.io/api/authorization/v1"
)

func main() {
	app := fiber.New()

	app.Post("/authorize", func(ctx *fiber.Ctx) error {
		var req authorizationv1.SubjectAccessReview
		ctx.BodyParser(&req)
		req.Status.Allowed = true

		if req.Spec.User == "system:serviceaccount:default:test-user" {
			if req.Spec.ResourceAttributes != nil {
				if req.Spec.ResourceAttributes.Verb == "get" || req.Spec.ResourceAttributes.Verb == "list" {
					req.Status.Allowed = true
					log.Printf("Allowed request for user %s on resource %s with verb %s", req.Spec.User, req.Spec.ResourceAttributes.Resource, req.Spec.ResourceAttributes.Verb)
				}

				if req.Spec.ResourceAttributes.Verb == "delete" {
					req.Status.Allowed = false
					req.Status.Reason = "You are not allowed to delete this resource"
					req.Status.Denied = true
					log.Printf("Denied request for user %s on resource %s with verb %s", req.Spec.User, req.Spec.ResourceAttributes.Resource, req.Spec.ResourceAttributes.Verb)
				}
			}

			if req.Spec.NonResourceAttributes != nil {
				if req.Spec.NonResourceAttributes.Verb == "get" || req.Spec.NonResourceAttributes.Verb == "list" {
					req.Status.Allowed = true
				}

				if req.Spec.NonResourceAttributes.Verb == "delete" {
					req.Status.Allowed = false
				}
			}
		}
		return ctx.JSON(req)
	})

	app.Get("/healthz", func(ctx *fiber.Ctx) error {
		fmt.Println("healthz")
		return ctx.SendStatus(200)
	})

	if err := app.ListenTLS(":443", "webhook.crt", "webhook.key"); err != nil {
		fmt.Println(err)
	}
}