import { createCookieSessionStorage } from "@remix-run/node";

type SessionData = {
  userInfo: string;
};

type SessionFlashData = {
  error: string;
};

const { getSession, commitSession, destroySession } =
  createCookieSessionStorage<SessionData, SessionFlashData>({
    cookie: {
      name: "__session",
      expires: new Date(Date.now() + 60_000),
      maxAge: 60,
      secrets: ["chat-app"],
      secure: true,
    },
  });

export { getSession, commitSession, destroySession };
