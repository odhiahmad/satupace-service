package service

import (
	"math"
	"run-sync/entity"
	"run-sync/repository"

	"github.com/google/uuid"
)

// ── Matching Configuration ──

const (
	DefaultMaxRadius      = 10.0 // km
	DefaultPaceTolerance  = 2.0  // min/km
	PaceWeight            = 0.7
	LocationWeight        = 0.3
	MinCompatibilityScore = 0.3
)

// ── Interface ──

type MatchingEngine interface {
	FindDirectCandidates(userId uuid.UUID) ([]CandidateResult, error)
	FindGroupCandidates(userId uuid.UUID) ([]GroupCandidateResult, error)
	CalculateCompatibility(profileA, profileB entity.RunnerProfile) float64
}

type CandidateResult struct {
	Profile       entity.RunnerProfile
	User          *entity.User
	Compatibility float64
	DistanceKm    float64
}

type GroupCandidateResult struct {
	Group         entity.RunGroup
	Compatibility float64
	DistanceKm    float64
	MemberCount   int64
}

// ── Implementation ──

type matchingEngine struct {
	profileRepo repository.RunnerProfileRepository
	matchRepo   repository.DirectMatchRepository
	groupRepo   repository.RunGroupRepository
	safetyRepo  repository.SafetyLogRepository
}

func NewMatchingEngine(
	profileRepo repository.RunnerProfileRepository,
	matchRepo repository.DirectMatchRepository,
	groupRepo repository.RunGroupRepository,
	safetyRepo repository.SafetyLogRepository,
) MatchingEngine {
	return &matchingEngine{
		profileRepo: profileRepo,
		matchRepo:   matchRepo,
		groupRepo:   groupRepo,
		safetyRepo:  safetyRepo,
	}
}

// FindDirectCandidates returns nearby compatible runners, excluding already
// matched/blocked users, sorted by compatibility score descending.
func (e *matchingEngine) FindDirectCandidates(userId uuid.UUID) ([]CandidateResult, error) {
	myProfile, err := e.profileRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	allProfiles, err := e.profileRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Build exclusion set: already matched or blocked
	excludeSet := map[uuid.UUID]bool{userId: true}
	matches, _ := e.matchRepo.FindUserMatches(userId)
	for _, m := range matches {
		excludeSet[m.User1Id] = true
		excludeSet[m.User2Id] = true
	}

	// Blocked users (from safety reports by this user)
	blockedLogs, _ := e.safetyRepo.FindByUserId(userId)
	for _, l := range blockedLogs {
		if l.Status == "blocked" {
			excludeSet[l.MatchId] = true // note: MatchId here refers to the reported user
		}
	}

	var candidates []CandidateResult
	for _, p := range allProfiles {
		// Skip self and excluded
		if excludeSet[p.UserId] || !p.IsActive {
			continue
		}

		// Women-only check
		if myProfile.WomenOnlyMode && !p.WomenOnlyMode {
			continue
		}

		// Distance check
		dist := Haversine(myProfile.Latitude, myProfile.Longitude, p.Latitude, p.Longitude)
		if dist > DefaultMaxRadius {
			continue
		}

		// Compatibility score
		score := e.CalculateCompatibility(*myProfile, p)
		if score < MinCompatibilityScore {
			continue
		}

		candidates = append(candidates, CandidateResult{
			Profile:       p,
			User:          p.User,
			Compatibility: math.Round(score*100) / 100,
			DistanceKm:    math.Round(dist*100) / 100,
		})
	}

	// Sort by compatibility descending
	sortCandidates(candidates)

	return candidates, nil
}

// FindGroupCandidates returns nearby open groups compatible with user's profile.
func (e *matchingEngine) FindGroupCandidates(userId uuid.UUID) ([]GroupCandidateResult, error) {
	myProfile, err := e.profileRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	groups, err := e.groupRepo.FindByStatus("open")
	if err != nil {
		return nil, err
	}

	var candidates []GroupCandidateResult
	for _, g := range groups {
		// Women-only check
		if myProfile.WomenOnlyMode && !g.IsWomenOnly {
			continue
		}
		if g.IsWomenOnly && !myProfile.WomenOnlyMode {
			continue
		}

		dist := Haversine(myProfile.Latitude, myProfile.Longitude, g.Latitude, g.Longitude)
		if dist > DefaultMaxRadius {
			continue
		}

		memberCount, _ := e.groupRepo.GetMemberCount(g.Id)
		if g.MaxMember > 0 && int(memberCount) >= g.MaxMember {
			continue // already full
		}

		score := calculateGroupCompatibility(*myProfile, g)
		if score < MinCompatibilityScore {
			continue
		}

		candidates = append(candidates, GroupCandidateResult{
			Group:         g,
			Compatibility: math.Round(score*100) / 100,
			DistanceKm:    math.Round(dist*100) / 100,
			MemberCount:   memberCount,
		})
	}

	sortGroupCandidates(candidates)

	return candidates, nil
}

// CalculateCompatibility computes a 0–1 score between two runner profiles.
//
//	compatibility = (paceScore * 0.7) + (locationScore * 0.3)
//	paceScore     = 1 - (abs(paceA - paceB) / maxTolerance)
//	locationScore = 1 - (distance / maxRadius)
func (e *matchingEngine) CalculateCompatibility(a, b entity.RunnerProfile) float64 {
	paceDiff := math.Abs(a.AvgPace - b.AvgPace)
	paceScore := 1.0 - (paceDiff / DefaultPaceTolerance)
	if paceScore < 0 {
		paceScore = 0
	}

	dist := Haversine(a.Latitude, a.Longitude, b.Latitude, b.Longitude)
	locationScore := 1.0 - (dist / DefaultMaxRadius)
	if locationScore < 0 {
		locationScore = 0
	}

	return (paceScore * PaceWeight) + (locationScore * LocationWeight)
}

func calculateGroupCompatibility(profile entity.RunnerProfile, group entity.RunGroup) float64 {
	// Compare profile pace against group's pace range midpoint
	groupMidPace := (group.MinPace + group.MaxPace) / 2.0
	paceDiff := math.Abs(profile.AvgPace - groupMidPace)
	paceScore := 1.0 - (paceDiff / DefaultPaceTolerance)
	if paceScore < 0 {
		paceScore = 0
	}

	dist := Haversine(profile.Latitude, profile.Longitude, group.Latitude, group.Longitude)
	locationScore := 1.0 - (dist / DefaultMaxRadius)
	if locationScore < 0 {
		locationScore = 0
	}

	return (paceScore * PaceWeight) + (locationScore * LocationWeight)
}

// ── Sorting helpers ──

func sortCandidates(c []CandidateResult) {
	for i := 1; i < len(c); i++ {
		for j := i; j > 0 && c[j].Compatibility > c[j-1].Compatibility; j-- {
			c[j], c[j-1] = c[j-1], c[j]
		}
	}
}

func sortGroupCandidates(c []GroupCandidateResult) {
	for i := 1; i < len(c); i++ {
		for j := i; j > 0 && c[j].Compatibility > c[j-1].Compatibility; j-- {
			c[j], c[j-1] = c[j-1], c[j]
		}
	}
}

// Haversine calculates the distance in km between two lat/lng points.
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
