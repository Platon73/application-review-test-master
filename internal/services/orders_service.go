package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"example/internal/dto"
	"example/internal/interfaces"
	"example/internal/repositories"
)

var (
	ErrEmptyUserID      = errors.New("user_id is required")
	ErrInvalidUser      = errors.New("invalid user")
	ErrRoomNotAvailable = errors.New("room not available")
	ErrInvalidDates     = errors.New("invalid date range")
)

type OrdersService struct {
	repo        interfaces.OrdersRepository
	userRepo    interfaces.UserRepository
	emailSender interfaces.EmailSender
	logger      interfaces.Logger
}

func NewOrdersService(
	repo interfaces.OrdersRepository,
	userRepo interfaces.UserRepository,
	emailSender interfaces.EmailSender,
	logger interfaces.Logger,
) *OrdersService {
	return &OrdersService{
		repo:        repo,
		userRepo:    userRepo,
		emailSender: emailSender,
		logger:      logger,
	}
}

func (s *OrdersService) Create(order *dto.Order) (*dto.Order, error) {
	// Валидация входных данных
	if order.UserID == "" {
		return nil, ErrEmptyUserID
	}
	if order.HotelID == "" || order.RoomTypeID == "" {
		return nil, errors.New("hotel_id and room_type_id are required")
	}
	if order.From.IsZero() || order.To.IsZero() {
		return nil, errors.New("dates are required")
	}
	if order.From.After(order.To) {
		return nil, ErrInvalidDates
	}

	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Получение пользователя
	user, err := s.userRepo.GetByID(ctx, order.UserID)
	if errors.Is(err, repositories.ErrUserNotFound) {
		return nil, fmt.Errorf("%w: user_id=%s", ErrInvalidUser, order.UserID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	order.User = user

	// Асинхронная отправка email с защитой от паники
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("panic in email sending",
					"recover", r,
					"user_id", user.ID)
			}
		}()

		content := fmt.Sprintf(
			"Уважаемый %s %s! Бронирование подтверждено.\nОтель: %s\nНомер: %s\nДаты: %s - %s",
			user.Firstname,
			user.Lastname,
			order.HotelID,
			order.RoomTypeID,
			order.From.Format("02.01.2006"),
			order.To.Format("02.01.2006"),
		)

		if err := s.emailSender.SendConfirmation(user.Email, content); err != nil {
			s.logger.Error("email sending failed",
				"error", err,
				"email", user.Email,
				"user_id", user.ID)
		}
	}()

	// Создание заказа
	err = s.repo.Create(ctx, order)
	if err != nil {
		if errors.Is(err, repositories.ErrRoomNotAvailable) {
			return nil, ErrRoomNotAvailable
		}
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}
