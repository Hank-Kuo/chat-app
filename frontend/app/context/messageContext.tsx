import React from "react";
import { UserInfo } from "~/lib/utils/session";
import { MessageType, ReplyType } from "~/apis/message";

interface MessageContextType {
  user: UserInfo;
  isReady: boolean;
  wsClient: WebSocket | null;
  message: MessageType;
  reply: ReplyType;
  // setUser: React.Dispatch<React.SetStateAction<UserInfo>>;
  setUserInfo: (userData: UserInfo) => void;
}

const MessageContext = React.createContext<MessageContextType>(null!);

export function MessageProvider({ children }: { children: React.ReactNode }) {
  const [isReady, setIsReady] = React.useState(false);
  const [message, setMessage] = React.useState<MessageType>({
    channel_id: "",
    message_id: 0,
    content: "",
    user_id: "",
    username: "",
    created_at: "",
  });

  const [reply, setReply] = React.useState<ReplyType>({
    reply_id: 0,
    message_id: 0,
    content: "",
    user_id: "",
    username: "",
    created_at: "",
  });

  const [wsClient, setWsClient] = React.useState<WebSocket | null>(null);

  const [isLogin, setIsLogin] = React.useState(false);
  let [user, setUser] = React.useState<UserInfo>({
    id: "",
    name: "",
    email: "",
    token: "",
    created_at: "",
  });

  React.useEffect(() => {
    if (isLogin) {
      const ws = new WebSocket("ws://127.0.0.1:9000/api/message/ws");
      ws.onopen = () => {
        ws.send("Authorization: Bearer " + user.token);
        setIsReady(true);
      };
      ws.onclose = () => {
        setIsReady(false);
        alert("Session disconnect");
      };

      ws.onmessage = function (event) {
        if (event.data.length > 0) {
          try {
            let data = JSON.parse(event.data);

            // check is message or reply
            if (data["status"] === "success") {
              if (typeof data["data"]["reply_id"] === "number") {
                let reply: ReplyType = data["data"];
                setReply(reply);
              } else if (typeof data["data"]["channel_id"] === "string") {
                let message: MessageType = data["data"];
                setMessage(message);
              }
            } else {
              alert("Got some error");
            }
          } catch {
            alert("Got some error");
          }
        }
      };

      setWsClient(ws);
    }
  }, [user]);

  const setUserInfo = (userData: UserInfo) => {
    setUser(userData);
    setIsLogin(true);
  };

  let value = { user, setUserInfo, isReady, wsClient, message, reply };

  return (
    <MessageContext.Provider value={value}>{children}</MessageContext.Provider>
  );
}

export function useMessage() {
  return React.useContext(MessageContext);
}
