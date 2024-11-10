# Ad Impression Counter - Code Challenge

## How to start
```bash
git clone https://github.com/entertainment-veks/ad-impression-counter
cd ad-impression-counter
go run main.go
```

## Config
| Configuration         | Type          | Environment Variable            | Default Value | Description                                                                 |
|-----------------------|---------------|----------------------------------|---------------|-----------------------------------------------------------------------------|
| ServerPort            | string        | `env:"SERVER_PORT"`              | "8080"        | ServerPort is the port the HTTP server will listen on.                      |
| TTL                   | time.Duration | `env:"TTL"`                      | "3600s"       | TTL is the duplicate impressions detection window duration in seconds.      |
| ImpressionsQueueSize  | int           | `env:"IMPRESSIONS_QUEUE_SIZE"`   | 1000          | ImpressionsQueueSize is the size of the impressions queue.                  |
| ImpressionWorkers     | int           | `env:"IMPRESSION_WORKERS"`       | 10            | ImpressionWorkers is the number of workers that will process impressions. |

In production `ImpressionWorkers` sould be set to the number of CPU cores.

## How to interact
Entire documentation is described in [swagger.yaml](./swagger.yaml).

To create ad campaign we have to call POST `/api/v1/campaigns` endpoint, example: 
```bash
curl --location 'http://localhost:8080/api/v1/campaigns' \
--header 'Content-Type: application/json' \
--data '{
           "name": "Test Campaign",
           "start_time": "2024-11-09T16:29:00Z"
         }'
``` 
Please, pay attention to the `start_time`, it has to be in the future. 
Otherwise we will receive an error. 
This endpoint work synchronously. 

To track impression we have to call POST `/api/v1/impressions` endpoint, example:
```bash
curl --location 'http://localhost:8080/api/v1/impressions' \
--header 'Content-Type: application/json' \
--data '{
           "campaign_id": "b7da5a35-0974-4b7b-ab5d-576e9e628d80",
           "user_id": "user-1",
           "ad_id": "ad-1"
         }'
```
Please make sure to set the `campaign_id` to the one retrieved from the response of Create Campaign endpoint.
This endpoint work asynchronously. 

To check statistics we have to call GET `/api/v1/campaigns/{campaign_id}/stats` endpoint, example:
```bash
curl --location 'http://localhost:8080/api/v1/campaigns/b7da5a35-0974-4b7b-ab5d-576e9e628d80/stats'
```
Please make sure to set the `campaign_id` in URL to the one retrieved from the response of Create Campaign endpoint.

## How it works?
Create campaign is a synchronic endpoint, that create a campaign for us.
It has `start_time` field, that has to be in the future 
and any provided impession tracks will be discarded until `start_time` is in future.

Track impression is an asynchronic endpoint, that put all calls in a queue, 
that furter going to be consumed by workers. Queue size can be configured with `IMPRESSIONS_QUEUE_SIZE` env variable, by default - 1000. 

Workers are working in a backgound, we can configure their quantity with `IMPRESSION_WORKERS` env variable, by default - 10. 
In case of any errors (campaign doesn't exist/didn't start yet) worker just discard an impression.

Get stats is a synchronic endpoint, that return us statistics for the campaign.

All impressions are stored as a concurrent-safe map from sync.Map, where `key` is campaign id and `value` is a slice of impressions, sorted by time added. 
As result we can save with O(1) complexity in the end of the slice that can be retrived with O(1) by capmaign id from sync.Map.
And we can count stats with O(n) complexity when we iterate the slice with only relevant impressions, that can be retrived with O(1) from the map.
Also, since we have sorted by added time impressions, we are able to stop on 24 hour old impressions and count total count with len(slice).

Such storage solution can be optimised with using lock-free map like [cornelk hashmap](https://github.com/cornelk/hashmap). 
But in accordance with the task I decided to not use external library.

# Task:

## Ad Impression Counter - Test Assignment
Create a concurrent service that tracks ad impressions across multiple campaigns in real-time, with mechanisms to avoid counting duplicate impressions from the same user within a specified time period.

**Note:** Consider the task as a guideline. You have great freedom to experiment, showcase interesting snippets, ideas, and workarounds.

### Task
The goal is to build an in-memory impression counter that supports concurrent processing, including handling duplicate impressions efficiently and providing a REST API to manage campaigns, track impressions, and retrieve statistics.

### Instructions

1. **Implement Core Data Structures:**
   - **Campaign:** Represents an advertising campaign.
   - **Impression:** Represents each ad view, including `userID` and `timestamp`.
   - **Stats:** Stores aggregated impression statistics for each campaign.

2. **Duplicate Handling:**
   - Track impressions uniquely for each `userID` to avoid duplicate counting.
   - Implement a TTL (Time-To-Live) mechanism to discard duplicate records after a set time period (e.g., one hour).
   - Ensure that each unique impression from a user is only counted once within the specified TTL.

3. **Concurrency:**
   - Use concurrency controls (e.g., channels, mutexes) to handle multiple impressions tracked simultaneously without race conditions.
   - Make the impression tracking system thread-safe.

4. **Build REST API Endpoints:**
   - **POST /api/v1/campaigns** — Register a new campaign.
      - **Request Body:** `CreateCampaignRequest`
   - **POST /api/v1/impressions** — Track a new impression.
      - **Request Body:** `TrackImpressionRequest`
   - **GET /api/v1/campaigns/{id}/stats** — Get impression statistics for a specific campaign.
      - **Response Body:** `Stats`

---

### Data Structures

```go
// Campaign represents an advertising campaign
type Campaign struct {
    ID string `json:"id"`
    Name string `json:"name"`
    StartTime time.Time `json:"start_time"`
}

// Impression represents a single ad view
type Impression struct {
    CampaignID string `json:"campaign_id"`
    Timestamp time.Time `json:"timestamp"`
    UserID string `json:"user_id"`
    AdID string `json:"ad_id"`
}

// Stats represents aggregated impression statistics
type Stats struct {
    CampaignID string `json:"campaign_id"`
    LastHour int64 `json:"last_hour"`
    LastDay int64 `json:"last_day"`
    TotalCount int64 `json:"total"`
}
```

---

### API Requests/Responses

```go
// CreateCampaignRequest represents the request structure for creating a campaign
type CreateCampaignRequest struct {
    Name string `json:"name"`
    StartTime time.Time `json:"start_time"`
}

// TrackImpressionRequest represents the request structure for tracking an impression
type TrackImpressionRequest struct {
    CampaignID string `json:"campaign_id"`
    UserID string `json:"user_id"`
    AdID string `json:"ad_id"`
}

// ErrorResponse defines the structure for error responses
type ErrorResponse struct {
    Error string `json:"error"`
    Code int `json:"code"`
}
```

---

### Technical Focus

- **Concurrency:** Use channels and mutexes to handle concurrent updates safely.
- **Duplicate Check:** Implement a TTL for `userID` records to prevent duplicate counting within a specified time period (e.g., one hour).
- **API Design:** Structure endpoints for intuitive use and handle edge cases.
- **Project Structure:** Organize code with scalability in mind, allowing for future growth.
- **Minimum Dependencies:** Prefer standard or well-known packages.

---

### Example Workflow

1. **Register a Campaign:**
   - Send a `POST` request to `/api/v1/campaigns` with `CreateCampaignRequest` in the JSON body.
   - Returns the new campaign’s ID and a success message.

2. **Track Impressions:**
   - Send a `POST` request to `/api/v1/impressions` with `TrackImpressionRequest` in the body.
   - Check for duplicate impressions (based on `userID` and TTL).
   - Returns a success message.

3. **Get Campaign Statistics:**
   - Send a `GET` request to `/api/v1/campaigns/{id}/stats`.
   - Returns the aggregated `Stats` for the campaign, including `LastHour`, `LastDay`, and `TotalCount`.

---

### Evaluation Criteria

1. **Concurrency Handling:**
   - Are channels and mutexes used effectively to manage concurrent updates and prevent race conditions?
   - Is the TTL for duplicate check implemented correctly?

2. **API Design:**
   - Are the endpoints structured in a RESTful manner?
   - Is duplicate detection efficient and scalable?

3. **Code Organization and Readability:**
   - Is the project organized and structured to support future expansion?
   - Is the code clear and well-commented?

4. **Error Handling and Documentation:**
   - Are error cases handled appropriately, with clear messages?
   - Is the code and API well-documented?

5. **Testing:**
   - Include tests to validate core functionality (e.g., concurrent impressions and duplicate check).
