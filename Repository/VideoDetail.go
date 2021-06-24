package repository

import (
	"context"
	model "homeplay/Model"
	"time"
)

func InsertDetails(i []model.VideoDetail) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ctxDB := model.DB.WithContext(ctx)
	tx := ctxDB.Begin()
	if err := tx.Where("videoNo=?", i[0].VideoNo).Delete(&model.VideoDetail{}).Error; err != nil {
		// 返回任何错误都会回滚事务
		tx.Rollback()
		return err
	}

	if err := tx.CreateInBatches(i, len(i)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 返回 nil 提交事务
	return tx.Commit().Error
}
