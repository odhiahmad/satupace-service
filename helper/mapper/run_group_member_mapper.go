package mapper

import (
	"run-sync/data/response"
	"run-sync/entity"
)

func MapRunGroupMember(m *entity.RunGroupMember) *response.RunGroupMemberResponse {
	if m == nil {
		return nil
	}

	return &response.RunGroupMemberResponse{
		Id:       m.Id.String(),
		GroupId:  m.GroupId.String(),
		UserId:   m.UserId.String(),
		Status:   m.Status,
		JoinedAt: m.JoinedAt,
	}
}

func MapRunGroupMemberDetail(m *entity.RunGroupMember, user *entity.User) *response.RunGroupMemberDetailResponse {
	if m == nil {
		return nil
	}

	return &response.RunGroupMemberDetailResponse{
		Id:       m.Id.String(),
		GroupId:  m.GroupId.String(),
		UserId:   m.UserId.String(),
		User:     MapUser(user),
		Status:   m.Status,
		JoinedAt: m.JoinedAt,
	}
}
