package main

import (
	"fmt"
	"log"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/analytics/v3"
)

var (
	metric string = "rt:activeUsers" // GA metric that we want, rt from Real Time
)

// Google Analytics access stuff
var (
	// CHANGE THESE!!!
	gaServiceAcctEmail string = "XXXXXXXX-compute@developer.gserviceaccount.com" // (json:"client_email") email address of registered application
	gaTableID          string = "ga:XXXXXXXX"                                       // namespaced profile (table) ID of your analytics account/property/profile
	tokenurl           string = "https://accounts.google.com/o/oauth2/token"         // (json:"token_uri") Google oauth2 Token URL
)

func main() {

	pk := "-----BEGIN PRIVATE KEY-----\nXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\n-----END PRIVATE KEY-----\n"

	// create a jwt.Config that we will subsequently use for our authenticated client/transport
	// relevant docs for all the oauth2 & json web token stuff at https://godoc.org/golang.org/x/oauth2 & https://godoc.org/golang.org/x/oauth2/jwt
	
	jwtc := jwt.Config{
		Email:      gaServiceAcctEmail,
		PrivateKey: []byte(pk),
		Scopes:     []string{analytics.AnalyticsReadonlyScope},
		TokenURL:   tokenurl,
	}

	// create our authenticated http client using the jwt.Config we just created
	clt := jwtc.Client(oauth2.NoContext)
	// create a new analytics service by passing in the authenticated http client
	as, err := analytics.New(clt)
	if err != nil {
		log.Fatal("Error creating Analytics Service at analytics.New() -", err)
	}

	rt := analytics.NewDataRealtimeService(as)
	rtSetup := rt.Get(gaTableID, metric)

	//rtSetup.Dimensions("rt:medium")
	rtSetup.Dimensions("rt:pageTitle,rt:pagePath")
	rtSetup.Sort("-" + metric)
	rtSetup.MaxResults(20)

	gtadata, err := rtSetup.Do()
	if err != nil {
		log.Fatal("Could not get real time data:", err)
	}

	//ui.UseTheme("helloworld")
	rtMostViewed := make([]string, 20)
	for i, data := range gtadata.Rows {
		if data[1] != "/" {
			rtMostViewed[i] = "[" + data[2] + "] " + data[0]
		}
	}

	fmt.Printf("%#v\n", rtMostViewed)

	rtTrafficType := rt.Get(gaTableID, metric)
	rtTrafficType.Dimensions("rt:trafficType")
	gtadata, err = rtTrafficType.Do()

	if err != nil {
		log.Fatal("Could not get real time data:", err)
	}

	// TOTAL active users
	rtTotalActiveUsers, _ := strconv.Atoi(gtadata.TotalsForAllResults["rt:activeUsers"])
	fmt.Println("Active Users:", rtTotalActiveUsers)

	trafficType := map[string]int{
		"CUSTOM":   0,
		"DIRECT":   0,
		"ORGANIC":  0,
		"REFERRAL": 0,
		"SOCIAL":   0,
	}
	for _, data := range gtadata.Rows {
		percentage, _ := strconv.Atoi(data[1])
		trafficType[data[0]] = percentage * 100 / rtTotalActiveUsers
	}

	fmt.Printf("%#v\n", trafficType)

}

