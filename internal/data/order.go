package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "zshop/api/order/v1"
	"zshop/internal/biz"
	"zshop/internal/types"
)

type orderRepo struct {
	data *Data
	log  *log.Helper
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *orderRepo) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (rsp *v1.CreateOrderReply, err error) {
	r.log.WithContext(ctx).Infof("creating order req:%v", req)
	db := r.data.db

	u := types.User{
		UserId:   req.GetUser().GetUserId(),
		UserName: req.GetUser().GetUserName(),
	}

	result := db.Create(&u)

	if result.RowsAffected != 1 {
		log.Errorf("insert table t_user error: param:%v", u)
		return nil, result.Error
	}

	tx := db.Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		return nil, err
	}

	// todo graceful solution
	order := types.Order{
		OrderId:       req.GetOrderInformation().GetOrderField().GetOrderId(),
		TransactionId: 0,
		ProductId:     0,
		ProductType:   0,
		Quantity:      0,
		Size:          "",
		Color:         "",
		Status:        0,
		RetryTime:     0,
	}

	ship := types.Shipping{
		Email:             "",
		Address:           "",
		FirstName:         "",
		LastName:          "",
		ApartmentSuiteEtc: "",
		City:              "",
		State:             0,
		ZipCode:           0,
		Phone:             0,
	}

	if err = tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Create(&ship).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if tx.Commit().Error != nil {
		return nil, err
	}

	//todo mq publish for other service
	return rsp, err
}
