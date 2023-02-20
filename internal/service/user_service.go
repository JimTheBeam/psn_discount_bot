package service

import (
	"psn_discount_bot/internal/model"
)

func (s *Service) CreateUser(newUser model.User) {
	if _, err := s.repo.UpsertUser(newUser); err != nil {
		s.log.WithError(err).WithField("user_telegram_id", newUser.TelegramID).Error("get upsert user")

		return
	}
}
