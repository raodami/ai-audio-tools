package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/webhook"
	"ai-audio-tools/internal/store"
)

// SetupWebhookRoutes sets up Stripe webhook endpoints
func SetupWebhookRoutes(r *gin.Engine, s *store.Store) {
	r.POST("/webhook", func(c *gin.Context) {
		webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
		if webhookSecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stripe webhook not configured"})
			return
		}

		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get raw data"})
			return
		}

		event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), webhookSecret)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Webhook error: %v", err)})
			return
		}

		switch event.Type {
		case "checkout.session.completed":
			handleCheckoutSessionCompleted(s, event)
		case "invoice.paid":
			handleInvoicePaid(s, event)
		case "customer.subscription.deleted":
			handleSubscriptionDeleted(s, event)
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
}

func handleCheckoutSessionCompleted(s *store.Store, event stripe.Event) {
	var session struct {
		Customer string          `json:"customer"`
		Metadata json.RawMessage `json:"metadata"`
	}
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return
	}

	var metadata map[string]string
	if err := json.Unmarshal(session.Metadata, &metadata); err != nil {
		return
	}

	userID, ok := metadata["user_id"]
	if !ok {
		return
	}

	// Mark user as pro
	if err := s.SetUserPro(userID, true); err != nil {
		fmt.Printf("Failed to update user %s: %v", userID, err)
	}
}

func handleInvoicePaid(s *store.Store, event stripe.Event) {
	var invoice struct {
		Subscription string `json:"subscription"`
	}
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return
	}
	// Can look up subscription to get customer ID
	_ = invoice
}

func handleSubscriptionDeleted(s *store.Store, event stripe.Event) {
	var sub struct {
		Customer string `json:"customer"`
	}
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return
	}
	if err := s.SetUserPro(sub.Customer, false); err != nil {
		fmt.Printf("Failed to update user %s: %v", sub.Customer, err)
	}
}
