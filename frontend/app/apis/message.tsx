interface getMessagesPayload {
  channelId: string;
}

interface getRepliesPayload {
  channelId: string;
  messageId: string;
}

interface addReplyPayload {
  channelId: string;
  messageId: string;
  userId: string;
  content: string;
  createdAt: string;
}
interface addMessagePayload {
  channelId: string;
  userId: string;
  content: string;
  createdAt: string;
}

export async function getMessagesAPI(data: getMessagesPayload) {
  return {
    status: true,
    messages: [
      {
        id: 1,
        name: "hank",
        creatdAt: "2023/10/1",
        message:
          "message message message message message message message message message message message message message message message message message message message message message message",
      },
      {
        id: 2,
        name: "hank",
        creatdAt: "2023/10/1",
        message: "message",
      },
      {
        id: 3,
        name: "hank",
        creatdAt: "2023/10/1",
        message: "message",
      },
      {
        id: 4,
        name: "hank",
        creatdAt: "2023/10/1",
        message: "message",
      },
      { id: 5, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 6, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 7, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 8, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 9, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 10, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 11, name: "hank", creatdAt: "2023/10/1", message: "message" },
      {
        id: 12,
        name: "hank",
        creatdAt: "2023/10/1",
        message:
          "message message message message message message message message message message message message message message message message message message message message message message",
      },
    ],
    errMessage: "",
  };
}

export async function getRepliesAPI(data: getRepliesPayload) {
  return {
    status: true,
    messages: [
      {
        id: 1,
        name: "hank",
        creatdAt: "2023/10/1",
        message:
          "message message message message message message message message message message message message message message message message message message message message message message",
      },
      {
        id: 2,
        name: "hank",
        creatdAt: "2023/10/1",
        message: "message",
      },
      {
        id: 3,
        name: "hank",
        creatdAt: "2023/10/1",
        message: "message",
      },
      {
        id: 4,
        name: "hank",
        creatdAt: "2023/10/1",
        message: "message",
      },
      { id: 5, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 6, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 7, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 8, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 9, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 10, name: "hank", creatdAt: "2023/10/1", message: "message" },
      { id: 11, name: "hank", creatdAt: "2023/10/1", message: "message" },
      {
        id: 12,
        name: "hank",
        creatdAt: "2023/10/1",
        message:
          "message message message message message message message message message message message message message message message message message message message message message message",
      },
    ],
    message: "",
  };
}

export async function addMessageAPI(data: addMessagePayload) {
  return {
    status: true,
    errMessage: "",
  };
}

export async function addReplyAPI(data: addReplyPayload) {
  return {
    status: true,
    message: "",
  };
}
