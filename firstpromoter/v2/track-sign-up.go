package firstpromoter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TrackSignUpInput struct {
	// ID to match the sale with the lead if the email can be changed before the first sale. Required if email is not provided.
	Uid *string
	// Email of the lead/sign-up. Required if uid is not provided.
	Email *string
	// Visitor tracking ID. It's set when the visitor tracking script tracks the referral visit on our system. The value is found inside _fprom_tid cookie. Required if ref_id is not provided.
	Tid *string
	// Default referral id of the promoter. Use this only when you want to assign the lead to a specific promoter. Required if tid is not provided.
	RefId *string
	// IP of the visitor who generated the sign up. It's used for fraud analysis.
	Ip *string
	// Date of the signup event
	CreatedAt *time.Time
	// Set this to true to skip email notifications. Default is false.
	SkipEmailNotification *bool
}

type Referral struct {
	Id    string  `json:"id"`
	Email *string `json:"email"`
	Uid   *string `json:"uid"`
}

type TrackSignUpOutput struct {
	Id                   int
	Etype                string
	SaleAmount           *int
	OriginalSaleAmount   *int
	OriginalSaleCurrency *string
	EventId              *string
	PlanId               *string
	BillingPeriod        *string
	CreatedAt            time.Time
	Referral             Referral
}

func (client Client) TrackSignUp(ctx context.Context, input TrackSignUpInput) (*TrackSignUpOutput, error) {
	var reqBody struct {
		Email                 string `json:"email,omitempty"`
		Uid                   string `json:"uid,omitempty"`
		Tid                   string `json:"tid,omitempty"`
		RefId                 string `json:"ref_id,omitempty"`
		Ip                    string `json:"ip,omitempty"`
		CreatedAt             string `json:"created_at,omitempty"`
		SkipEmailNotification bool   `json:"skip_email_notification"`
	}

	if input.Uid != nil {
		reqBody.Uid = *input.Uid
	}

	if input.Email != nil {
		reqBody.Email = *input.Email
	}

	if input.Tid != nil {
		reqBody.Tid = *input.Tid
	}

	if input.RefId != nil {
		reqBody.RefId = *input.RefId
	}

	if input.Ip != nil {
		reqBody.Ip = *input.Ip
	}

	if input.CreatedAt != nil {
		reqBody.CreatedAt = input.CreatedAt.UTC().Format(time.RFC3339)
	}

	if input.SkipEmailNotification != nil {
		reqBody.SkipEmailNotification = *input.SkipEmailNotification
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal signup request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://v2.firstpromoter.com/api/v2/track/signup", bytes.NewReader(reqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ACCOUNT-ID", client.accountId)
	req.Header.Set("Authorization", "Bearer "+client.apiKey)

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slurp, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("signup: unexpected status %s: %s", resp.Status, string(slurp))
	}

	var respBody struct {
		Id                   int      `json:"id"`
		Etype                string   `json:"etype"`
		SaleAmount           *int     `json:"sale_amount"`
		OriginalSaleAmount   *int     `json:"original_sale_amount"`
		OriginalSaleCurrency *string  `json:"original_sale_currency"`
		EventId              *string  `json:"event_id"`
		PlanId               *string  `json:"plan_id"`
		BillingPeriod        *string  `json:"billing_period"`
		CreatedAt            string   `json:"created_at"`
		Referral             Referral `json:"referral"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil && err != io.EOF {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, respBody.CreatedAt)
	if err == nil {
		return nil, fmt.Errorf("Error parsing CreatedAt")
	}

	return &TrackSignUpOutput{
		Id:                   respBody.Id,
		Etype:                respBody.Etype,
		SaleAmount:           respBody.SaleAmount,
		OriginalSaleAmount:   respBody.OriginalSaleAmount,
		OriginalSaleCurrency: respBody.OriginalSaleCurrency,
		EventId:              respBody.EventId,
		PlanId:               respBody.PlanId,
		BillingPeriod:        respBody.BillingPeriod,
		CreatedAt:            createdAt,
		Referral:             respBody.Referral,
	}, nil
}
