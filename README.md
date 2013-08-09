tzone
=================
Timezone conversion API with account tokens and API usage, etc.

Usage
-------------
curl http://localhost:9000/timeZones
curl http://localhost:9000/convertCurrent?to=America/Los_Angeles&token=xxx
curl -H "Accept: application/json" http://localhost:9000/current/America/Los_Angeles\?token=xxx
