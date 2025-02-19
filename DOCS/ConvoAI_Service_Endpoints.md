# ConvoAI Service API Endpoints

This document provides details about the ConvoAI Service API endpoints.

## Invite Agent

Invites an AI agent to join a channel.

### Endpoint

`POST /agent/invite`

### Request Body

```json
{
  "channelName": "string",
  "requesterId": "string",
  "inputModalities": ["text"],
  "outputModalities": ["text", "audio"]
}
```

### Response

```json
{
  "agent_id": "string",
  "create_ts": number,
  "status": "RUNNING"
}
```

## Remove Agent

Removes an AI agent from a channel.

### Endpoint

`POST /agent/remove`

### Request Body

```json
{
  "channelName": "string",
  "requesterId": "string",
  "agentId": "string"
}
```

### Response

```json
{
  "success": boolean,
  "agent_id": "string"
}
```
