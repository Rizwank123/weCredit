package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/weCredit/internal/domain"
	"github.com/weCredit/internal/pkg/config"
)

var daysOfWeek = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}

// AppUtil ... Service providing utility functionality throughout the application
type AppUtil interface {
	// GetCurrentTime ... Get the current time from the system
	GetCurrentTime() time.Time

	// GenerateOTP ... Generate a one time password
	GenerateOTP(length int) string
	// GenerateUniqueToken ... Generate a unique token
	GenerateUniqueToken() string
	// GetExpiryTimeForDuration ... Get an expiry time based on the duration (in hours) passed
	GetExpiryTimeForDuration(duration int) time.Time

	// CompareSlices ... Find the elements in one array of string but not in the other
	CompareSlices(a, b []string) []string

	// ParseStringForTime ... Parse string into time.RFC3339 format
	ParseStringForTime(date string) (time.Time, error)
	// ParseStringForTimeWithLocation ... Parse string into time.RFC3339 format
	ParseStringForTimeWithLocation(date string, loc *time.Location) (time.Time, error)
	// FormatDate ... Format date to get day of month with suffix
	FormatDate(t time.Time) string
	// ParseWeekday ... Parses a string and returns corresponding weekday for it
	//
	// For example, if you pass "Monday", it returns "1"
	ParseWeekday(v string) (time.Weekday, error)
	// IsTimeExpired ... Validate if the specified time has expired based on the current time
	IsTimeExpired(t time.Time) bool
	// SendOtp sends the otp using twilio
	SendOtp(cfg config.WeCreditConfig, sms domain.OtpMessage) error
}

// NewAppUtil ... Creates a new AppUtil
func NewAppUtil() AppUtil {
	return &simpleAppUtil{}
}

type simpleAppUtil struct{}

func (as *simpleAppUtil) GetCurrentTime() time.Time {
	return time.Now()
}

func (as *simpleAppUtil) GenerateOTP(length int) string {
	const digits = "1234567890"
	if length <= 0 {
		return ""
	}
	buffer := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, buffer, length)
	if n != length || err != nil {
		return ""
	}
	for i := 0; i < len(buffer); i++ {
		buffer[i] = digits[int(buffer[i])%len(digits)]
	}
	return string(buffer)
}

func (as *simpleAppUtil) GenerateUniqueToken() string {
	code, _ := uuid.NewV4()
	return code.String()
}

func (as *simpleAppUtil) GetExpiryTimeForDuration(duration int) time.Time {
	t := as.GetCurrentTime().Add(time.Hour*time.Duration(duration) + time.Minute*0 + time.Second*0)
	return t
}

func (as *simpleAppUtil) CompareSlices(a, b []string) (diff []string) {
	m := make(map[string]bool)
	for _, item := range b {
		m[item] = true
	}
	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func (as *simpleAppUtil) ParseStringForTime(date string) (time.Time, error) {
	tm, err := time.Parse(time.RFC3339, date)
	return tm, err
}

func (as *simpleAppUtil) ParseStringForTimeWithLocation(date string, loc *time.Location) (time.Time, error) {
	tm, err := time.ParseInLocation(time.RFC3339, date, loc)
	return tm, err
}

func (as *simpleAppUtil) IsTimeExpired(t time.Time) bool {
	return time.Since(t).Seconds() > 0
}

func (as *simpleAppUtil) FormatDate(t time.Time) string {
	suffix := "th"
	switch t.Day() {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return t.Format("2" + suffix + " Jan")
}

func (as *simpleAppUtil) ParseWeekday(v string) (time.Weekday, error) {
	if d, ok := daysOfWeek[v]; ok {
		return d, nil
	}

	return time.Sunday, fmt.Errorf("invalid weekday '%s'", v)
}

func (as *simpleAppUtil) SendOtp(cfg config.WeCreditConfig, sms domain.OtpMessage) error {
	// TODO: implement this
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSSID,
		Password: cfg.AccountAuthToken,
	})

	// Define a predefined SMS template
	messageBody := fmt.Sprintf("Your OTP is: %s. Please use this to complete your login. Do not share this code with anyone.", sms.Otp)

	// Prepare the message parameters
	params := &openapi.CreateMessageParams{}
	params.SetTo(sms.To)
	params.SetFrom(cfg.TwilioNumber)
	params.SetBody(messageBody)

	// Send the message
	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send OTP: %w", err)
	}

	// Log the response for debugging
	log.Printf("OTP sent successfully! SID: %s", *resp.Sid)
	return nil

}
