/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package model

import "github.com/jinzhu/gorm"

type Video struct {
	gorm.Model
	Title         string `json:"title"`
	AuthorId      uint   `json:"author_id"` // TODO: 这里的外键约束交由 Service 层保证
	PlayUrl       string `json:"play_url"`  // 播放外链
	CoverUrl      string `json:"cover_url"` // 封面外链
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
}
