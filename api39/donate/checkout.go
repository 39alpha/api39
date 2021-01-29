package donate

import (
	"encoding/json"
	"fmt"
	"github.com/39alpha/api39/api39"
	"github.com/kataras/iris/v12"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

type DonationRequest struct {
	Currency    string  `json:"currency"`
	Email       string  `json:"email"`
	ProductName string  `json:"product"`
	UnitAmount  float64 `json:"unit_amount"`
}

func ParseDonationRequest(ctx iris.Context) (*DonationRequest, error) {
	if ctx.Method() == "POST" {
		body, err := ctx.GetBody()
		if err != nil {
			return nil, err
		}

		var parsed DonationRequest
		if err := json.Unmarshal(body, &parsed); err != nil {
			return nil, err
		} else {
			return &parsed, nil
		}
	}
	return nil, fmt.Errorf("%s requests are not supported", ctx.Method())
}

type CheckoutSessionResponse struct {
	SessionId string `json:"id"`
}

func Checkout(ctx iris.Context) {
	req, err := ParseDonationRequest(ctx)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("invalid request").
			Type("error").
			Detail(err.Error()))
		return
	}

	if req.UnitAmount < 1.0 {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("invalid unit_amount").
			Type("amount-too-low"))
		return
	} else if req.UnitAmount >= 515.24 {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("invalid unit_amount").
			Type("amount-too-high"))
		return
	}

	cfg, ok := ctx.Values().Get("config").(*api39.Config)
	if !ok {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("failed to load configuration").
			Type("error"))
		return
	}

	stripe.Key = cfg.Stripe.Apikey
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		CustomerEmail:            &req.Email,
		SubmitType:               stripe.String("donate"),
		BillingAddressCollection: stripe.String("required"),
		Mode:                     stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(cfg.Stripe.Currency),

					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(cfg.Stripe.Product),
					},
					UnitAmount: stripe.Int64((int64)(100 * req.UnitAmount)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(cfg.Stripe.SuccessURL),
		CancelURL:  stripe.String(cfg.Stripe.CancelURL),
	}

	session, err := session.New(params)
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	data := CheckoutSessionResponse{
		SessionId: session.ID,
	}

	_, _ = ctx.JSON(data)
}
