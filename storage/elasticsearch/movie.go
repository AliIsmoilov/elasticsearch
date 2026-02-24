package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"elasticsearch/storage/repo"

	es "github.com/elastic/go-elasticsearch/v8"
)

type movieES struct {
	client *es.Client
}

func NewMovie(client *es.Client) repo.ElasticMovieI {
	return &movieES{client: client}
}

func (r *movieES) Search(ctx context.Context, req repo.GetAllMoviesReq) (repo.GetAllMoviesResp, error) {

	var query map[string]interface{}
	if req.Query == "" {
		query = map[string]interface{}{"match_all": map[string]interface{}{}}
	} else {
		query = map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"movie", "overview", "director", "actors", "screenwriter"},
			},
		}
	}

	body := map[string]interface{}{
		"query": query,
	}

	// pagination
	if req.Limit > 0 {
		body["size"] = req.Limit
		if req.Page > 0 {
			from := (req.Page - 1) * req.Limit
			body["from"] = from
		}
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return repo.GetAllMoviesResp{}, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("movies"),
		r.client.Search.WithBody(buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return repo.GetAllMoviesResp{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return repo.GetAllMoviesResp{}, fmt.Errorf("elasticsearch search error: %s", res.String())
	}

	var rbody struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source map[string]interface{} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&rbody); err != nil {
		return repo.GetAllMoviesResp{}, err
	}

	movies := make([]repo.Movie, 0, len(rbody.Hits.Hits))
	for _, h := range rbody.Hits.Hits {
		src := h.Source
		m := repo.Movie{}
		if v, ok := src["id"]; ok {
			switch val := v.(type) {
			case float64:
				m.Id = int64(val)
			case string:
				// try parse
				// ignore parse errors
			}
		}
		if v, ok := src["rating"]; ok {
			if f, ok := v.(float64); ok {
				m.Rating = int(f)
			}
		}
		if v, ok := src["movie"]; ok {
			if s, ok := v.(string); ok {
				m.Movie = s
			}
		}
		if v, ok := src["year"]; ok {
			if f, ok := v.(float64); ok {
				m.Year = int(f)
			}
		}
		if v, ok := src["country"]; ok {
			if s, ok := v.(string); ok {
				m.Country = s
			}
		}
		if v, ok := src["rating_ball"]; ok {
			if f, ok := v.(float64); ok {
				m.RatingBall = f
			}
		}
		if v, ok := src["overview"]; ok {
			if s, ok := v.(string); ok {
				m.Overview = s
			}
		}
		if v, ok := src["director"]; ok {
			if s, ok := v.(string); ok {
				m.Director = s
			}
		}
		if v, ok := src["screenwriter"]; ok {
			switch arr := v.(type) {
			case []interface{}:
				sa := make([]string, 0, len(arr))
				for _, it := range arr {
					if s, ok := it.(string); ok {
						sa = append(sa, s)
					}
				}
				m.Screenwriter = sa
			case string:
				m.Screenwriter = []string{arr}
			}
		}
		if v, ok := src["actors"]; ok {
			switch arr := v.(type) {
			case []interface{}:
				sa := make([]string, 0, len(arr))
				for _, it := range arr {
					if s, ok := it.(string); ok {
						sa = append(sa, s)
					}
				}
				m.Actors = sa
			case string:
				m.Actors = []string{arr}
			}
		}
		if v, ok := src["url_logo"]; ok {
			if s, ok := v.(string); ok {
				m.UrlLogo = s
			}
		}

		movies = append(movies, m)
	}

	return repo.GetAllMoviesResp{Movies: movies, Count: rbody.Hits.Total.Value}, nil
}
