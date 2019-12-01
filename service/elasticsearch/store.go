package elasticsearch

import (
	"context"

	pb "github.com/meateam/search-service/proto"
	es "github.com/olivere/elastic/v7"
)

// Store holds the elasticsearch and implements Store interface.
type Store struct {
	client *es.Client
	index  string
}

func newStore(cfg []es.ClientOptionFunc, index string) (*Store, error) {
	client, err := es.NewClient(cfg...)
	if err != nil {
		return nil, err
	}

	return &Store{client: client, index: index}, nil
}

// HealthCheck checks the health of the service, returns true if healthy, or false otherwise.
func (s Store) HealthCheck(ctx context.Context) (bool, error) {
	_, err := s.client.CatHealth().Do(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetAll finds all files that matches the query and Index,
// if successful returns a file slice, and a nil error,
// otherwise returns nil and non-nil error if any occured.
func (s Store) GetAll(ctx context.Context, query es.Query) ([]string, error) {
	res, err := s.client.Search().
		Index(s.index).
		Query(query).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, res.TotalHits())
	for _, hit := range res.Hits.Hits {
		ids = append(ids, hit.Id)
	}

	return ids, nil
}

// Create creates a file.
// If successful returns the file id and a nil error,
// otherwise returns empty string and non-nil error if any occured.
func (s Store) Create(ctx context.Context, file *pb.File) (string, error) {
	res, err := s.client.Index().
		Index(s.index).
		Id(file.GetId()).
		BodyJson(file).
		Do(ctx)

	if err != nil {
		return "", err
	}

	return res.Id, nil
}
