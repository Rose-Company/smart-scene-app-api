package server

import (
	"smart-scene-app-api/common"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ACList struct {
	ActionId string         `json:"action_id"`
	RoleIds  pq.StringArray `json:"role_ids" gorm:"type:[]text"`
	UserID   pq.StringArray `json:"user_id" gorm:"type:[]text"`
}

type AuthorizationConfig map[string][]string

func (r AuthorizationConfig) CheckValidValidRole(roleID string, userID string, actionId string) error {
	auData, ok := r[actionId]
	if !ok {
		return common.ErrActionNotAllowed
	}
	for _, role := range auData {
		if role == roleID || role == userID {
			return nil
		}
	}
	return common.ErrActionNotAllowed
}

func (s *server) InitAuthorizationData() {
	db := s.GetService(common.PREFIX_MAIN_POSTGRES).(*gorm.DB)
	var acList []ACList
	ad := AuthorizationConfig{}
	err := db.Raw(`
	 SELECT 
		action_id,
		array_agg(role_id) FILTER (WHERE role_id IS NOT NULL) as role_ids,
		array_agg(user_id) FILTER (WHERE user_id IS NOT NULL) as user_ids
	FROM 
		PUBLIC.action_control_list
	WHERE 
		status = 1
	GROUP BY 
		action_id
	`).Scan(&acList).Error
	if err != nil {
		panic("Fetch ACL data error")
	}

	for _, action := range acList {
		ad[action.ActionId] = append(action.RoleIds, action.UserID...)
	}
	s.authorization = ad
}
