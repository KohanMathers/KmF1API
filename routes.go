package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	kmf1 "github.com/KohanMathers/KmF1/go"
)

var sessionSlugs = map[string]string{
	"r":                 kmf1.SessionRace,
	"race":              kmf1.SessionRace,
	"q":                 kmf1.SessionQualifying,
	"qualifying":        kmf1.SessionQualifying,
	"s":                 kmf1.SessionSprint,
	"sprint":            kmf1.SessionSprint,
	"sq":                kmf1.SessionSprintQualifying,
	"sprint-qualifying": kmf1.SessionSprintQualifying,
	"fp1":               kmf1.SessionPractice1,
	"practice1":         kmf1.SessionPractice1,
	"practice-1":        kmf1.SessionPractice1,
	"fp2":               kmf1.SessionPractice2,
	"practice2":         kmf1.SessionPractice2,
	"practice-2":        kmf1.SessionPractice2,
	"fp3":               kmf1.SessionPractice3,
	"practice3":         kmf1.SessionPractice3,
	"practice-3":        kmf1.SessionPractice3,
}

var eventAliases = map[string]string{
	"bahrain":        "Bahrain Grand Prix",
	"saudi_arabia":   "Saudi Arabian Grand Prix",
	"australia":      "Australian Grand Prix",
	"japan":          "Japanese Grand Prix",
	"china":          "Chinese Grand Prix",
	"miami":          "Miami Grand Prix",
	"emilia_romagna": "Emilia Romagna Grand Prix",
	"monaco":         "Monaco Grand Prix",
	"canada":         "Canadian Grand Prix",
	"spain":          "Spanish Grand Prix",
	"austria":        "Austrian Grand Prix",
	"silverstone":    "British Grand Prix",
	"hungary":        "Hungarian Grand Prix",
	"belgium":        "Belgian Grand Prix",
	"netherlands":    "Dutch Grand Prix",
	"italy":          "Italian Grand Prix",
	"azerbaijan":     "Azerbaijan Grand Prix",
	"singapore":      "Singapore Grand Prix",
	"texas":          "United States Grand Prix",
	"mexico":         "Mexico City Grand Prix",
	"brazil":         "São Paulo Grand Prix",
	"las_vegas":      "Las Vegas Grand Prix",
	"qatar":          "Qatar Grand Prix",
	"abu_dhabi":      "Abu Dhabi Grand Prix",
}

type apiHandler struct {
	client   *kmf1.Client
	logger   *slog.Logger
	sessions *sessionCache
}

type routeSpec struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"handler"`
}

var apiRoutes = []routeSpec{
	{Method: "GET", Path: "/health", Handler: "health"},
	{Method: "GET", Path: "/routes", Handler: "listRoutes"},
	{Method: "GET", Path: "/api/v1/drivers", Handler: "listDriverReferences"},
	{Method: "GET", Path: "/api/v1/drivers/{driverId}", Handler: "getDriverReference"},
	{Method: "GET", Path: "/api/v1/constructors", Handler: "listConstructors"},
	{Method: "GET", Path: "/api/v1/constructors/{constructorId}", Handler: "getConstructor"},
	{Method: "GET", Path: "/api/v1/circuits", Handler: "listCircuitReferences"},
	{Method: "GET", Path: "/api/v1/circuits/{circuitId}", Handler: "getCircuitReference"},
	{Method: "GET", Path: "/api/v1/finishing-statuses", Handler: "listFinishingStatuses"},
	{Method: "GET", Path: "/api/v1/season/{year}/events", Handler: "listEvents"},
	{Method: "GET", Path: "/api/v1/season/{year}/teams", Handler: "listTeamStyles"},
	{Method: "GET", Path: "/api/v1/season/{year}/teams/{team}", Handler: "teamStyle"},
	{Method: "GET", Path: "/api/v1/season/{year}/compounds", Handler: "listCompounds"},
	{Method: "GET", Path: "/api/v1/season/{year}/compounds/{compound}", Handler: "compoundColor"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/info", Handler: "eventInfo"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/sessions", Handler: "eventSessions"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/circuit", Handler: "eventCircuit"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/circuit/corners", Handler: "circuitCorners"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/circuit/marshal-lights", Handler: "circuitMarshalLights"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/circuit/marshal-sectors", Handler: "circuitMarshalSectors"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/standings/drivers", Handler: "standingsDrivers"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/standings/constructors", Handler: "standingsConstructors"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/info", Handler: "sessionInfo"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers", Handler: "sessionDrivers"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/laps", Handler: "sessionLaps"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/results", Handler: "sessionResults"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/best-stats", Handler: "sessionBestStats"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/weather", Handler: "sessionWeather"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/track-status", Handler: "sessionTrackStatus"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/progress", Handler: "sessionProgress"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/status-events", Handler: "sessionStatusEvents"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/clock", Handler: "sessionClock"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/archive-status", Handler: "sessionArchiveStatus"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/race-control-messages", Handler: "sessionRaceControlMessages"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/team-radio", Handler: "sessionTeamRadio"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/pit-stops", Handler: "sessionPitStops"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/top-three", Handler: "sessionTopThree"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/championship-prediction", Handler: "sessionChampionshipPrediction"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/broadcast", Handler: "sessionBroadcast"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/broadcast/audio", Handler: "sessionBroadcastAudio"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/broadcast/content", Handler: "sessionBroadcastContent"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}", Handler: "driverSummary"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/info", Handler: "driverInfo"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/result", Handler: "driverResult"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/reference", Handler: "driverReference"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/best-stats", Handler: "driverBestStats"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/laps", Handler: "driverLaps"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/laps/{n}", Handler: "driverLap"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/laps/{n}/telemetry", Handler: "driverLapTelemetry"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/telemetry", Handler: "driverTelemetry"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/intervals", Handler: "driverIntervals"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/stints", Handler: "driverStints"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/pit-stops", Handler: "driverPitStops"},
	{Method: "GET", Path: "/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/team-radio", Handler: "driverTeamRadio"},
}

func newMux(client *kmf1.Client, logger *slog.Logger) *http.ServeMux {
	h := &apiHandler{client: client, logger: logger, sessions: newSessionCache()}
	mux := http.NewServeMux()

	for _, route := range apiRoutes {
		switch route.Handler {
		case "health":
			mux.HandleFunc(route.Method+" "+route.Path, h.health)
		case "listRoutes":
			mux.HandleFunc(route.Method+" "+route.Path, h.listRoutes)
		case "listDriverReferences":
			mux.HandleFunc(route.Method+" "+route.Path, h.listDriverReferences)
		case "getDriverReference":
			mux.HandleFunc(route.Method+" "+route.Path, h.getDriverReference)
		case "listConstructors":
			mux.HandleFunc(route.Method+" "+route.Path, h.listConstructors)
		case "getConstructor":
			mux.HandleFunc(route.Method+" "+route.Path, h.getConstructor)
		case "listCircuitReferences":
			mux.HandleFunc(route.Method+" "+route.Path, h.listCircuitReferences)
		case "getCircuitReference":
			mux.HandleFunc(route.Method+" "+route.Path, h.getCircuitReference)
		case "listFinishingStatuses":
			mux.HandleFunc(route.Method+" "+route.Path, h.listFinishingStatuses)
		case "listEvents":
			mux.HandleFunc(route.Method+" "+route.Path, h.listEvents)
		case "listTeamStyles":
			mux.HandleFunc(route.Method+" "+route.Path, h.listTeamStyles)
		case "teamStyle":
			mux.HandleFunc(route.Method+" "+route.Path, h.teamStyle)
		case "listCompounds":
			mux.HandleFunc(route.Method+" "+route.Path, h.listCompounds)
		case "compoundColor":
			mux.HandleFunc(route.Method+" "+route.Path, h.compoundColor)
		case "eventInfo":
			mux.HandleFunc(route.Method+" "+route.Path, h.eventInfo)
		case "eventSessions":
			mux.HandleFunc(route.Method+" "+route.Path, h.eventSessions)
		case "eventCircuit":
			mux.HandleFunc(route.Method+" "+route.Path, h.eventCircuit)
		case "circuitCorners":
			mux.HandleFunc(route.Method+" "+route.Path, h.circuitCorners)
		case "circuitMarshalLights":
			mux.HandleFunc(route.Method+" "+route.Path, h.circuitMarshalLights)
		case "circuitMarshalSectors":
			mux.HandleFunc(route.Method+" "+route.Path, h.circuitMarshalSectors)
		case "standingsDrivers":
			mux.HandleFunc(route.Method+" "+route.Path, h.standingsDrivers)
		case "standingsConstructors":
			mux.HandleFunc(route.Method+" "+route.Path, h.standingsConstructors)
		case "sessionInfo":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionInfo)
		case "sessionDrivers":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionDrivers)
		case "sessionLaps":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionLaps)
		case "sessionResults":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionResults)
		case "sessionBestStats":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionBestStats)
		case "sessionWeather":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionWeather)
		case "sessionTrackStatus":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionTrackStatus)
		case "sessionProgress":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionProgress)
		case "sessionStatusEvents":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionStatusEvents)
		case "sessionClock":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionClock)
		case "sessionArchiveStatus":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionArchiveStatus)
		case "sessionRaceControlMessages":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionRaceControlMessages)
		case "sessionTeamRadio":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionTeamRadio)
		case "sessionPitStops":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionPitStops)
		case "sessionTopThree":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionTopThree)
		case "sessionChampionshipPrediction":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionChampionshipPrediction)
		case "sessionBroadcast":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionBroadcast)
		case "sessionBroadcastAudio":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionBroadcastAudio)
		case "sessionBroadcastContent":
			mux.HandleFunc(route.Method+" "+route.Path, h.sessionBroadcastContent)
		case "driverSummary":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverSummary)
		case "driverInfo":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverInfo)
		case "driverResult":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverResult)
		case "driverReference":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverReference)
		case "driverBestStats":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverBestStats)
		case "driverLaps":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverLaps)
		case "driverLap":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverLap)
		case "driverLapTelemetry":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverLapTelemetry)
		case "driverTelemetry":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverTelemetry)
		case "driverIntervals":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverIntervals)
		case "driverStints":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverStints)
		case "driverPitStops":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverPitStops)
		case "driverTeamRadio":
			mux.HandleFunc(route.Method+" "+route.Path, h.driverTeamRadio)
		}
	}

	return mux
}

const (
	cacheImmutable = "public, max-age=31536000, immutable"
	cacheLive      = "public, max-age=5"
)

func seasonCacheControl(year int) string {
	if year < time.Now().Year() {
		return cacheImmutable
	}
	return cacheLive
}

func sessionCacheControl(session *kmf1.Session) string {
	if GetSessionArchiveStatus(session).Status == "Complete" {
		return cacheImmutable
	}
	return cacheLive
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONCached(w http.ResponseWriter, status int, v any, cacheControl string) {
	w.Header().Set("Cache-Control", cacheControl)
	writeJSON(w, status, v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *apiHandler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *apiHandler) listRoutes(w http.ResponseWriter, r *http.Request) {
	writeJSONCached(w, http.StatusOK, map[string]any{
		"routes": apiRoutes,
		"count":  len(apiRoutes),
	}, cacheImmutable)
}

func pathYear(w http.ResponseWriter, r *http.Request) (int, bool) {
	year, err := strconv.Atoi(r.PathValue("year"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid year: "+r.PathValue("year"))
		return 0, false
	}
	return year, true
}

type testingEvent int

func pathEventIdentifier(r *http.Request) any {
	raw := r.PathValue("event")
	if n, err := strconv.Atoi(raw); err == nil {
		return n
	}
	if n, ok := parseTestingIdentifier(raw); ok {
		return testingEvent(n)
	}
	if name, ok := eventAliases[strings.ToLower(raw)]; ok {
		return name
	}
	return raw
}

func parseTestingIdentifier(raw string) (int, bool) {
	s := strings.ToLower(raw)
	for _, prefix := range []string{"testing", "test"} {
		if !strings.HasPrefix(s, prefix) {
			continue
		}
		rest := strings.TrimLeft(strings.TrimPrefix(s, prefix), "-_")
		if rest == "" {
			return 1, true
		}
		if n, err := strconv.Atoi(rest); err == nil && n > 0 {
			return n, true
		}
		return 0, false
	}
	return 0, false
}

func pathSessionType(w http.ResponseWriter, r *http.Request, event *kmf1.Event) (string, bool) {
	raw := r.PathValue("session")
	if sessionType, ok := sessionSlugs[strings.ToLower(raw)]; ok {
		return sessionType, true
	}

	want := normalizeSessionName(raw)
	sessions := GetEventSessions(event)
	names := make([]string, len(sessions))
	for i, s := range sessions {
		if normalizeSessionName(s.SessionType) == want {
			return s.SessionType, true
		}
		names[i] = s.SessionType
	}
	writeError(w, http.StatusNotFound, "unknown session "+strconv.Quote(raw)+"; this event has: "+strings.Join(names, ", "))
	return "", false
}

func normalizeSessionName(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case ' ', '-', '_':
			return -1
		}
		return unicode.ToLower(r)
	}, s)
}

type link struct {
	Href string `json:"href"`
}

func pathLink(parts ...string) string {
	escaped := make([]string, len(parts))
	for i, part := range parts {
		escaped[i] = url.PathEscape(part)
	}
	return "/" + strings.Join(escaped, "/")
}

func (h *apiHandler) loadEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) (*kmf1.Event, bool) {
	year, ok := pathYear(w, r)
	if !ok {
		return nil, false
	}
	season := GetSeason(h.client, year)
	if _, err := GetEvents(ctx, season); err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return nil, false
	}
	var event *kmf1.Event
	var err error
	switch id := pathEventIdentifier(r).(type) {
	case testingEvent:
		event, err = GetTestingEvent(ctx, season, int(id))
	case int:
		if id == 0 {
			err = fmt.Errorf("round 0 is used by testing events; address them as testing-1, testing-2, ...")
			break
		}
		event, err = GetEvent(ctx, season, id)
	default:
		event, err = GetEvent(ctx, season, id)
	}
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return nil, false
	}
	w.Header().Set("Cache-Control", seasonCacheControl(year))
	return event, true
}

func (h *apiHandler) loadSession(ctx context.Context, w http.ResponseWriter, r *http.Request) (*kmf1.Session, bool) {
	event, ok := h.loadEvent(ctx, w, r)
	if !ok {
		return nil, false
	}
	sessionType, ok := pathSessionType(w, r, event)
	if !ok {
		return nil, false
	}
	ref, err := GetEventSession(event, sessionType)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return nil, false
	}

	info := GetEventInfo(event)
	key := sessionCacheKey{year: info.Season, round: info.RoundNumber, event: info.EventName, sessionType: sessionType}
	if len(info.Sessions) > 0 {
		// Testing events share round 0 and can share a name, so the start time tells them apart.
		key.eventStart = info.Sessions[0].StartDateUtc
	}
	session, err := h.sessions.getOrLoad(ctx, key, func() (*kmf1.Session, error) {
		return LoadSession(ctx, ref)
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return nil, false
	}
	w.Header().Set("Cache-Control", sessionCacheControl(session))
	return session, true
}

func (h *apiHandler) loadDriver(ctx context.Context, w http.ResponseWriter, r *http.Request) (*kmf1.Driver, bool) {
	session, ok := h.loadSession(ctx, w, r)
	if !ok {
		return nil, false
	}
	driver := GetSessionDriver(session, r.PathValue("driver"))
	if !driver.Exists() {
		writeError(w, http.StatusNotFound, "driver not found in session: "+r.PathValue("driver"))
		return nil, false
	}
	return driver, true
}

func (h *apiHandler) listDriverReferences(w http.ResponseWriter, r *http.Request) {
	refs, err := GetDriverReferences(r.Context(), h.client)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSONCached(w, http.StatusOK, refs, cacheImmutable)
}

func (h *apiHandler) getDriverReference(w http.ResponseWriter, r *http.Request) {
	ref, err := GetDriverReference(r.Context(), h.client, r.PathValue("driverId"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSONCached(w, http.StatusOK, ref, cacheImmutable)
}

func (h *apiHandler) listConstructors(w http.ResponseWriter, r *http.Request) {
	refs, err := GetConstructors(r.Context(), h.client)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSONCached(w, http.StatusOK, refs, cacheImmutable)
}

func (h *apiHandler) getConstructor(w http.ResponseWriter, r *http.Request) {
	ref, err := GetConstructor(r.Context(), h.client, r.PathValue("constructorId"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSONCached(w, http.StatusOK, ref, cacheImmutable)
}

func (h *apiHandler) listCircuitReferences(w http.ResponseWriter, r *http.Request) {
	refs, err := GetCircuitReferences(r.Context(), h.client)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSONCached(w, http.StatusOK, refs, cacheImmutable)
}

func (h *apiHandler) getCircuitReference(w http.ResponseWriter, r *http.Request) {
	ref, err := GetCircuitReference(r.Context(), h.client, r.PathValue("circuitId"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSONCached(w, http.StatusOK, ref, cacheImmutable)
}

func (h *apiHandler) listFinishingStatuses(w http.ResponseWriter, r *http.Request) {
	statuses, err := GetFinishingStatuses(r.Context(), h.client)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSONCached(w, http.StatusOK, statuses, cacheImmutable)
}

func (h *apiHandler) listEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	year, ok := pathYear(w, r)
	if !ok {
		return
	}
	events, err := GetEvents(ctx, GetSeason(h.client, year))
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	infos := make([]any, len(events))
	for i := range events {
		infos[i] = GetEventInfo(&events[i])
	}
	writeJSONCached(w, http.StatusOK, infos, seasonCacheControl(year))
}

func (h *apiHandler) listTeamStyles(w http.ResponseWriter, r *http.Request) {
	year, ok := pathYear(w, r)
	if !ok {
		return
	}
	writeJSONCached(w, http.StatusOK, GetTeamStyles(GetSeason(h.client, year)), seasonCacheControl(year))
}

func (h *apiHandler) teamStyle(w http.ResponseWriter, r *http.Request) {
	year, ok := pathYear(w, r)
	if !ok {
		return
	}
	style, ok := GetTeamStyle(GetSeason(h.client, year), r.PathValue("team"))
	if !ok {
		writeError(w, http.StatusNotFound, "no styling known for team: "+r.PathValue("team"))
		return
	}
	writeJSONCached(w, http.StatusOK, style, seasonCacheControl(year))
}

func (h *apiHandler) listCompounds(w http.ResponseWriter, r *http.Request) {
	year, ok := pathYear(w, r)
	if !ok {
		return
	}
	writeJSONCached(w, http.StatusOK, GetCompounds(GetSeason(h.client, year)), seasonCacheControl(year))
}

func (h *apiHandler) compoundColor(w http.ResponseWriter, r *http.Request) {
	year, ok := pathYear(w, r)
	if !ok {
		return
	}
	compound := r.PathValue("compound")
	color, ok := GetCompoundColor(GetSeason(h.client, year), compound)
	if !ok {
		writeError(w, http.StatusNotFound, "no color known for compound: "+compound)
		return
	}
	writeJSONCached(w, http.StatusOK, map[string]string{
		"compound": compound,
		"color":    color,
	}, seasonCacheControl(year))
}

func (h *apiHandler) eventInfo(w http.ResponseWriter, r *http.Request) {
	event, ok := h.loadEvent(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetEventInfo(event))
}

func (h *apiHandler) eventSessions(w http.ResponseWriter, r *http.Request) {
	event, ok := h.loadEvent(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetEventSessions(event))
}

func (h *apiHandler) eventCircuit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, ok := h.loadEvent(ctx, w, r)
	if !ok {
		return
	}
	circuit, err := GetEventCircuit(event)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	corners, err := GetCircuitCorners(ctx, circuit)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	lights, err := GetCircuitMarshalLights(ctx, circuit)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	sectors, err := GetCircuitMarshalSectors(ctx, circuit)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	base := pathLink(
		"api", "v1", "season",
		r.PathValue("year"),
		"events",
		r.PathValue("event"),
		"circuit",
	)
	writeJSON(w, http.StatusOK, map[string]any{
		"_links": map[string]any{
			"self":           link{Href: base},
			"corners":        link{Href: base + "/corners"},
			"marshalLights":  link{Href: base + "/marshal-lights"},
			"marshalSectors": link{Href: base + "/marshal-sectors"},
		},
		"corners":        corners,
		"marshalLights":  lights,
		"marshalSectors": sectors,
	})
}

func (h *apiHandler) circuitCorners(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, ok := h.loadEvent(ctx, w, r)
	if !ok {
		return
	}
	circuit, err := GetEventCircuit(event)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	corners, err := GetCircuitCorners(ctx, circuit)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, corners)
}

func (h *apiHandler) circuitMarshalLights(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, ok := h.loadEvent(ctx, w, r)
	if !ok {
		return
	}
	circuit, err := GetEventCircuit(event)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	lights, err := GetCircuitMarshalLights(ctx, circuit)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, lights)
}

func (h *apiHandler) circuitMarshalSectors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, ok := h.loadEvent(ctx, w, r)
	if !ok {
		return
	}
	circuit, err := GetEventCircuit(event)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	sectors, err := GetCircuitMarshalSectors(ctx, circuit)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sectors)
}

func (h *apiHandler) standingsDrivers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, ok := h.loadEvent(ctx, w, r)
	if !ok {
		return
	}
	drivers, err := GetStandingsDrivers(ctx, GetEventStandings(event))
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, drivers)
}

func (h *apiHandler) standingsConstructors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event, ok := h.loadEvent(ctx, w, r)
	if !ok {
		return
	}
	constructors, err := GetStandingsConstructors(ctx, GetEventStandings(event))
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, constructors)
}

func (h *apiHandler) sessionInfo(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionInfo(session))
}

func (h *apiHandler) sessionDrivers(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionDrivers(session))
}

func (h *apiHandler) sessionLaps(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionLaps(session))
}

func (h *apiHandler) sessionResults(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionResults(session))
}

func (h *apiHandler) sessionBestStats(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionBestStats(session))
}

func (h *apiHandler) sessionWeather(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionWeather(session))
}

func (h *apiHandler) sessionTrackStatus(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionTrackStatus(session))
}

func (h *apiHandler) sessionProgress(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionProgress(session))
}

func (h *apiHandler) sessionStatusEvents(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionStatusEvents(session))
}

func (h *apiHandler) sessionClock(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionClock(session))
}

func (h *apiHandler) sessionArchiveStatus(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionArchiveStatus(session))
}

func (h *apiHandler) sessionRaceControlMessages(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionRaceControlMessages(session))
}

func (h *apiHandler) sessionTeamRadio(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionTeamRadio(session))
}

func (h *apiHandler) sessionPitStops(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionPitStops(session))
}

func (h *apiHandler) sessionTopThree(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetSessionTopThree(session))
}

func (h *apiHandler) sessionChampionshipPrediction(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	cp := GetSessionChampionshipPrediction(session)
	if cp == nil {
		writeError(w, http.StatusNotFound, "championship prediction is only available for race sessions")
		return
	}
	base := pathLink(
		"api", "v1", "season",
		r.PathValue("year"),
		"events",
		r.PathValue("event"),
		r.PathValue("session"),
		"championship-prediction",
	)
	writeJSON(w, http.StatusOK, map[string]any{
		"_links": map[string]any{
			"self":    link{Href: base},
			"drivers": link{Href: base + "/drivers"},
			"teams":   link{Href: base + "/teams"},
		},
		"drivers": GetChampionshipPredictionDrivers(cp),
		"teams":   GetChampionshipPredictionTeams(cp),
	})
}

func (h *apiHandler) sessionBroadcast(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	b := GetSessionBroadcast(session)
	base := pathLink(
		"api", "v1", "season",
		r.PathValue("year"),
		"events",
		r.PathValue("event"),
		r.PathValue("session"),
		"broadcast",
	)
	writeJSON(w, http.StatusOK, map[string]any{
		"_links": map[string]any{
			"self":    link{Href: base},
			"audio":   link{Href: base + "/audio"},
			"content": link{Href: base + "/content"},
		},
		"audio":   GetBroadcastAudio(b),
		"content": GetBroadcastContent(b),
	})
}

func (h *apiHandler) sessionBroadcastAudio(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetBroadcastAudio(GetSessionBroadcast(session)))
}

func (h *apiHandler) sessionBroadcastContent(w http.ResponseWriter, r *http.Request) {
	session, ok := h.loadSession(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetBroadcastContent(GetSessionBroadcast(session)))
}

func (h *apiHandler) driverSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	driver, ok := h.loadDriver(ctx, w, r)
	if !ok {
		return
	}
	base := pathLink(
		"api", "v1", "season",
		r.PathValue("year"),
		"events",
		r.PathValue("event"),
		r.PathValue("session"),
		"drivers",
		r.PathValue("driver"),
	)
	resp := map[string]any{
		"_links": map[string]any{
			"self":      link{Href: base},
			"info":      link{Href: base + "/info"},
			"result":    link{Href: base + "/result"},
			"reference": link{Href: base + "/reference"},
			"bestStats": link{Href: base + "/best-stats"},
			"laps":      link{Href: base + "/laps"},
			"telemetry": link{Href: base + "/telemetry"},
			"intervals": link{Href: base + "/intervals"},
			"stints":    link{Href: base + "/stints"},
			"pitStops":  link{Href: base + "/pit-stops"},
			"teamRadio": link{Href: base + "/team-radio"},
		},
		"info": GetDriverInfo(driver),
	}
	if result, ok := GetDriverResult(driver); ok {
		resp["result"] = result
	}
	if stats, ok := GetDriverBestStats(driver); ok {
		resp["bestStats"] = stats
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *apiHandler) driverInfo(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.loadDriver(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetDriverInfo(driver))
}

func (h *apiHandler) driverResult(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.loadDriver(r.Context(), w, r)
	if !ok {
		return
	}
	result, ok := GetDriverResult(driver)
	if !ok {
		writeError(w, http.StatusNotFound, "no result recorded for driver")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *apiHandler) driverReference(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	driver, ok := h.loadDriver(ctx, w, r)
	if !ok {
		return
	}
	ref, err := GetDriverJolpicaReference(ctx, driver)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ref)
}

func (h *apiHandler) driverBestStats(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.loadDriver(r.Context(), w, r)
	if !ok {
		return
	}
	stats, ok := GetDriverBestStats(driver)
	if !ok {
		writeError(w, http.StatusNotFound, "no best-stats recorded for driver")
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *apiHandler) driverLaps(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.loadDriver(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetDriverLaps(driver))
}

func (h *apiHandler) driverLap(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.loadDriver(r.Context(), w, r)
	if !ok {
		return
	}
	n, err := strconv.Atoi(r.PathValue("n"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid lap number: "+r.PathValue("n"))
		return
	}
	lap, ok := GetDriverLap(driver, n)
	if !ok {
		writeError(w, http.StatusNotFound, "lap not found")
		return
	}
	writeJSON(w, http.StatusOK, lap)
}

func (h *apiHandler) driverLapTelemetry(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.loadDriver(r.Context(), w, r)
	if !ok {
		return
	}
	n, err := strconv.Atoi(r.PathValue("n"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid lap number: "+r.PathValue("n"))
		return
	}
	lap, ok := GetDriverLap(driver, n)
	if !ok {
		writeError(w, http.StatusNotFound, "lap not found")
		return
	}
	writeJSON(w, http.StatusOK, GetLapTelemetry(lap))
}

func (h *apiHandler) driverTelemetry(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.loadDriver(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetDriverTelemetry(driver))
}

func (h *apiHandler) driverIntervals(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.loadDriver(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetDriverIntervals(driver))
}

func (h *apiHandler) driverStints(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.loadDriver(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetDriverStints(driver))
}

func (h *apiHandler) driverPitStops(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.loadDriver(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetDriverPitStops(driver))
}

func (h *apiHandler) driverTeamRadio(w http.ResponseWriter, r *http.Request) {
	driver, ok := h.loadDriver(r.Context(), w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, GetDriverTeamRadio(driver))
}
