package Laas

import (
	"database/sql"
	"fmt"
	"errors"
	"net"
	"strconv"
	"user_agent"
	"math/rand"
	_ "github.com/go-sql-driver/mysql"
)

const (
	POSSIBLE_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	creds = "<INSERT SQL SERVER CREDENTIALS>"
)

type SQLDal struct {
	sqlConn *sql.DB
}

func  (dal *SQLDal) ConnToDb() error {
	var err error
	dal.sqlConn, err = sql.Open("mysql", creds)
	if err != nil {
		fmt.Println("Error opening database connection")
		return err
	}

	err = dal.sqlConn.Ping()
	if err != nil {
		fmt.Println("Error opening database connection")
		return err
	}

	//If here the connection is successful
	return nil
}

func (dal *SQLDal) Close() {
	dal.sqlConn.Close()
}

func (dal *SQLDal) SaveRedirectRequest(req *RedirectRequest) error{
	//If there are no agent types, the types dont begin with mozilla or opera, or they contain bot, ignore the request
	defer dal.Close()
	if len(req.AgentTypes) == 0 {
		return nil
	}

	agent := user_agent.UserAgent{}
	agent.Parse(req.AgentTypes[0])
	if agent.Bot(){
		return nil
	}


	linkId, err := dal.GetLinkId(req.Path, req.Host)
	if err != nil {
		fmt.Print(err)
		return err
	}

	ipNumber := inet_aton(net.ParseIP(req.IpAddress))
	country, code, err := dal.GetCountry(ipNumber)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("LinkID: " + strconv.Itoa(linkId))
	resp, err := dal.sqlConn.Exec(`INSERT INTO FrontEnd_linkstat (Link_id,Referer,IpAddress,CountryCode,Country,TimeClicked) 
											VALUES(?,?,?,?,?,NOW())`,
											linkId, req.RefererUrl, req.IpAddress, code, country)

	if err != nil {
		fmt.Print(err);
		return err
	}else{

		//Now save all of the Agent Types and languages
		statId, err := resp.LastInsertId()
		if err != nil{
			//There was an error somewhere, exit
			return err
		}

		//increment the link total
		err = dal.IncrementTotal(linkId)
		if err != nil {
			return err
		}

		err = dal.saveAgentTypes(&agent, statId, req.AgentTypes[0])
		if err != nil {
			return err
		}

		//Last call, just return
		return dal.saveLanguages(req, statId)
	}
}

func (dal *SQLDal) IncrementTotal(linkId int) error{
	row, err := dal.sqlConn.Query(`SELECT id FROM FrontEnd_linkclicktotal WHERE Link_id=? AND Date=CURDATE()`, linkId)
	if row != nil && row.Next() {
		fmt.Print("Got linkId for click total")
		//Update Old entry
		var statId int
		err = row.Scan(&statId)
		if err != nil {
			return err
		}
		_ , err = dal.sqlConn.Exec(`UPDATE FrontEnd_linkclicktotal SET TotalClicked=TotalClicked + 1 WHERE id=?`, statId)
		return err
	} else {
		fmt.Print("No linkId found for click total")
		_ , err := dal.sqlConn.Exec(`INSERT INTO FrontEnd_linkclicktotal(Link_id,TotalClicked,Date) VALUES(?,1,CURDATE())`, linkId)
		return err
	}

}

func (dal *SQLDal) GetCountry(ipNumber int64) (string, string, error){
	rows, err := dal.sqlConn.Query(`SELECT Country, CountryCode FROM linkstr.FrontEnd_lightlocationinfo WHERE IpNumberStart <= ? AND IpNumberEnd >= ?`, ipNumber, ipNumber)
	if err != nil {
		return "", "", err
	}

	if rows.Next() {
		var country, code string
		err = rows.Scan(&country, &code)
		if err != nil {
			return	"", "", err
		}
		return country, code, nil
	}
	return "", "", errors.New("No IpAddress Lookup found")
}

func (dal *SQLDal) GetLinkId(path string, host string) (int, error) {
	var linkId int
	linkId = -1
	err := dal.sqlConn.QueryRow(`SELECT FrontEnd_redirectlink.id FROM FrontEnd_redirectlink LEFT JOIN 
											FrontEnd_supporteddomain ON FrontEnd_redirectlink.Domain_Id=FrontEnd_supporteddomain.Id 
											WHERE Domain LIKE ? AND UrlKey LIKE ?`, host, path).Scan(&linkId)

	return linkId, err
}

func (dal *SQLDal) GetRedirectUrl(host string, path string) (string, Interstitial, error){
	fmt.Println("Path : " + host)
	var redirUrl string
	var id int
	err := dal.sqlConn.QueryRow(`SELECT RedirectUrl, FrontEnd_redirectlink.id FROM FrontEnd_redirectlink LEFT JOIN 
											FrontEnd_supporteddomain ON FrontEnd_redirectlink.Domain_Id=FrontEnd_supporteddomain.id 
											WHERE Domain LIKE ? AND UrlKey LIKE ?`, host, path).Scan(&redirUrl, &id)
	if err != nil {
		return "http://www.golinkstr.com/Redirect404", Interstitial{Id:-1},  err
	}

	intersticial, err := dal.getDomainInterstitial(host)
	intersticial.LinkId = id

	return redirUrl, intersticial, err
}

func (dal *SQLDal) LogMessage(message string){
	fmt.Println(message)
}

func (dal *SQLDal) saveAgentTypes(agent *user_agent.UserAgent, linkClickId int64, uaString string) error {
	browser, version := agent.Browser()
	browserString := browser + "/" + version
	_, err := dal.sqlConn.Exec("INSERT INTO FrontEnd_linkagenttype(Stat_Id,AgentType,Browser,OS,Device) VALUES(?, ?, ?, ?, ?)", 
												linkClickId, uaString, browserString, agent.OS(), agent.Platform())

	return err
}

func (dal *SQLDal) saveLanguages(req *RedirectRequest, linkClickId int64) error {
	for _, value := range req.Languages {
		_ , err := dal.sqlConn.Exec("INSERT INTO FrontEnd_linklanguage(Stat_Id,Language) VALUES(?, ?)", linkClickId, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dal *SQLDal) getDomainInterstitial(host string) (Interstitial, error) {
	url := ""
	adUrl := ""
	var (
		displayChance int32
		interstitialId int
	)
	interstitialId = -1
	displayChance = 0
	err := dal.sqlConn.QueryRow(`SELECT Url, DisplayChance, FrontEnd_intersticial.id, AdClickUrl FROM FrontEnd_supporteddomain
							    LEFT JOIN FrontEnd_intersticial ON FrontEnd_supporteddomain.intersticial_id = FrontEnd_intersticial.id
								where FrontEnd_supporteddomain.Domain LIKE ?`, host).Scan(&url, &displayChance, &interstitialId, &adUrl)
	if err != nil {
		return Interstitial{Id:-1}, err
	}
	if url != "" && displayChance != 0 {
		numb := rand.Int31n(100) + 1
		if numb <= displayChance {
			return Interstitial{url, interstitialId, -1, adUrl} , nil
		}
	}
	return Interstitial{Id:-1}, nil
}
