package main

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/subscription/mgmt/subscription"
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var (
	subscriptions         = map[string]string{}
	secureScoreClients    = map[string]security.SecureScoresClient{}
	secureScorePointGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "azure_security_center",
		Name:      "secure_score_point",
		Help:      "Azure Security Center Secure Score as point",
	}, []string{"subscription_id"})
	secureScorePercentageGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "azure_security_center",
		Name:      "secure_score_percentage",
		Help:      "Azure Security Center Secure Score as percentage",
	}, []string{"subscription_id"})
)

func init() {
	// Register prometheus metrics
	prometheus.MustRegister(secureScorePointGauge)
	prometheus.MustRegister(secureScorePercentageGauge)

	//Authorize for SDK
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Fatalln(err)
	}

	// List subscriptions
	subscriptionClient := subscription.NewSubscriptionsClient()
	subscriptionClient.Authorizer = authorizer
	resultIterator, err := subscriptionClient.ListComplete(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	for resultIterator.NotDone() {
		sub := resultIterator.Value()
		log.Printf("subscriptionId=%v displayName=%v\n", *sub.SubscriptionID, *sub.DisplayName)
		subscriptions[*sub.SubscriptionID] = *sub.DisplayName
		err := resultIterator.NextWithContext(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Create map of SecureScoreClient
	for subscriptionId, _ := range subscriptions {
		secureScoreClient := security.NewSecureScoresClient(subscriptionId, "")
		secureScoreClient.Authorizer = authorizer
		secureScoreClients[subscriptionId] = secureScoreClient
	}
}

func refreshMetrics(secureScoreClient security.SecureScoresClient) {
	subscriptionName := subscriptions[secureScoreClient.SubscriptionID]
	secureScoresListPage, err := secureScoreClient.List(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	for secureScoresListPage.NotDone() {
		secureScoreItems := secureScoresListPage.Values()

		for _, secureScoreItem := range secureScoreItems {
			log.Printf("Subscription=%v Current=%v Percentage=%v\n", subscriptionName, *secureScoreItem.ScoreDetails.Current, *secureScoreItem.ScoreDetails.Percentage)
			gauge, err := secureScorePercentageGauge.GetMetricWithLabelValues(subscriptionName)
			if err != nil {
				log.Fatalln(err)
			}
			gauge.Set(*secureScoreItem.ScoreDetails.Percentage)

			gauge, err = secureScorePointGauge.GetMetricWithLabelValues(subscriptionName)
			if err != nil {
				log.Fatalln(err)
			}
			gauge.Set(*secureScoreItem.ScoreDetails.Current)
		}

		err := secureScoresListPage.NextWithContext(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func main() {
	for _, secureScoreClient := range secureScoreClients {
		go func(secureScoreClient security.SecureScoresClient) {
			for range time.NewTicker(5 * time.Minute).C {
				refreshMetrics(secureScoreClient)
			}
		}(secureScoreClient)
		refreshMetrics(secureScoreClient)
	}

	// Expose prometheus handler
	log.Println("Listen on :8080")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))

}
