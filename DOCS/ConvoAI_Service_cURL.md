# ConvoAI Service cURL Examples

This document provides cURL examples for testing the ConvoAI Service endpoints.

## Invite Agent

```bash
curl -X POST http://localhost:8080/agent/invite \
  -H "Content-Type: application/json" \
  -d '{
    "channel_name": "test-channel",
    "requester_id": "0",
    "input_modalities": ["text"],
    "output_modalities": ["text", "audio"]
  }'
```

### Example Response

```json
{
  "agent_id": "1NTAGR0RA8YEUVTCMGLOGWHIPQYPXAZO",
  "create_ts": 1739905500,
  "status": "RUNNING"
}
```

## Remove Agent

```bash
curl -X POST http://localhost:8080/agent/remove \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "1NTAGR0RA8YEUVTCMGLOGWHIPQYPXAZO"
  }'
```

### Example Response

```json
{
  "success": true,
  "agent_id": "1NT29X0XUN1CFS1VJBS11RAFSJFYBMOW"
}
```
