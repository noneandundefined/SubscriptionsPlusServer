package main

import (
	"context"
	"time"
)

func (s *httpServer) startCronUpdSubscriptions() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if err := s.store.Subscriptions.Update_SubscriptionsMounth(context.Background()); err != nil {
				s.logger.Error("%v", err)
			}
		}
	}()
}

func (s *httpServer) startCronUpdUserSubscriptionBeforEndSub() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if err := s.store.Users.Update_UserSubscriptionBeforEndSub(context.Background()); err != nil {
				s.logger.Error("%v", err)
			}
		}
	}()
}

func (s *httpServer) startCronAutoActivateExpiredTransactions() {
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if err := s.store.Transactions.AutoActivateExpiredTransactions(context.Background()); err != nil {
				s.logger.Error("%v", err)
			}
		} 
	}()
}
