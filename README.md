# KmF1 API

An API wrapping [KmF1](https://github.com/KohanMathers/KmF1).

That's it.

That's the whole thing.

Deploy it yourself or use it at https://raceplayback.com

Routes list:

- `/health`
- `/routes`
- `/api/v1/drivers`
- `/api/v1/drivers/{driverId}`
- `/api/v1/constructors`
- `/api/v1/constructors/{constructorId}`
- `/api/v1/circuits`
- `/api/v1/circuits/{circuitId}`
- `/api/v1/finishing-statuses`
- `/api/v1/season/{year}/events`
- `/api/v1/season/{year}/teams`
- `/api/v1/season/{year}/teams/{team}`
- `/api/v1/season/{year}/compounds`
- `/api/v1/season/{year}/compounds/{compound}`
- `/api/v1/season/{year}/events/{event}/info`
- `/api/v1/season/{year}/events/{event}/sessions`
- `/api/v1/season/{year}/events/{event}/circuit`
- `/api/v1/season/{year}/events/{event}/circuit/corners`
- `/api/v1/season/{year}/events/{event}/circuit/marshal-lights`
- `/api/v1/season/{year}/events/{event}/circuit/marshal-sectors`
- `/api/v1/season/{year}/events/{event}/standings/drivers`
- `/api/v1/season/{year}/events/{event}/standings/constructors`
- `/api/v1/season/{year}/events/{event}/{session}/info`
- `/api/v1/season/{year}/events/{event}/{session}/drivers`
- `/api/v1/season/{year}/events/{event}/{session}/laps`
- `/api/v1/season/{year}/events/{event}/{session}/results`
- `/api/v1/season/{year}/events/{event}/{session}/best-stats`
- `/api/v1/season/{year}/events/{event}/{session}/weather`
- `/api/v1/season/{year}/events/{event}/{session}/track-status`
- `/api/v1/season/{year}/events/{event}/{session}/progress`
- `/api/v1/season/{year}/events/{event}/{session}/status-events`
- `/api/v1/season/{year}/events/{event}/{session}/clock`
- `/api/v1/season/{year}/events/{event}/{session}/archive-status`
- `/api/v1/season/{year}/events/{event}/{session}/race-control-messages`
- `/api/v1/season/{year}/events/{event}/{session}/team-radio`
- `/api/v1/season/{year}/events/{event}/{session}/pit-stops`
- `/api/v1/season/{year}/events/{event}/{session}/top-three`
- `/api/v1/season/{year}/events/{event}/{session}/championship-prediction`
- `/api/v1/season/{year}/events/{event}/{session}/broadcast`
- `/api/v1/season/{year}/events/{event}/{session}/broadcast/audio`
- `/api/v1/season/{year}/events/{event}/{session}/broadcast/content`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/info`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/result`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/reference`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/best-stats`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/laps`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/laps/{n}`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/laps/{n}/telemetry`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/telemetry`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/intervals`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/stints`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/pit-stops`
- `/api/v1/season/{year}/events/{event}/{session}/drivers/{driver}/team-radio`
