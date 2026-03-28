package checkout

import (
	"github.com/adyen/adyen-go-api-library/v21/src/adyen"
	adyen_checkout "github.com/adyen/adyen-go-api-library/v21/src/checkout"
	"github.com/adyen/adyen-go-api-library/v21/src/common"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// checkout page is hosted on adyen's website, the frontend does not need to deal with any of this
const checkoutType = "hosted"

type CheckoutController struct {
	apiKey          string
	merchantAccount string
	themeId         string
	returnUrl       string

	adyenClient *adyen.APIClient
}

func NewCheckoutController(apiKey, merchantAccount, themeId, returnUrl string) CheckoutController {
	cc := CheckoutController{
		apiKey:          apiKey,
		merchantAccount: merchantAccount,
		themeId:         themeId,
		returnUrl:       returnUrl,
	}

	client := adyen.NewClient(&common.Config{
		ApiKey:      cc.apiKey,
		Environment: common.TestEnv,
	})

	cc.adyenClient = client

	return cc
}

func (cc *CheckoutController) TestCheckout(ctx *echo.Context) error {
	service := cc.adyenClient.Checkout()

	createCheckoutSessionRequest := adyen_checkout.CreateCheckoutSessionRequest{
		Reference:       uuid.New().String(),
		Mode:            common.PtrString(checkoutType),
		Amount:          eurAmount(10, 50),
		MerchantAccount: cc.merchantAccount,
		CountryCode:     common.PtrString("NL"),
		ThemeId:         common.PtrString(cc.themeId),
		ReturnUrl:       cc.returnUrl,
	}

	req := service.PaymentsApi.
		SessionsInput().
		IdempotencyKey(uuid.New().String()).
		CreateCheckoutSessionRequest(createCheckoutSessionRequest)

	res, httpRes, err := service.PaymentsApi.Sessions(ctx.Request().Context(), req)
	// we dont do anything special with the http response, but it could be useful for logging/debugging
	_ = httpRes

	if err != nil {
		log.Error("Could not complete adyen request", "error", err)
		return ctx.String(500, "could not complete request")
	}

	log.Info(res.Url)
	if res.Url == nil {
		log.Error("adyen's response url was empty")
		return ctx.String(500, "could not complete request")
	}

	return ctx.String(200, *res.Url)
}

func eurAmount(euros int64, cents int64) adyen_checkout.Amount {
	return adyen_checkout.Amount{
		Currency: "EUR",
		Value:    cents + (euros * 100),
	}
}
