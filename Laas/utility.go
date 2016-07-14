package Laas

import (
	"net/http"
	"strings"
	"strconv"
	"net"
)

func GetDataFromRequest(r *http.Request) *RedirectRequest{
	reqData := RedirectRequest{}
	reqData.Host = r.Host
	reqData.Path = strings.TrimLeft(r.URL.Path, "/")
	langs, ok := r.Header["Accept-Language"]
	if ok {
		reqData.Languages = langs
	}
	aTypes, ok := r.Header["User-Agent"]
	if ok {
		reqData.AgentTypes = aTypes
	}
	referer, ok := r.Header["Referer"]
	if ok && len(referer) > 0 {
		reqData.RefererUrl = referer[0]
	}
	colIndex := strings.IndexAny(r.RemoteAddr, ":")
	if colIndex > -1{
		reqData.IpAddress = r.RemoteAddr[:colIndex]
	}else{
		reqData.IpAddress = r.RemoteAddr
	}
	return &reqData
}

// Convert net.IP to int64
func inet_aton(ipnr net.IP) int64 {      
    bits := strings.Split(ipnr.String(), ".")
    
    b0, _ := strconv.Atoi(bits[0])
    b1, _ := strconv.Atoi(bits[1])
    b2, _ := strconv.Atoi(bits[2])
    b3, _ := strconv.Atoi(bits[3])

    var sum int64
    
    sum += int64(b0) << 24
    sum += int64(b1) << 16
    sum += int64(b2) << 8
    sum += int64(b3)
    
    return sum
}