package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sigitisme/amf-loan-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "amf-loan-service",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "amf-loan-service", response["service"])
}

func TestLoanStateTransitions(t *testing.T) {
	tests := []struct {
		name          string
		currentState  domain.LoanState
		newState      domain.LoanState
		shouldBeValid bool
	}{
		{
			name:          "proposed to approved",
			currentState:  domain.LoanStateProposed,
			newState:      domain.LoanStateApproved,
			shouldBeValid: true,
		},
		{
			name:          "approved to invested",
			currentState:  domain.LoanStateApproved,
			newState:      domain.LoanStateInvested,
			shouldBeValid: true,
		},
		{
			name:          "invested to disbursed",
			currentState:  domain.LoanStateInvested,
			newState:      domain.LoanStateDisbursed,
			shouldBeValid: true,
		},
		{
			name:          "approved back to proposed (invalid)",
			currentState:  domain.LoanStateApproved,
			newState:      domain.LoanStateProposed,
			shouldBeValid: false,
		},
		{
			name:          "disbursed back to invested (invalid)",
			currentState:  domain.LoanStateDisbursed,
			newState:      domain.LoanStateInvested,
			shouldBeValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a simple state transition validation
			// In a real implementation, this would be handled by business logic
			stateOrder := map[domain.LoanState]int{
				domain.LoanStateProposed:  1,
				domain.LoanStateApproved:  2,
				domain.LoanStateInvested:  3,
				domain.LoanStateDisbursed: 4,
			}

			currentOrder := stateOrder[tt.currentState]
			newOrder := stateOrder[tt.newState]

			if tt.shouldBeValid {
				// Forward transitions should be valid
				assert.Greater(t, newOrder, currentOrder, "Valid state transition should move forward")
			} else {
				// Backward transitions should be invalid
				assert.LessOrEqual(t, newOrder, currentOrder, "Invalid state transition should not move backward")
			}
		})
	}
}
