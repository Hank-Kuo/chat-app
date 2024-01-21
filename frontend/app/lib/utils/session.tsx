import { createCookieSessionStorage } from "@remix-run/node";

export interface UserInfo {
  id: string;
  name: string;
  email: string;
  token: string;
  created_at: string;
}

type SessionData = {
  userInfo: UserInfo;
};

type SessionFlashData = {
  error: string;
};

const { getSession, commitSession, destroySession } =
  createCookieSessionStorage<SessionData, SessionFlashData>({
    cookie: {
      name: "__session",
      expires: new Date(Date.now() + 8 * 60 * 60 * 1000),
      secrets: ["chat-app"],
      secure: true,
    },
  });

export { getSession, commitSession, destroySession };
