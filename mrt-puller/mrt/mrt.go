package mrt

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/mrt-puller/config"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/mrt-puller/model"
	"log/slog"
)

type Api interface {
	GetJourney(req model.GetJourneyRequest) (model.GetJourneyResponse, error)
}

type ApiImpl struct {
	client   *resty.Client
	email    string
	password string
}

func NewApiImpl(c config.MrtApiConfig) *ApiImpl {
	return &ApiImpl{
		client: resty.New().
			SetBaseURL(c.BaseURL).
			SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Safari/605.1.15").
			SetHeader("Cache-Control", "no-cache"),
		email:    c.Email,
		password: c.Password,
	}
}

func (a ApiImpl) GetJourney(ctx context.Context, req model.GetJourneyRequest) (model.GetJourneyResponse, error) {
	var resp model.GetJourneyResponse
	_, err := a.client.R().
		SetBody(req).
		SetResult(&resp).
		SetContext(ctx).
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Language", string(req.Lang)).
		Post("journey")

	if err != nil {
		return model.GetJourneyResponse{}, err
	}

	if resp.Meta.ResponseCode != 600 {
		slog.ErrorContext(ctx, "Unsuccessful response code", slog.Any("response", resp))
		return model.GetJourneyResponse{}, errors.New("unsuccessful response code")
	}

	return resp, nil
}

func (a ApiImpl) GetAccessToken(ctx context.Context) (string, error) {
	req := model.LoginRequest{
		Email:    a.email,
		Password: a.password,
	}
	var resp model.LoginResponse

	_, err := a.client.R().
		SetBody(req).
		SetResult(&resp).
		SetContext(ctx).
		Post("login")
	if err != nil {
		return "", err
	}

	if resp.Meta.ResponseCode != 600 {
		slog.ErrorContext(ctx, "Unsuccessful response code", slog.Any("response", resp))
		return "", errors.New("unsuccessful response code")
	}

	if resp.Data.AccessToken == "" {
		slog.Error("Empty access token", slog.Any("response", resp))
		return "", errors.New("empty access token")
	}

	return resp.Data.AccessToken, nil
}
