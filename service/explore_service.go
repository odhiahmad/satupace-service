package service

import (
	"math"
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/repository"

	"github.com/google/uuid"
)

type ExploreService interface {
	FindNearbyRunners(userId uuid.UUID, req request.ExploreRunnersRequest) ([]response.ExploreRunnerResponse, error)
	FindNearbyGroups(userId uuid.UUID, req request.ExploreGroupsRequest) ([]response.ExploreGroupResponse, error)
}

type exploreService struct {
	profileRepo     repository.RunnerProfileRepository
	groupRepo       repository.RunGroupRepository
	directMatchRepo repository.DirectMatchRepository
	memberRepo      repository.RunGroupMemberRepository
}

func NewExploreService(
	profileRepo repository.RunnerProfileRepository,
	groupRepo repository.RunGroupRepository,
	directMatchRepo repository.DirectMatchRepository,
	memberRepo repository.RunGroupMemberRepository,
) ExploreService {
	return &exploreService{
		profileRepo:     profileRepo,
		groupRepo:       groupRepo,
		directMatchRepo: directMatchRepo,
		memberRepo:      memberRepo,
	}
}

func (s *exploreService) FindNearbyRunners(userId uuid.UUID, req request.ExploreRunnersRequest) ([]response.ExploreRunnerResponse, error) {
	if req.RadiusKm <= 0 {
		req.RadiusKm = 10
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	profiles, err := s.profileRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Get existing matches if we want to exclude them
	matchedIds := make(map[uuid.UUID]bool)
	if req.ExcludeMatchedId {
		matches, _ := s.directMatchRepo.FindUserMatches(userId)
		for _, m := range matches {
			matchedIds[m.User1Id] = true
			matchedIds[m.User2Id] = true
		}
	}

	var results []response.ExploreRunnerResponse
	for _, p := range profiles {
		// Skip self
		if p.UserId == userId {
			continue
		}

		// Skip already matched
		if req.ExcludeMatchedId && matchedIds[p.UserId] {
			continue
		}

		// Skip inactive profiles
		if !p.IsActive {
			continue
		}

		// Calculate distance
		dist := haversine(req.Latitude, req.Longitude, p.Latitude, p.Longitude)
		if dist > req.RadiusKm {
			continue
		}

		// Pace filter
		if req.MinPace > 0 && p.AvgPace < req.MinPace {
			continue
		}
		if req.MaxPace > 0 && p.AvgPace > req.MaxPace {
			continue
		}

		// Preferred time filter
		if req.PreferredTime != "" && p.PreferredTime != req.PreferredTime {
			continue
		}

		// Gender filter (via User)
		if req.Gender != "" && p.User != nil && p.User.Gender != nil && *p.User.Gender != req.Gender {
			continue
		}

		// Women only filter
		if req.WomenOnly && !p.WomenOnlyMode {
			continue
		}

		var name *string
		var gender *string
		if p.User != nil {
			name = p.User.Name
			gender = p.User.Gender
		}

		results = append(results, response.ExploreRunnerResponse{
			UserId:            p.UserId.String(),
			Name:              name,
			Gender:            gender,
			AvgPace:           p.AvgPace,
			PreferredDistance: p.PreferredDistance,
			PreferredTime:     p.PreferredTime,
			Image:             p.Image,
			DistanceKm:        math.Round(dist*100) / 100,
			WomenOnlyMode:     p.WomenOnlyMode,
		})

		if len(results) >= req.Limit {
			break
		}
	}

	return results, nil
}

func (s *exploreService) FindNearbyGroups(userId uuid.UUID, req request.ExploreGroupsRequest) ([]response.ExploreGroupResponse, error) {
	if req.RadiusKm <= 0 {
		req.RadiusKm = 10
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Status == "" {
		req.Status = "open"
	}

	// Get groups the user has already joined or created
	memberships, _ := s.memberRepo.FindByUserId(userId)
	joinedGroupIds := make(map[uuid.UUID]bool)
	for _, m := range memberships {
		if m.Status == "joined" {
			joinedGroupIds[m.GroupId] = true
		}
	}

	groups, err := s.groupRepo.FindByStatus(req.Status)
	if err != nil {
		return nil, err
	}

	var results []response.ExploreGroupResponse
	for _, g := range groups {
		// Skip groups the user has already joined
		if joinedGroupIds[g.Id] {
			continue
		}

		dist := haversine(req.Latitude, req.Longitude, g.Latitude, g.Longitude)
		if dist > req.RadiusKm {
			continue
		}

		if req.MinPace > 0 && g.MaxPace < req.MinPace {
			continue
		}
		if req.MaxPace > 0 && g.MinPace > req.MaxPace {
			continue
		}

		if req.WomenOnly && !g.IsWomenOnly {
			continue
		}

		memberCount, _ := s.groupRepo.GetMemberCount(g.Id)

		results = append(results, response.ExploreGroupResponse{
			GroupId:           g.Id.String(),
			Name:              g.Name,
			MinPace:           g.MinPace,
			MaxPace:           g.MaxPace,
			PreferredDistance: g.PreferredDistance,
			MeetingPoint:      g.MeetingPoint,
			ScheduledAt:       g.ScheduledAt,
			MaxMember:         g.MaxMember,
			CurrentMembers:    int(memberCount),
			IsWomenOnly:       g.IsWomenOnly,
			Status:            g.Status,
			DistanceKm:        math.Round(dist*100) / 100,
			CreatedBy:         g.CreatedBy.String(),
		})

		if len(results) >= req.Limit {
			break
		}
	}

	return results, nil
}

// haversine calculates the distance in km between two lat/lng points.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // Earth radius in km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
