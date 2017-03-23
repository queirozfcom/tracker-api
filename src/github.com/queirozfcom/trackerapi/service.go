package trackerapi

import (
	"context"
	"errors"
	"sync"
	"github.com/google/go-github/github"
	"fmt"
)

// Service is a simple CRUD interface for user profiles.
type Service interface {
	GetWatchedRepos(ctx context.Context, username string) ([]RepoInformation, error)
	GetStarredRepos(ctx context.Context, username string) ([]RepoInformation, error)
	//GetMyWatchedRepos(ctx context.Context) ([]interface{}, error)
	//PostProfile(ctx context.Context, p Profile) error
	//GetProfile(ctx context.Context, id string) (Profile, error)
	//PutProfile(ctx context.Context, id string, p Profile) error
	//PatchProfile(ctx context.Context, id string, p Profile) error
	//DeleteProfile(ctx context.Context, id string) error
	//GetAddresses(ctx context.Context, profileID string) ([]Address, error)
	//GetAddress(ctx context.Context, profileID string, addressID string) (Address, error)
	//PostAddress(ctx context.Context, profileID string, a Address) error
	//DeleteAddress(ctx context.Context, profileID string, addressID string) error
}

type RepoInformation struct {
	FullName string
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type inmemService struct {
	mtx      sync.RWMutex
	m        map[string](map[string]string)
	ghClient github.Client
}

func NewInmemService(ghClient github.Client) Service {
	return &inmemService{
		m:        map[string](map[string]string){},
		ghClient: ghClient,
	}
}

func (s *inmemService) GetWatchedRepos(ctx context.Context, username string) ([]RepoInformation, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	opt := &github.ListOptions{PerPage: 10}

	var allRepos []RepoInformation

	for {

		repos, resp, err := s.ghClient.Activity.ListWatched(ctx, username, opt)

		if err != nil {
			return []RepoInformation{}, err
		}

		for _, element := range repos {
			allRepos = append(allRepos, RepoInformation{FullName: *element.FullName})
		}

		if resp.NextPage == 0 {
			fmt.Println(resp.Response.Header.Get("X-Ratelimit-Remaining"))

			break
		}
		opt.Page = resp.NextPage

	}


	return allRepos, nil

}

func (s *inmemService) GetMyStarredRepos(ctx context.Context) ([]RepoInformation, error) {
	return s.GetStarredRepos(ctx, "")
}

func (s *inmemService) GetStarredRepos(ctx context.Context, username string) ([]RepoInformation, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	opt := &github.ActivityListStarredOptions{}

	var allRepos []RepoInformation

	for {
		repos, resp, err := s.ghClient.Activity.ListStarred(ctx, username, opt)

		if err != nil {
			return []RepoInformation{}, err
		}

		for _, element := range repos {
			allRepos = append(allRepos, RepoInformation{FullName: *element.Repository.FullName})
		}

		if resp.NextPage == 0 {
			fmt.Println(resp.Response.Header.Get("X-Ratelimit-Remaining"))

			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil

}

func (s *inmemService) GetMyWatchedRepos(ctx context.Context) ([]RepoInformation, error) {
	return s.GetWatchedRepos(ctx, "")
}

//func (s *inmemService) PostProfile(ctx context.Context, p Profile) error {
//	s.mtx.Lock()
//	defer s.mtx.Unlock()
//	if _, ok := s.m[p.ID]; ok {
//		return ErrAlreadyExists // POST = create, don't overwrite
//	}
//	s.m[p.ID] = p
//	return nil
//}
//
//func (s *inmemService) GetProfile(ctx context.Context, id string) (Profile, error) {
//	s.mtx.RLock()
//	defer s.mtx.RUnlock()
//	p, ok := s.m[id]
//	if !ok {
//		return Profile{}, ErrNotFound
//	}
//	return p, nil
//}
//
//func (s *inmemService) PutProfile(ctx context.Context, id string, p Profile) error {
//	if id != p.ID {
//		return ErrInconsistentIDs
//	}
//	s.mtx.Lock()
//	defer s.mtx.Unlock()
//	s.m[id] = p // PUT = create or update
//	return nil
//}
//
//func (s *inmemService) PatchProfile(ctx context.Context, id string, p Profile) error {
//	if p.ID != "" && id != p.ID {
//		return ErrInconsistentIDs
//	}
//
//	s.mtx.Lock()
//	defer s.mtx.Unlock()
//
//	existing, ok := s.m[id]
//	if !ok {
//		return ErrNotFound // PATCH = update existing, don't create
//	}
//
//	// We assume that it's not possible to PATCH the ID, and that it's not
//	// possible to PATCH any field to its zero value. That is, the zero value
//	// means not specified. The way around this is to use e.g. Name *string in
//	// the Profile definition. But since this is just a demonstrative example,
//	// I'm leaving that out.
//
//	if p.Name != "" {
//		existing.Name = p.Name
//	}
//	if len(p.Addresses) > 0 {
//		existing.Addresses = p.Addresses
//	}
//	s.m[id] = existing
//	return nil
//}
//
//func (s *inmemService) DeleteProfile(ctx context.Context, id string) error {
//	s.mtx.Lock()
//	defer s.mtx.Unlock()
//	if _, ok := s.m[id]; !ok {
//		return ErrNotFound
//	}
//	delete(s.m, id)
//	return nil
//}
//
//func (s *inmemService) GetAddresses(ctx context.Context, profileID string) ([]Address, error) {
//	s.mtx.RLock()
//	defer s.mtx.RUnlock()
//	p, ok := s.m[profileID]
//	if !ok {
//		return []Address{}, ErrNotFound
//	}
//	return p.Addresses, nil
//}
//
//func (s *inmemService) GetAddress(ctx context.Context, profileID string, addressID string) (Address, error) {
//	s.mtx.RLock()
//	defer s.mtx.RUnlock()
//	p, ok := s.m[profileID]
//	if !ok {
//		return Address{}, ErrNotFound
//	}
//	for _, address := range p.Addresses {
//		if address.ID == addressID {
//			return address, nil
//		}
//	}
//	return Address{}, ErrNotFound
//}
//
//func (s *inmemService) PostAddress(ctx context.Context, profileID string, a Address) error {
//	s.mtx.Lock()
//	defer s.mtx.Unlock()
//	p, ok := s.m[profileID]
//	if !ok {
//		return ErrNotFound
//	}
//	for _, address := range p.Addresses {
//		if address.ID == a.ID {
//			return ErrAlreadyExists
//		}
//	}
//	p.Addresses = append(p.Addresses, a)
//	s.m[profileID] = p
//	return nil
//}
//
//func (s *inmemService) DeleteAddress(ctx context.Context, profileID string, addressID string) error {
//	s.mtx.Lock()
//	defer s.mtx.Unlock()
//	p, ok := s.m[profileID]
//	if !ok {
//		return ErrNotFound
//	}
//	newAddresses := make([]Address, 0, len(p.Addresses))
//	for _, address := range p.Addresses {
//		if address.ID == addressID {
//			continue // delete
//		}
//		newAddresses = append(newAddresses, address)
//	}
//	if len(newAddresses) == len(p.Addresses) {
//		return ErrNotFound
//	}
//	p.Addresses = newAddresses
//	s.m[profileID] = p
//	return nil
//}
