package model

type Lang string

const (
	LangEN = "en_EN"
	LangTH = "th_TH"
)

type MrtApiMeta struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Language        string `json:"language"`
	Version         string `json:"version"`
	Time            string `json:"time"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Ref1     string `json:"ref1"`
}

type LoginResponse struct {
	Data struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		IsRequireOTP bool   `json:"isRequireOTP"`
		User         struct {
			FirstName       string `json:"firstName"`
			LastName        string `json:"lastName"`
			Email           string `json:"email"`
			PhoneNumber     string `json:"phoneNumber"`
			IsAcceptConsent bool   `json:"isAcceptConsent"`
			IsRenewPolicy   bool   `json:"isRenewPolicy"`
			IsRenewTerm     bool   `json:"isRenewTerm"`
			IsConsentFail   bool   `json:"isConsentFail"`
		} `json:"user"`
	} `json:"data"`
	Meta MrtApiMeta `json:"meta"`
}

type GetJourneyRequest struct {
	CardID      string `json:"cardId"`
	PageNo      int    `json:"pageNo"`
	PageSize    int    `json:"pageSize"`
	AccessToken string `json:"-"`
	Lang        Lang   `json:"-"`
}

type GetJourneyResponse struct {
	Data struct {
		List []struct {
			TravelDate string        `json:"travelDate"`
			Journeys   []JourneyData `json:"journeys"`
		} `json:"list"`
		PageNo    int `json:"pageNo"`
		PageSize  int `json:"pageSize"`
		TotalPage int `json:"totalPage"`
	} `json:"data"`
	Meta MrtApiMeta `json:"meta"`
}
type JourneyData struct {
	JourneyID string `json:"journeyId"`
	From      struct {
		StationCode  string `json:"stationCode"`
		StationName  string `json:"stationName"`
		StationColor string `json:"stationColor"`
		Date         string `json:"date"`
	} `json:"from"`
	To struct {
		StationCode  string `json:"stationCode"`
		StationName  string `json:"stationName"`
		StationColor string `json:"stationColor"`
		Date         string `json:"date"`
	} `json:"to"`
	CardNumber  string `json:"cardNumber"`
	Date        string `json:"date"`
	Status      string `json:"status"`
	StatusText  string `json:"statusText"`
	TotalTime   string `json:"totalTime"`
	TotalAmount int    `json:"totalAmount"`
	PassName    string `json:"passName"`
	PassID      string `json:"passId"`
}
