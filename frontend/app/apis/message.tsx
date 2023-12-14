interface getMessagesPayload {
  channelId: string;
}

interface getRepliesPayload {
  messageId: number;
}

interface addReplyPayload {
  message_id: number;
  user_id: string;
  username: string;
  content: string;
}

interface addMessagePayload {
  channel_id: string;
  userId: string;
  user_id: string;
  username: string;
  content: string;
}

export async function getMessagesAPI(
  data: getMessagesPayload,
  header: HeadersInit
) {
  const res = await fetch(
    `http://localhost:9000/api/message?channel_id=${data.channelId}`,
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
  const res = await fetch(
    `http://localhost:9000/api/reply?message_id=${data.messageId}`,
    {
      headers: header,
    }
  );

  return await res.json();
}

export async function addMessageAPI(
  data: addMessagePayload,
  header: HeadersInit
) {
  const res = await fetch(`http://localhost:9000/api/message`, {
    method: "POST",
    body: JSON.stringify(data),
    headers: header,
  });

  return await res.json();
}

export async function addReplyAPI(data: addReplyPayload, header: HeadersInit) {
  const res = await fetch(`http://localhost:9000/api/reply`, {
    method: "POST",
    body: JSON.stringify(data),
    headers: header,
  });

  return await res.json();
}
