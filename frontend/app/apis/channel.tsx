export interface ChannelType {
  id: string;
  name: string;
  createdAt: string;
}

interface addChannelPayload {
  name: string;
}

interface joinChannelPayload {
  channel_id: string;
}

export async function getAllChannelsAPI(header: HeadersInit) {
  const res = await fetch(`http://localhost:9000/api/channel`, {
    headers: header,
  });

  return await res.json();
}

export async function getUserChannelsAPI(header: HeadersInit) {
  const res = await fetch(`http://localhost:9000/api/user/channel`, {
    headers: header,
  });

  return await res.json();
}

export async function addChannelAPI(
  data: addChannelPayload,
  header: HeadersInit
) {
  const res = await fetch(`http://localhost:9000/api/channel`, {
    method: "POST",
    body: JSON.stringify(data),
    headers: header,
  });

  return await res.json();
}

export async function joinChannelAPI(
  data: joinChannelPayload,
  header: HeadersInit
) {
  const res = await fetch(`http://localhost:9000/api/channel/join`, {
    method: "POST",
    body: JSON.stringify(data),
    headers: header,
  });

  return await res.json();
}
