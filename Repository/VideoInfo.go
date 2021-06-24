package repository

import (
	"context"
	model "homeplay/Model"
	"time"
)

func InsertInfos(d *model.VideoInfo) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ctxDB := model.DB.WithContext(ctx)
	tx := ctxDB.Begin()
	if err := tx.Where("videoNo=?", d.VideoNo).WithContext(ctx).Delete(&model.VideoInfo{}).Error; err != nil {
		// 返回任何错误都会回滚事务
		tx.Rollback()
		return err
	}

	if err := tx.Create(d).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 返回 nil 提交事务
	return tx.Commit().Error
}
