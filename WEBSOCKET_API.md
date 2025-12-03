# WebSocket 실시간 채팅 API 문서

## 개요

실시간 채팅 기능을 위한 WebSocket API입니다. 메시지 전송, 읽음 처리, 멤버 추가/제거가 모두 실시간으로 처리됩니다.

---

## WebSocket 연결

### WS /ws/chat/:roomId

채팅방에 WebSocket으로 연결합니다.

#### 연결 URL

```
ws://localhost:3000/ws/chat/:roomId?user_id=USER_ID
```

#### Parameters

| 파라미터 | 위치 | 타입 | 필수 | 설명 |
|----------|------|------|------|------|
| roomId | Path | uint | O | 채팅방 ID |
| user_id | Query | uint | O | 사용자 ID |

#### 연결 예시

**JavaScript (브라우저):**
```javascript
const roomId = 1;
const userId = 123;
const ws = new WebSocket(`ws://localhost:3000/ws/chat/${roomId}?user_id=${userId}`);

ws.onopen = () => {
  console.log('WebSocket 연결됨');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('받은 메시지:', message);

  // 메시지 타입별 처리
  switch(message.type) {
    case 'message':
      // 새 메시지 수신
      handleNewMessage(message.data);
      break;
    case 'read':
      // 읽음 처리
      handleReadReceipt(message.data);
      break;
    case 'member_join':
      // 멤버 추가
      handleMemberJoin(message.data);
      break;
    case 'member_leave':
      // 멤버 제거
      handleMemberLeave(message.data);
      break;
    case 'member_online':
      // 멤버 접속
      handleMemberOnline(message.data);
      break;
    case 'member_offline':
      // 멤버 접속 종료
      handleMemberOffline(message.data);
      break;
  }
};

ws.onerror = (error) => {
  console.error('WebSocket 에러:', error);
};

ws.onclose = () => {
  console.log('WebSocket 연결 종료');
};
```

**Node.js:**
```javascript
const WebSocket = require('ws');

const roomId = 1;
const userId = 123;
const ws = new WebSocket(`ws://localhost:3000/ws/chat/${roomId}?user_id=${userId}`);

ws.on('open', () => {
  console.log('WebSocket 연결됨');
});

ws.on('message', (data) => {
  const message = JSON.parse(data);
  console.log('받은 메시지:', message);
});
```

---

## 메시지 타입

### 1. message (새 메시지)

누군가가 메시지를 보냈을 때 수신되는 메시지입니다.

**구조:**
```json
{
  "type": "message",
  "room_id": 1,
  "user_id": 2,
  "data": {
    "id": 15,
    "chat_room_id": 1,
    "user_id": 2,
    "message": "안녕하세요!",
    "message_type": "text",
    "file_url": null,
    "is_read": false,
    "created_at": "2024-11-13T16:30:00Z",
    "updated_at": "2024-11-13T16:30:00Z",
    "user": {
      "id": 2,
      "email": "user2@example.com",
      "name": "김철수"
    }
  }
}
```

**처리 예시:**
```javascript
function handleNewMessage(messageData) {
  // UI에 메시지 추가
  const messageElement = createMessageElement(messageData);
  chatContainer.appendChild(messageElement);

  // 스크롤을 맨 아래로
  chatContainer.scrollTop = chatContainer.scrollHeight;
}
```

### 2. read (읽음 처리)

누군가가 메시지를 읽었을 때 수신되는 메시지입니다.

**구조:**
```json
{
  "type": "read",
  "room_id": 1,
  "user_id": 3,
  "data": {
    "user_id": 3,
    "last_read_at": "2024-11-13T16:35:00Z"
  }
}
```

**처리 예시:**
```javascript
function handleReadReceipt(readData) {
  // 읽음 표시 업데이트
  updateReadStatus(readData.user_id, readData.last_read_at);
}
```

### 3. member_join (멤버 추가)

새로운 멤버가 채팅방에 추가되었을 때 수신되는 메시지입니다.

**구조:**
```json
{
  "type": "member_join",
  "room_id": 1,
  "user_id": 5,
  "data": {
    "id": 10,
    "chat_room_id": 1,
    "user_id": 5,
    "role": "member",
    "joined_at": "2024-11-13T16:40:00Z",
    "user": {
      "id": 5,
      "email": "user5@example.com",
      "name": "이영희"
    }
  }
}
```

**처리 예시:**
```javascript
function handleMemberJoin(memberData) {
  // 멤버 목록에 추가
  addMemberToList(memberData.user);

  // 시스템 메시지 표시
  showSystemMessage(`${memberData.user.name}님이 입장했습니다.`);
}
```

### 4. member_leave (멤버 제거)

멤버가 채팅방에서 제거되었을 때 수신되는 메시지입니다.

**구조:**
```json
{
  "type": "member_leave",
  "room_id": 1,
  "user_id": 5,
  "data": {
    "user_id": 5
  }
}
```

**처리 예시:**
```javascript
function handleMemberLeave(leaveData) {
  // 멤버 목록에서 제거
  removeMemberFromList(leaveData.user_id);

  // 시스템 메시지 표시
  showSystemMessage(`사용자가 퇴장했습니다.`);
}
```

### 5. member_online (멤버 접속)

멤버가 WebSocket에 연결되었을 때 수신되는 메시지입니다.

**구조:**
```json
{
  "type": "member_online",
  "room_id": 1,
  "user_id": 4,
  "data": {
    "user_id": 4,
    "status": "online"
  }
}
```

**처리 예시:**
```javascript
function handleMemberOnline(onlineData) {
  // 멤버 상태를 온라인으로 표시
  updateMemberStatus(onlineData.user_id, 'online');
}
```

### 6. member_offline (멤버 접속 종료)

멤버가 WebSocket 연결을 종료했을 때 수신되는 메시지입니다.

**구조:**
```json
{
  "type": "member_offline",
  "room_id": 1,
  "user_id": 4,
  "data": {
    "user_id": 4,
    "status": "offline"
  }
}
```

**처리 예시:**
```javascript
function handleMemberOffline(offlineData) {
  // 멤버 상태를 오프라인으로 표시
  updateMemberStatus(offlineData.user_id, 'offline');
}
```

---

## 실시간 처리 흐름

### 1. 메시지 전송

```
1. 클라이언트 A → HTTP POST /api/v1/chat/rooms/1/messages
2. 서버 → DB에 메시지 저장
3. 서버 → WebSocket Hub를 통해 브로드캐스트
4. 서버 → 모든 연결된 클라이언트에게 'message' 이벤트 전송
5. 클라이언트 B, C, D → 실시간으로 메시지 수신
```

### 2. 읽음 처리

```
1. 클라이언트 A → HTTP POST /api/v1/chat/rooms/1/read
2. 서버 → DB에 읽음 시간 업데이트
3. 서버 → WebSocket Hub를 통해 브로드캐스트
4. 서버 → 모든 연결된 클라이언트에게 'read' 이벤트 전송
5. 클라이언트 B, C, D → 실시간으로 읽음 표시 업데이트
```

### 3. 멤버 추가

```
1. 클라이언트 A → HTTP POST /api/v1/chat/rooms/1/members
2. 서버 → DB에 멤버 추가
3. 서버 → WebSocket Hub를 통해 브로드캐스트
4. 서버 → 모든 연결된 클라이언트에게 'member_join' 이벤트 전송
5. 클라이언트 B, C, D → 실시간으로 멤버 목록 업데이트
```

### 4. 멤버 제거

```
1. 클라이언트 A → HTTP DELETE /api/v1/chat/rooms/1/members/5
2. 서버 → DB에서 멤버 제거
3. 서버 → WebSocket Hub를 통해 브로드캐스트
4. 서버 → 모든 연결된 클라이언트에게 'member_leave' 이벤트 전송
5. 클라이언트 B, C, D → 실시간으로 멤버 목록 업데이트
```

---

## 전체 예제 (React)

```javascript
import { useEffect, useState, useRef } from 'react';

function ChatRoom({ roomId, userId }) {
  const [messages, setMessages] = useState([]);
  const [members, setMembers] = useState([]);
  const wsRef = useRef(null);

  useEffect(() => {
    // WebSocket 연결
    const ws = new WebSocket(`ws://localhost:3000/ws/chat/${roomId}?user_id=${userId}`);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('채팅방 연결됨');
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);

      switch(message.type) {
        case 'message':
          // 새 메시지 추가
          setMessages(prev => [...prev, message.data]);
          break;

        case 'read':
          // 읽음 표시 업데이트
          console.log(`User ${message.data.user_id} read messages`);
          break;

        case 'member_join':
          // 멤버 추가
          setMembers(prev => [...prev, message.data.user]);
          break;

        case 'member_leave':
          // 멤버 제거
          setMembers(prev => prev.filter(m => m.id !== message.data.user_id));
          break;

        case 'member_online':
          // 온라인 상태 업데이트
          updateMemberStatus(message.data.user_id, 'online');
          break;

        case 'member_offline':
          // 오프라인 상태 업데이트
          updateMemberStatus(message.data.user_id, 'offline');
          break;
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket 에러:', error);
    };

    ws.onclose = () => {
      console.log('채팅방 연결 종료');
    };

    // 컴포넌트 언마운트 시 연결 종료
    return () => {
      ws.close();
    };
  }, [roomId, userId]);

  const sendMessage = async (text) => {
    // HTTP API로 메시지 전송
    const response = await fetch(`http://localhost:3000/api/v1/chat/rooms/${roomId}/messages`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        user_id: userId,
        message: text,
        message_type: 'text',
      }),
    });

    // 성공하면 WebSocket을 통해 자동으로 브로드캐스트됨
  };

  const markAsRead = async () => {
    await fetch(`http://localhost:3000/api/v1/chat/rooms/${roomId}/read`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        user_id: userId,
      }),
    });

    // 성공하면 WebSocket을 통해 자동으로 브로드캐스트됨
  };

  return (
    <div className="chat-room">
      <div className="members">
        <h3>멤버 목록</h3>
        {members.map(member => (
          <div key={member.id}>{member.name}</div>
        ))}
      </div>

      <div className="messages">
        {messages.map(msg => (
          <div key={msg.id} className="message">
            <strong>{msg.user.name}:</strong> {msg.message}
          </div>
        ))}
      </div>

      <button onClick={() => sendMessage('안녕하세요!')}>
        메시지 보내기
      </button>
      <button onClick={markAsRead}>
        읽음 처리
      </button>
    </div>
  );
}
```

---

## 주요 특징

### ✅ 실시간 처리
- **메시지 전송**: 메시지를 보내면 모든 접속 중인 사용자에게 즉시 전달
- **읽음 처리**: 읽음 처리하면 다른 사용자들에게 즉시 표시
- **멤버 관리**: 멤버 추가/제거가 실시간으로 반영

### ✅ 자동 상태 관리
- **접속/종료**: WebSocket 연결/종료 시 자동으로 상태 브로드캐스트
- **채팅방별 격리**: 각 채팅방의 메시지는 해당 채팅방 멤버에게만 전송

### ✅ 안정성
- **연결 검증**: 채팅방 멤버만 WebSocket 연결 가능
- **자동 정리**: 연결 종료 시 자동으로 리소스 정리

---

## 주의사항

1. **인증**: 현재는 user_id를 쿼리 파라미터로 받지만, 실제 환경에서는 JWT 토큰을 사용해야 합니다.

2. **재연결**: WebSocket 연결이 끊어진 경우 자동 재연결 로직을 클라이언트에서 구현해야 합니다.

3. **메시지 순서**: WebSocket은 메시지 순서를 보장하지만, HTTP API로 메시지를 보내는 경우 DB의 created_at으로 정렬합니다.

4. **부하 관리**: 많은 사용자가 접속할 경우 서버 리소스를 고려해야 합니다.

---

## 에러 처리

### WebSocket 연결 실패

- **잘못된 roomId**: 연결 즉시 종료
- **잘못된 user_id**: 연결 즉시 종료
- **멤버가 아님**: 연결 즉시 종료

### 연결 종료

WebSocket 연결이 예기치 않게 종료될 수 있습니다. 클라이언트에서 재연결 로직을 구현하세요:

```javascript
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;

function connectWebSocket() {
  const ws = new WebSocket(`ws://localhost:3000/ws/chat/${roomId}?user_id=${userId}`);

  ws.onclose = () => {
    if (reconnectAttempts < maxReconnectAttempts) {
      reconnectAttempts++;
      console.log(`재연결 시도 ${reconnectAttempts}/${maxReconnectAttempts}`);
      setTimeout(() => connectWebSocket(), 1000 * reconnectAttempts);
    } else {
      console.error('재연결 실패');
    }
  };

  ws.onopen = () => {
    reconnectAttempts = 0; // 성공 시 카운터 리셋
  };

  return ws;
}
```

---

## 성능 최적화

### 1. 메시지 배치 처리

많은 메시지를 받을 때 UI 업데이트를 배치로 처리:

```javascript
let messageQueue = [];
let updateTimer = null;

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  messageQueue.push(message);

  if (!updateTimer) {
    updateTimer = setTimeout(() => {
      setMessages(prev => [...prev, ...messageQueue]);
      messageQueue = [];
      updateTimer = null;
    }, 100); // 100ms마다 배치 업데이트
  }
};
```

### 2. 메시지 페이지네이션

초기 로드 시 최근 메시지만 가져오고, 스크롤 시 더 로드:

```javascript
// 처음에는 최근 50개만
const response = await fetch(
  `http://localhost:3000/api/v1/chat/rooms/1/messages?limit=50&offset=0`
);
```
