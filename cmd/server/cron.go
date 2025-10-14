package main

import (
	"context"
	"fmt"
	"subscriptionplus/server/pkg/notify"
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

func (s *httpServer) startCronSubscriptionNotifications() {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			duration := next.Sub(now)
			time.Sleep(duration)

			ctx := context.Background()

			subs, err := s.store.Subscriptions.Get_SubscriptionsForNotify(ctx)
			if err != nil {
				s.logger.Error("failed to fetch subscriptions: %v", err)

				time.Sleep(1 * time.Hour)
				continue
			}

			if subs == nil || len(*subs) == 0 {
				continue
			}

			for _, sub := range *subs {
				token, err := s.store.Notifications.Get_NotificationTokenByUuid(ctx, sub.UserUUID)
				if err != nil {
					s.logger.Error("failed to get token for user %s: %v", sub.UserUUID, err)
					continue
				}

				if token != nil {
					_ = notify.SendPushNotification(token.Token, fmt.Sprintf("Upcoming subscription: %s", sub.Name), fmt.Sprintf("Your plan will renew for %.2f RUB on %s", sub.Price, sub.DatePay.Format("January 2, 2006")))
				}
			}
		}
	}()
}
