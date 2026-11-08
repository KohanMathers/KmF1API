package main

import (
	"context"

	kmf1 "github.com/KohanMathers/KmF1/go"
)

func GetDriverReferences[T any](ctx context.Context, client interface {
	DriverRefs(context.Context) ([]T, error)
}) ([]T, error) {
	return client.DriverRefs(ctx)
}

func GetDriverReference[T any](ctx context.Context, client interface {
	DriverRef(context.Context, string) (T, error)
}, driverId string) (T, error) {
	return client.DriverRef(ctx, driverId)
}

func GetConstructors[T any](ctx context.Context, client interface {
	Constructors(context.Context) ([]T, error)
}) ([]T, error) {
	return client.Constructors(ctx)
}

func GetConstructor[T any](ctx context.Context, client interface {
	Constructor(context.Context, string) (T, error)
}, constructorId string) (T, error) {
	return client.Constructor(ctx, constructorId)
}

func GetCircuitReferences[T any](ctx context.Context, client interface {
	CircuitRefs(context.Context) ([]T, error)
}) ([]T, error) {
	return client.CircuitRefs(ctx)
}

func GetCircuitReference[T any](ctx context.Context, client interface {
	CircuitRef(context.Context, string) (T, error)
}, circuitId string) (T, error) {
	return client.CircuitRef(ctx, circuitId)
}

func GetFinishingStatuses[T any](ctx context.Context, client interface {
	FinishingStatuses(context.Context) ([]T, error)
}) ([]T, error) {
	return client.FinishingStatuses(ctx)
}

func GetSeason(client *kmf1.Client, year int) *kmf1.Season {
	return client.Year(year)
}

func GetEvents(ctx context.Context, season *kmf1.Season) ([]kmf1.Event, error) {
	return season.Events(ctx)
}

func GetEvent(ctx context.Context, season *kmf1.Season, identifier any) (*kmf1.Event, error) {
	return season.Event(ctx, identifier)
}

func GetTestingEvent(ctx context.Context, season *kmf1.Season, n int) (*kmf1.Event, error) {
	return season.TestingEvent(ctx, n)
}

func GetTeamStyle[T any](season interface {
	TeamStyle(string) (T, bool)
}, identifier string) (T, bool) {
	return season.TeamStyle(identifier)
}

func GetTeamStyles[T any](season interface{ TeamStyles() []T }) []T {
	return season.TeamStyles()
}

func GetCompoundColor(season *kmf1.Season, compound string) (string, bool) {
	return season.CompoundColor(compound)
}

func GetCompounds(season *kmf1.Season) []string {
	return season.Compounds()
}

func GetEventInfo[T any](event interface{ Info() T }) T {
	return event.Info()
}

func GetEventSessions[T any](event interface{ Sessions() []T }) []T {
	return event.Sessions()
}

func GetEventSession(event *kmf1.Event, sessionType string) (*kmf1.SessionRef, error) {
	return event.Session(sessionType)
}

func GetEventCircuit(event *kmf1.Event) (*kmf1.Circuit, error) {
	return event.Circuit()
}

func GetEventStandings(event *kmf1.Event) *kmf1.Standings {
	return event.Standings()
}

func LoadSession(ctx context.Context, ref *kmf1.SessionRef) (*kmf1.Session, error) {
	return ref.Load(ctx)
}

func GetSession(ctx context.Context, client *kmf1.Client, year int, eventIdentifier any, sessionType string) (*kmf1.Session, error) {
	event, err := GetEvent(ctx, GetSeason(client, year), eventIdentifier)
	if err != nil {
		return nil, err
	}
	ref, err := GetEventSession(event, sessionType)
	if err != nil {
		return nil, err
	}
	return LoadSession(ctx, ref)
}

func GetSessionInfo[T any](session interface{ Info() *T }) *T {
	return session.Info()
}

func GetSessionDrivers[T any](session interface{ Drivers() []T }) []T {
	return session.Drivers()
}

func GetSessionDriver(session *kmf1.Session, identifier string) *kmf1.Driver {
	return session.Driver(identifier)
}

func GetSessionLaps[T any](session interface{ Laps() []T }) []T {
	return session.Laps()
}

func GetSessionResults[T any](session interface{ Results() []T }) []T {
	return session.Results()
}

func GetSessionBestStats[T any](session interface{ BestStats() []T }) []T {
	return session.BestStats()
}

func GetSessionWeather[T any](session interface{ Weather() []T }) []T {
	return session.Weather()
}

func GetSessionTrackStatus[T any](session interface{ TrackStatus() []T }) []T {
	return session.TrackStatus()
}

func GetSessionProgress[T any](session interface{ Progress() []T }) []T {
	return session.Progress()
}

func GetSessionStatusEvents[T any](session interface{ StatusEvents() []T }) []T {
	return session.StatusEvents()
}

func GetSessionClock[T any](session interface{ Clock() []T }) []T {
	return session.Clock()
}

func GetSessionArchiveStatus[T any](session interface{ ArchiveStatus() T }) T {
	return session.ArchiveStatus()
}

func GetSessionRaceControlMessages[T any](session interface{ RaceControlMessages() []T }) []T {
	return session.RaceControlMessages()
}

func GetSessionTeamRadio[T any](session interface{ TeamRadio() []T }) []T {
	return session.TeamRadio()
}

func GetSessionPitStops[T any](session interface{ PitStops() []T }) []T {
	return session.PitStops()
}

func GetSessionTopThree[T any](session interface{ TopThree() []T }) []T {
	return session.TopThree()
}

func GetSessionChampionshipPrediction(session *kmf1.Session) *kmf1.ChampionshipPrediction {
	return session.ChampionshipPrediction()
}

func GetSessionBroadcast(session *kmf1.Session) *kmf1.Broadcast {
	return session.Broadcast()
}

func GetDriverInfo[T any](driver interface{ Info() T }) T {
	return driver.Info()
}

func GetDriverResult[T any](driver interface{ Result() (T, bool) }) (T, bool) {
	return driver.Result()
}

func GetDriverJolpicaReference[T any](ctx context.Context, driver interface {
	Reference(context.Context) (T, error)
}) (T, error) {
	return driver.Reference(ctx)
}

func GetDriverBestStats[T any](driver interface{ BestStats() (T, bool) }) (T, bool) {
	return driver.BestStats()
}

func GetDriverLaps[T any](driver interface{ Laps() []T }) []T {
	return driver.Laps()
}

func GetDriverLap(driver *kmf1.Driver, n int) (*kmf1.Lap, bool) {
	return driver.Lap(n)
}

func GetDriverTelemetry[T any](driver interface{ Telemetry() []T }) []T {
	return driver.Telemetry()
}

func GetDriverIntervals[T any](driver interface{ Intervals() []T }) []T {
	return driver.Intervals()
}

func GetDriverStints[T any](driver interface{ Stints() []T }) []T {
	return driver.Stints()
}

func GetDriverPitStops[T any](driver interface{ PitStops() []T }) []T {
	return driver.PitStops()
}

func GetDriverTeamRadio[T any](driver interface{ TeamRadio() []T }) []T {
	return driver.TeamRadio()
}

func GetLapTelemetry[T any](lap interface{ Telemetry() []T }) []T {
	return lap.Telemetry()
}

func GetCircuitCorners[T any](ctx context.Context, circuit interface {
	Corners(context.Context) ([]T, error)
}) ([]T, error) {
	return circuit.Corners(ctx)
}

func GetCircuitMarshalLights[T any](ctx context.Context, circuit interface {
	MarshalLights(context.Context) ([]T, error)
}) ([]T, error) {
	return circuit.MarshalLights(ctx)
}

func GetCircuitMarshalSectors[T any](ctx context.Context, circuit interface {
	MarshalSectors(context.Context) ([]T, error)
}) ([]T, error) {
	return circuit.MarshalSectors(ctx)
}

func GetStandingsDrivers[T any](ctx context.Context, standings interface {
	Drivers(context.Context) ([]T, error)
}) ([]T, error) {
	return standings.Drivers(ctx)
}

func GetStandingsConstructors[T any](ctx context.Context, standings interface {
	Constructors(context.Context) ([]T, error)
}) ([]T, error) {
	return standings.Constructors(ctx)
}

func GetChampionshipPredictionDrivers[T any](cp interface{ Drivers() []T }) []T {
	return cp.Drivers()
}

func GetChampionshipPredictionTeams[T any](cp interface{ Teams() []T }) []T {
	return cp.Teams()
}

func GetBroadcastAudio[T any](b interface{ Audio() []T }) []T {
	return b.Audio()
}

func GetBroadcastContent[T any](b interface{ Content() []T }) []T {
	return b.Content()
}
