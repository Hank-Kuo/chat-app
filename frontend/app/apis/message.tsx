export interface MessageType {
  channel_id: string;
  message_id: number;
  content: string;
  user_id: string;
  username: string;
  created_at: string;
}

export interface ReplyType {
  message_id: number;
  reply_id: number;
  content: string;
  user_id: string;
  username: string;
  created_at: string;
}

interface getMessagesPayload {
  channelId: string;
  nextCursor?: string;
}

interface getRepliesPayload {
  messageId: number;
  nextCursor?: string;
}

export async function getMessagesAPI(
  data: getMessagesPayload,
  header: HeadersInit
) {
  const queryParams = data.nextCursor ? `&cursor=${data.nextCursor}` : "";
  const res = await fetch(
    `http://localhost:9000/api/message?channel_id=${data.channelId}${queryParams}`,
    {
      headers: header,
    }
  );

  return await res.json();
}

export async function getRepliesAPI(
  data: getRepliesPayload,
  header: HeadersInit
) {
  const queryParams = data.nextCursor ? `&cursor=${data.nextCursor}` : "";
  const res = await fetch(
    `http://localhost:9000/api/reply?message_id=${data.messageId}${queryParams}`,
    {
      headers: header,
    }
  );

  return await res.json();
}
