# 그룹 채팅 API 문서

## 목차
1. [채팅방 생성](#채팅방-생성)
2. [채팅방 목록 조회](#채팅방-목록-조회)
3. [채팅방 상세 조회](#채팅방-상세-조회)
4. [메시지 전송](#메시지-전송)
5. [메시지 목록 조회](#메시지-목록-조회)
6. [메시지 읽음 처리](#메시지-읽음-처리)
7. [멤버 추가](#멤버-추가)
8. [멤버 제거](#멤버-제거)

---

## 채팅방 생성

### POST /api/v1/chat/rooms

새로운 그룹 채팅방을 생성합니다.

#### Request

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

**Body:**
```json
{
  "name": "등산 동호회 채팅방",
  "description": "주말 등산을 함께하는 사람들의 채팅방",
  "club_id": 1,
  "room_type": "club",
  "member_ids": [2, 3, 4, 5]
}
```

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| name | string | O | 채팅방 이름 |
| description | string | X | 채팅방 설명 |
| club_id | uint | X | 클럽 ID (클럽 채팅방인 경우) |
| room_type | string | X | 채팅방 타입 (group, club, direct) - 기본값: group |
| member_ids | []uint | X | 초대할 멤버 ID 목록 |

#### Response

**성공 (201 Created):**
```json
{
  "success": true,
  "message": "Chat room created successfully",
  "data": {
    "id": 1,
    "name": "등산 동호회 채팅방",
    "description": "주말 등산을 함께하는 사람들의 채팅방",
    "club_id": 1,
    "room_type": "club",
    "created_by": 1,
    "member_count": 5,
    "last_message": null,
    "last_message_at": null,
    "created_at": "2024-11-13T10:00:00Z",
    "updated_at": "2024-11-13T10:00:00Z",
    "creator": {
      "id": 1,
      "email": "user1@example.com",
      "name": "홍길동"
    },
    "members": [
      {
        "id": 1,
        "chat_room_id": 1,
        "user_id": 1,
        "role": "admin",
        "joined_at": "2024-11-13T10:00:00Z",
        "user": {
          "id": 1,
          "email": "user1@example.com",
          "name": "홍길동"
        }
      },
      {
        "id": 2,
        "chat_room_id": 1,
        "user_id": 2,
        "role": "member",
        "joined_at": "2024-11-13T10:00:00Z",
        "user": {
          "id": 2,
          "email": "user2@example.com",
          "name": "김철수"
        }
      }
    ]
  }
}
```

#### cURL 예제

```bash
curl -X POST http://localhost:3000/api/v1/chat/rooms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "등산 동호회 채팅방",
    "description": "주말 등산을 함께하는 사람들의 채팅방",
    "club_id": 1,
    "room_type": "club",
    "member_ids": [2, 3, 4, 5]
  }'
```

---

## 채팅방 목록 조회

### GET /api/v1/chat/rooms

사용자가 속한 모든 채팅방 목록을 조회합니다.

#### Request

**Headers:**
```
Authorization: Bearer {JWT_TOKEN}
```

**Query Parameters:**
| 파라미터 | 타입 | 필수 | 설명 |
|----------|------|------|------|
| user_id | uint | O | 사용자 ID (현재는 쿼리로 전달, 추후 JWT에서 추출) |

#### Response

**성공 (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "등산 동호회 채팅방",
      "description": "주말 등산을 함께하는 사람들의 채팅방",
      "club_id": 1,
      "room_type": "club",
      "created_by": 1,
      "member_count": 5,
      "last_message": "다음 주말에 만나요!",
      "last_message_at": "2024-11-13T15:30:00Z",
      "created_at": "2024-11-13T10:00:00Z",
      "updated_at": "2024-11-13T15:30:00Z",
      "creator": {
        "id": 1,
        "email": "user1@example.com",
        "name": "홍길동"
      },
      "club": {
        "id": 1,
        "name": "등산 동호회",
        "description": "주말마다 산을 오르는 모임"
      }
    }
  ]
}
```

#### cURL 예제

```bash
curl -X GET "http://localhost:3000/api/v1/chat/rooms?user_id=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 채팅방 상세 조회

### GET /api/v1/chat/rooms/:id

특정 채팅방의 상세 정보를 조회합니다.

#### Request

**Headers:**
```
Authorization: Bearer {JWT_TOKEN}
```

**Path Parameters:**
| 파라미터 | 타입 | 필수 | 설명 |
|----------|------|------|------|
| id | uint | O | 채팅방 ID |

#### Response

**성공 (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "등산 동호회 채팅방",
    "description": "주말 등산을 함께하는 사람들의 채팅방",
    "club_id": 1,
    "room_type": "club",
    "created_by": 1,
    "member_count": 5,
    "last_message": "다음 주말에 만나요!",
    "last_message_at": "2024-11-13T15:30:00Z",
    "created_at": "2024-11-13T10:00:00Z",
    "updated_at": "2024-11-13T15:30:00Z",
    "creator": {
      "id": 1,
      "email": "user1@example.com",
      "name": "홍길동"
    },
    "club": {
      "id": 1,
      "name": "등산 동호회"
    },
    "members": [
      {
        "id": 1,
        "chat_room_id": 1,
        "user_id": 1,
        "role": "admin",
        "joined_at": "2024-11-13T10:00:00Z",
        "last_read_at": "2024-11-13T15:30:00Z",
        "unread_count": 0,
        "user": {
          "id": 1,
          "email": "user1@example.com",
          "name": "홍길동"
        }
      }
    ]
  }
}
```

**실패 (404 Not Found):**
```json
{
  "success": false,
  "error": "Chat room not found"
}
```

#### cURL 예제

```bash
curl -X GET http://localhost:3000/api/v1/chat/rooms/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 메시지 전송

### POST /api/v1/chat/rooms/:id/messages

채팅방에 메시지를 전송합니다.

#### Request

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

**Path Parameters:**
| 파라미터 | 타입 | 필수 | 설명 |
|----------|------|------|------|
| id | uint | O | 채팅방 ID |

**Body:**
```json
{
  "user_id": 1,
  "message": "안녕하세요! 다음 주말에 등산 가실 분?",
  "message_type": "text",
  "file_url": ""
}
```

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| user_id | uint | O | 발신자 ID (추후 JWT에서 추출) |
| message | string | O | 메시지 내용 |
| message_type | string | X | 메시지 타입 (text, image, file, system) - 기본값: text |
| file_url | string | X | 파일/이미지 URL |

#### Response

**성공 (201 Created):**
```json
{
  "success": true,
  "message": "Message sent successfully",
  "data": {
    "id": 1,
    "chat_room_id": 1,
    "user_id": 1,
    "message": "안녕하세요! 다음 주말에 등산 가실 분?",
    "message_type": "text",
    "file_url": null,
    "is_read": false,
    "created_at": "2024-11-13T15:30:00Z",
    "updated_at": "2024-11-13T15:30:00Z",
    "user": {
      "id": 1,
      "email": "user1@example.com",
      "name": "홍길동"
    }
  }
}
```

**실패 (403 Forbidden):**
```json
{
  "success": false,
  "error": "User is not a member of this chat room"
}
```

**실패 (404 Not Found):**
```json
{
  "success": false,
  "error": "Chat room not found"
}
```

#### cURL 예제

```bash
curl -X POST http://localhost:3000/api/v1/chat/rooms/1/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "user_id": 1,
    "message": "안녕하세요! 다음 주말에 등산 가실 분?",
    "message_type": "text"
  }'
```

---

## 메시지 목록 조회

### GET /api/v1/chat/rooms/:id/messages

채팅방의 메시지 목록을 조회합니다 (페이지네이션 지원).

#### Request

**Headers:**
```
Authorization: Bearer {JWT_TOKEN}
```

**Path Parameters:**
| 파라미터 | 타입 | 필수 | 설명 |
|----------|------|------|------|
| id | uint | O | 채팅방 ID |

**Query Parameters:**
| 파라미터 | 타입 | 필수 | 기본값 | 설명 |
|----------|------|------|--------|------|
| limit | int | X | 50 | 한 번에 가져올 메시지 수 |
| offset | int | X | 0 | 건너뛸 메시지 수 |

#### Response

**성공 (200 OK):**
```json
{
  "success": true,
  "data": {
    "messages": [
      {
        "id": 3,
        "chat_room_id": 1,
        "user_id": 2,
        "message": "저도 갈게요!",
        "message_type": "text",
        "file_url": null,
        "is_read": false,
        "created_at": "2024-11-13T15:35:00Z",
        "updated_at": "2024-11-13T15:35:00Z",
        "user": {
          "id": 2,
          "email": "user2@example.com",
          "name": "김철수"
        }
      },
      {
        "id": 2,
        "chat_room_id": 1,
        "user_id": 1,
        "message": "다음 주말 북한산 어떠세요?",
        "message_type": "text",
        "file_url": null,
        "is_read": true,
        "created_at": "2024-11-13T15:32:00Z",
        "updated_at": "2024-11-13T15:32:00Z",
        "user": {
          "id": 1,
          "email": "user1@example.com",
          "name": "홍길동"
        }
      },
      {
        "id": 1,
        "chat_room_id": 1,
        "user_id": 1,
        "message": "안녕하세요! 다음 주말에 등산 가실 분?",
        "message_type": "text",
        "file_url": null,
        "is_read": true,
        "created_at": "2024-11-13T15:30:00Z",
        "updated_at": "2024-11-13T15:30:00Z",
        "user": {
          "id": 1,
          "email": "user1@example.com",
          "name": "홍길동"
        }
      }
    ],
    "total": 3,
    "limit": 50,
    "offset": 0
  }
}
```

#### cURL 예제

```bash
# 최근 50개 메시지 조회
curl -X GET "http://localhost:3000/api/v1/chat/rooms/1/messages?limit=50&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 다음 페이지 (51~100번째 메시지)
curl -X GET "http://localhost:3000/api/v1/chat/rooms/1/messages?limit=50&offset=50" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 메시지 읽음 처리

### POST /api/v1/chat/rooms/:id/read

사용자가 채팅방의 메시지를 읽음 처리합니다 (읽지 않은 메시지 수 초기화).

#### Request

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

**Path Parameters:**
| 파라미터 | 타입 | 필수 | 설명 |
|----------|------|------|------|
| id | uint | O | 채팅방 ID |

**Body:**
```json
{
  "user_id": 1
}
```

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| user_id | uint | O | 사용자 ID (추후 JWT에서 추출) |

#### Response

**성공 (200 OK):**
```json
{
  "success": true,
  "message": "Messages marked as read"
}
```

**실패 (404 Not Found):**
```json
{
  "success": false,
  "error": "Membership not found"
}
```

#### cURL 예제

```bash
curl -X POST http://localhost:3000/api/v1/chat/rooms/1/read \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "user_id": 1
  }'
```

---

## 멤버 추가

### POST /api/v1/chat/rooms/:id/members

채팅방에 새로운 멤버를 추가합니다.

#### Request

**Headers:**
```
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

**Path Parameters:**
| 파라미터 | 타입 | 필수 | 설명 |
|----------|------|------|------|
| id | uint | O | 채팅방 ID |

**Body:**
```json
{
  "user_id": 6
}
```

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| user_id | uint | O | 추가할 사용자 ID |

#### Response

**성공 (201 Created):**
```json
{
  "success": true,
  "message": "Member added successfully",
  "data": {
    "id": 6,
    "chat_room_id": 1,
    "user_id": 6,
    "role": "member",
    "joined_at": "2024-11-13T16:00:00Z",
    "last_read_at": null,
    "unread_count": 0,
    "created_at": "2024-11-13T16:00:00Z",
    "user": {
      "id": 6,
      "email": "user6@example.com",
      "name": "이영희"
    }
  }
}
```

**실패 (409 Conflict):**
```json
{
  "success": false,
  "error": "User is already a member"
}
```

**실패 (404 Not Found):**
```json
{
  "success": false,
  "error": "Chat room not found"
}
```

#### cURL 예제

```bash
curl -X POST http://localhost:3000/api/v1/chat/rooms/1/members \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "user_id": 6
  }'
```

---

## 멤버 제거

### DELETE /api/v1/chat/rooms/:id/members/:userId

채팅방에서 멤버를 제거합니다.

#### Request

**Headers:**
```
Authorization: Bearer {JWT_TOKEN}
```

**Path Parameters:**
| 파라미터 | 타입 | 필수 | 설명 |
|----------|------|------|------|
| id | uint | O | 채팅방 ID |
| userId | uint | O | 제거할 사용자 ID |

#### Response

**성공 (200 OK):**
```json
{
  "success": true,
  "message": "Member removed successfully"
}
```

**실패 (404 Not Found):**
```json
{
  "success": false,
  "error": "Member not found"
}
```

#### cURL 예제

```bash
curl -X DELETE http://localhost:3000/api/v1/chat/rooms/1/members/6 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 데이터 모델

### ChatRoom (채팅방)

| 필드 | 타입 | 설명 |
|------|------|------|
| id | uint | 채팅방 ID |
| name | string | 채팅방 이름 |
| description | string | 채팅방 설명 |
| club_id | uint | 클럽 ID (nullable) |
| room_type | string | 채팅방 타입 (group, club, direct) |
| created_by | uint | 생성자 ID |
| member_count | int | 멤버 수 |
| last_message | string | 마지막 메시지 |
| last_message_at | timestamp | 마지막 메시지 시간 |
| created_at | timestamp | 생성 시간 |
| updated_at | timestamp | 수정 시간 |

### ChatRoomMember (채팅방 멤버)

| 필드 | 타입 | 설명 |
|------|------|------|
| id | uint | 멤버십 ID |
| chat_room_id | uint | 채팅방 ID |
| user_id | uint | 사용자 ID |
| role | string | 역할 (admin, member) |
| joined_at | timestamp | 가입 시간 |
| last_read_at | timestamp | 마지막 읽은 시간 |
| unread_count | int | 읽지 않은 메시지 수 |
| created_at | timestamp | 생성 시간 |

### ChatMessage (채팅 메시지)

| 필드 | 타입 | 설명 |
|------|------|------|
| id | uint | 메시지 ID |
| chat_room_id | uint | 채팅방 ID |
| user_id | uint | 발신자 ID |
| message | string | 메시지 내용 |
| message_type | string | 메시지 타입 (text, image, file, system) |
| file_url | string | 파일/이미지 URL (nullable) |
| is_read | bool | 읽음 여부 |
| created_at | timestamp | 생성 시간 |
| updated_at | timestamp | 수정 시간 |

---

## 전체 플로우 예시

### 1. 채팅방 생성 및 멤버 초대
```bash
# 1-1. 채팅방 생성
curl -X POST http://localhost:3000/api/v1/chat/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "name": "등산 동호회 채팅방",
    "description": "주말 등산 모임",
    "club_id": 1,
    "room_type": "club",
    "member_ids": [2, 3, 4]
  }'

# Response: 채팅방 ID 1 생성됨
```

### 2. 메시지 전송
```bash
# 2-1. 첫 메시지 전송
curl -X POST http://localhost:3000/api/v1/chat/rooms/1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "message": "안녕하세요! 다음 주말에 등산 가실 분?"
  }'

# 2-2. 다른 사용자가 답장
curl -X POST http://localhost:3000/api/v1/chat/rooms/1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 2,
    "message": "저도 갈게요!"
  }'
```

### 3. 메시지 조회 및 읽음 처리
```bash
# 3-1. 메시지 목록 조회
curl -X GET "http://localhost:3000/api/v1/chat/rooms/1/messages?limit=50"

# 3-2. 읽음 처리
curl -X POST http://localhost:3000/api/v1/chat/rooms/1/read \
  -H "Content-Type: application/json" \
  -d '{"user_id": 2}'
```

### 4. 멤버 관리
```bash
# 4-1. 멤버 추가
curl -X POST http://localhost:3000/api/v1/chat/rooms/1/members \
  -H "Content-Type: application/json" \
  -d '{"user_id": 5}'

# 4-2. 멤버 제거
curl -X DELETE http://localhost:3000/api/v1/chat/rooms/1/members/5
```

---

## 주요 기능

### ✅ 구현된 기능
- 그룹 채팅방 생성
- 채팅방 목록 조회 (사용자별)
- 채팅방 상세 정보 조회 (멤버 목록 포함)
- 메시지 전송 (텍스트, 이미지, 파일)
- 메시지 목록 조회 (페이지네이션)
- 읽지 않은 메시지 수 관리
- 메시지 읽음 처리
- 채팅방 멤버 추가/제거
- 마지막 메시지 및 시간 자동 업데이트

### 🔜 추후 개선 가능한 기능
- WebSocket을 이용한 실시간 메시지 전송
- 메시지 검색 기능
- 파일 업로드 기능
- 메시지 삭제/수정 기능
- 채팅방 나가기 기능
- 푸시 알림 연동
- 메시지 타입별 필터링

---

## 에러 코드

| HTTP 상태 코드 | 설명 |
|----------------|------|
| 200 OK | 요청 성공 |
| 201 Created | 리소스 생성 성공 |
| 400 Bad Request | 잘못된 요청 |
| 403 Forbidden | 권한 없음 (채팅방 멤버가 아님) |
| 404 Not Found | 리소스를 찾을 수 없음 |
| 409 Conflict | 중복 (이미 멤버임) |
| 500 Internal Server Error | 서버 내부 오류 |
