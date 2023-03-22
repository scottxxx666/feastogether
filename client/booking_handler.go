package client

import (
	"encoding/json"
	"feastogether/config"
	"feastogether/fetch"
	"log"
)

// api
const (
	LOGIN_API      = "https://www.feastogether.com.tw/api/994f5388-d001-4ca4-a7b1-72750d4211cf/custSignIn"
	SAVE_SEATS_API = "https://www.feastogether.com.tw/api/booking/saveSeatss"
	SAVE_SAETA_API = "https://www.feastogether.com.tw/api/booking/saveSaetas"
	BOOKNIG_API    = "https://www.feastogether.com.tw/api/booking/booknigs"
	BOOKING_API    = "https://www.feastogether.com.tw/api/booking/bookings"

	// old api
	SAVE_SAETS_API = "https://www.feastogether.com.tw/api/booking/saveSaets"
	B00KING_API    = "https://www.feastogether.com.tw/api/booking/b00king"
)

var MealSeqMap = map[string]int{
	"11:30": 1,
	"12:30": 1,
	"14:30": 3,
	"17:30": 4,
	"18:00": 4,
	"18:30": 4,
}

// 取得 Token
func GetToken(user config.UserConfig) string {

	payload := Login{
		Act: user.Account,
		Pwd: user.Password,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal struct to JSON:%v\n", err)
		return ""
	}

	resp, err := fetch.Post(LOGIN_API, payloadBytes, user, "")
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()

	var data Response
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode response body: %v\n", err)
		return ""
	}
	if data.StatusCode != 1000 {
		log.Println(data.Result.Msg)
		return ""
	}
	return data.Result.CustomerLoginResp.Token
}

// 取得驗證序號
func GetSaveSaets(user config.UserConfig, token string) string {

	resp, err := fetch.Post(SAVE_SAETA_API, nil, user, token)
	if err != nil {
		log.Println(err)
		return ""
	}

	defer resp.Body.Close()

	var data SaveSaetsResponse
	if json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode response body: %v\n", err)
		return ""
	}

	if data.StatusCode != 1000 {
		log.Printf("取得驗證序號失敗 : %v\n", data)
		return ""
	}
	return data.Result
}

// 立即定位
func GetSaveSeats(user config.UserConfig, token string, payload config.RestaurantConfig) string {

	saveSeats := SaveSeats{
		StoreID:     payload.StoreID,
		PeopleCount: payload.PeopleCount,
		MealPeriod:  payload.MealPeriod,
		MealDate:    payload.MealDate,
		MealTime:    payload.MealTime,
		MealSeq:     MealSeqMap[payload.MealTime],
		// Zkde:        nil,

		// 3/20 新參數
		Zked: nil,
	}

	if saets := GetSaveSaets(user, token); saets != "" {
		// saveSeats.Zkde = strings.ReplaceAll(saets, "I", "l")

		// 3/20 更新序號邏輯
		saveSeats.Zked = saets
	}

	payloadBytes, err := json.Marshal(saveSeats)
	if err != nil {
		log.Printf("Failed to marshal struct to JSON:%v\n", err)
		return ""
	}

	resp, err := fetch.Post(SAVE_SEATS_API, payloadBytes, user, token)
	if err != nil {
		log.Println(err)
		return ""
	}

	defer resp.Body.Close()

	var data Response
	if json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode response body: %v\n", err)
		return ""
	}

	if data.StatusCode != 1000 {
		log.Println(data.Result.Msg)
		return ""
	}

	return data.Result.ExpirationTime
}

// 取得驗證序號
func GetB00king(user config.UserConfig, token string) string {

	resp, err := fetch.Post(BOOKNIG_API, nil, user, token)
	if err != nil {
		log.Println(err)
		return ""
	}

	defer resp.Body.Close()

	var data SaveSaetsResponse
	if json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode response body: %v\n", err)
		return ""
	}

	if data.StatusCode != 1000 {
		log.Printf("取得驗證序號失敗 : %v\n", data)
		return ""
	}
	return data.Result
}

// 送出定位
func SaveBooking(user config.UserConfig, token string, payload config.RestaurantConfig) string {

	booking := Booking{
		StoreID:     payload.StoreID,
		MealPeriod:  payload.MealPeriod,
		MealDate:    payload.MealDate,
		MealTime:    payload.MealTime,
		MealPurpose: "",
		MealSeq:     MealSeqMap[payload.MealTime],
		Special:     0,
		ChildSeat:   0,
		Adult:       payload.PeopleCount,
		Child:       0,
		ChargeList: []struct {
			Seq   int "json:\"seq\""
			Count int "json:\"count\""
		}{
			// 大人
			{
				Seq:   201,
				Count: payload.PeopleCount,
			},
			// 小孩
			{
				Seq:   202,
				Count: 0,
			},
		},
		StoreCode:    "NTBQ",
		RedirectType: "iEat_card",
		Domain:       "https://www.feastogether.com.tw",
		PathFir:      "booking",
		PathSec:      "result",
		// YuuO:         GetB00king(user, token),

		// 3/20 新參數
		Yuuu: GetB00king(user, token),
	}

	payloadBytes, err := json.Marshal(booking)
	if err != nil {
		log.Printf("Failed to marshal struct to JSON:%v\n", err)
		return ""
	}

	resp, err := fetch.Post(BOOKING_API, payloadBytes, user, token)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()

	var data Response
	if json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode response body: %v\n", err)
		return ""
	}

	if data.StatusCode != 1000 {
		log.Println(data.Result.Msg)
		return ""
	}
	return data.Message
}
