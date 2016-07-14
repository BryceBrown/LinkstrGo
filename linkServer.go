package main

import (
	"fmt"
	"net"
    "net/http/fcgi"
	"net/http"
	"strings"
	"Laas"
	"strconv"
)

const ( 
	//BASE_URL = "http://www.golinkstr.com/"
	BASE_URL = "http://localhost:8080/"
	INTERSTITIAL_FRAME = `
		<html>
		<head>
		<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.1/jquery.min.js"></script>
		<link rel="stylesheet" type="text/css" href="//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap.min.css">
		<script type="text/javascript">
		var baseUrl = "{BASEURL}";
		var urlToRedirectTo = "{LOCATION}";
		var intervalTime = {TIME};
		var interstitialId = {INTERID};
		var linkId = {LINKID};
		var addUrl = "{AD_URL}";
		</script>
		<script src="{BASEURL}static/js/frontend/interstitial.js"></script>
		<link rel="stylesheet" type="text/css" href="{BASEURL}static/css/interstitial.css">
		</head>
		<body>
		<div style="height:100px;width:100%;"><div id='RedirectContent'><button class='btn' id='NextButton'>Continue</button>
		<br /><br /><div style='display: inline' id='RedirectCountDown'></div></div></div>
		<div id='ad_frame'></div>
		<iframe width="100%" height="100%" src="{FRAME}"></iframe>
		</body>
		</html>
	`
)

func main() {
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
	http.HandleFunc("/", handleRedirect)
	http.HandleFunc("/robots.txt", handleRobots)
 	listener, _ := net.Listen("tcp", "127.0.0.1:9001")
	fcgi.Serve(listener, nil)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
	logger := Laas.GAELogger{}
	logger.Init(r)

	dal := Laas.SQLDal{}
	err := dal.ConnToDb()
	if err != nil {
		logger.LogMessage("Error connecting to db - %v\n", Laas.LogLevelCritical, err)
	}

	redirUrl, interstitial, err := dal.GetRedirectUrl(r.Host, strings.TrimLeft(r.URL.Path, "/"))
	if err != nil {
		logger.LogMessage("Redirect URL is not valid - %v\n", Laas.LogLevelWarning, err)
	}

	//Grab data from request
	reqData := Laas.GetDataFromRequest(r)
	logger.LogMessage("IpAddress- %v\n", Laas.LogLevelDebug, reqData.IpAddress)
	go dal.SaveRedirectRequest(reqData)
	//send response
	logger.LogMessage("Redirect Url - %v\n", Laas.LogLevelDebug, redirUrl)
	if interstitial.Id != -1 {
		var resp string
		resp = strings.Replace(INTERSTITIAL_FRAME, "{LOCATION}", redirUrl, 1)
		resp = strings.Replace(resp, "{TIME}", "15000", 1)
		resp = strings.Replace(resp, "{FRAME}", interstitial.Url, 2)
		resp = strings.Replace(resp, "{LINKID}", strconv.Itoa(interstitial.LinkId), 1)
		resp = strings.Replace(resp, "{INTERID}", strconv.Itoa(interstitial.Id), 1)
		resp = strings.Replace(resp, "{AD_URL}", interstitial.AdUrl, 1)
		resp = strings.Replace(resp, "{BASEURL}", BASE_URL, 4)
		w.Write([]byte(resp))
	}else{
		http.Redirect(w, r, redirUrl, 302)
	}
}

func handleRobots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("User-agent: * \nDisallow: /"))
}