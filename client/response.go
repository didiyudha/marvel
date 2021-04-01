package client

type Item struct {
	Resourceuri string `json:"resourceURI"`
	Name        string `json:"name"`
}

type Comics struct {
	Available     int    `json:"available"`
	Returned      int    `json:"returned"`
	Collectionuri string `json:"collectionURI"`
	Items         []Item `json:"items"`
}

type Thumbnail struct {
	Path      string `json:"path"`
	Extension string `json:"extension"`
}

type Url struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type StoriesItem struct {
	Resourceuri string `json:"resourceURI"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

type Stories struct {
	Available     int           `json:"available"`
	Returned      int           `json:"returned"`
	Collectionuri string        `json:"collectionURI"`
	Items         []StoriesItem `json:"items"`
}

type EventsItem struct {
	Resourceuri string `json:"resourceURI"`
	Name        string `json:"name"`
}

type Events struct {
	Available     int          `json:"available"`
	Returned      int          `json:"returned"`
	Collectionuri string       `json:"collectionURI"`
	Items         []EventsItem `json:"items"`
}

type SeriesItem struct {
	Resourceuri string `json:"resourceURI"`
	Name        string `json:"name"`
}

type Series struct {
	Available     int          `json:"available"`
	Returned      int          `json:"returned"`
	Collectionuri string       `json:"collectionURI"`
	Items         []SeriesItem `json:"items"`
}

type Result struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Modified    string    `json:"modified"`
	Resourceuri string    `json:"resourceURI"`
	Urls        []Url     `json:"urls"`
	Thumbnail   Thumbnail `json:"thumbnail"`
	Comics      Comics    `json:"comics"`
	Stories     Stories   `json:"stories"`
	Events      Events    `json:"events"`
	Series      Series    `json:"series"`
}

type Data struct {
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
	Total   int      `json:"total"`
	Count   int      `json:"count"`
	Results []Result `json:"results"`
}

type CharacterResponse struct {
	Code            int    `json:"code"`
	Status          string `json:"status"`
	Copyright       string `json:"copyright"`
	Attributiontext string `json:"attributionText"`
	Attributionhtml string `json:"attributionHTML"`
	Data            Data   `json:"data"`
	Etag            string `json:"etag"`
}

type ErrResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
