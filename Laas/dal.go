package Laas

import (
	"time"
	"database/sql"
)

//TODO Make Agreegate stats in memory. We might want to move this to a background thing
type LinkStat struct {
	DomainId int32
	RedirectUrl string
	UrlKey string
	RefererUrl string
	ClickTime time.Time
}

type OverallStats struct {
	TotalLinkClicks int32
	LanguageStats map[string]int32
	AgentStats map[string]int32
}

type RedirectLink struct {
	RedirectUrl string
	UrlKey string
}

type Domain struct {
	DomainId int32
	Name string
	Links []RedirectLink
}

type ErrorMessage struct {
	Message string
}

type RedirectRequest struct {
	RefererUrl string
	Languages []string
	AgentTypes []string
	Host string
	Path string
	IpAddress string
}

type SocialInfo struct {
	Url string
	FacebookId string
}

type AuthResponse struct {
	Message string
	Token string
}

type LinkResponse struct {
	Message string
	Url string
}

type Interstitial struct {
	Url string
	Id int
	LinkId int
	AdUrl string
}

type SqlAbstraction interface {
	Initialize(creds string) error
	Exec(sql string) (sql.Result, error)
	QueryRow(sql string) (sql.Rows, error)
}

type LaasLogger interface {
	GetLogFilePath() string
	WriteLog(format string, args ...interface{})
}

type DAL interface {
	GetRedirectUrl(host string, path string) (string, error)
	SaveRedirectRequest(req *RedirectRequest) error
	SaveSocialInfo(info *SocialInfo) error
	ConnToDb(creds string) error
	GenerateUrl(refUrl string, host string, authToken string) (string, error)
	Login(username string, password string) (string, error)
	LogMessage(message string)

	GetLinkStatsByDate(domain string, urlKey string, startTime time.Time, endTime time.Time) ([]LinkStat, error)
	GetLinks(page int, pagecount int, domain string) ([]RedirectLink, error)
}