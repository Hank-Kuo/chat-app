import React from "react";
import {
  MetaFunction,
  json,
  redirect,
  ActionFunctionArgs,
  LoaderFunctionArgs,
} from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";

import Sidebar from "../../components/sidebar";
import ChannelModal from "./channelModal";
import ReplyModal from "./replyModal";
import MessageContainer from "./messageContainer";
import { getSession, commitSession } from "../../lib/utils/session";
import {
  getAllChannelsAPI,
  getUserChannelsAPI,
  addChannelAPI,
  joinChannelAPI,
} from "../../apis/channel";
import { S } from "./styles";

export const meta: MetaFunction = () => {
  return [
    { title: "Chat App" },
    { name: "description", content: "Welcome to Chat app!" },
  ];
};

interface UserInfo {
  id: string;
  name: string;
  email: string;
  token: string;
  created_at: string;
}
export async function loader({ request }: LoaderFunctionArgs) {
  const session = await getSession(request.headers.get("Cookie"));
  const userInfoStr = session.get("userInfo") as string;
  if (!userInfoStr || userInfoStr.length === 0) {
    return redirect("/login");
  }

  const userInfo: UserInfo = JSON.parse(userInfoStr);
  const header = new Headers();

  header.append("Content-Type", "application/json");
  header.append("Authorization", `Bearer ${userInfo.token}`);
  const channelsRes = await getAllChannelsAPI(header);

  let channels = [];
  let userChannels = [];
  let err = "";
  if (channelsRes["status"] === "success") {
    channels = channelsRes["data"];
  } else {
    err = channelsRes["message"];
  }

  const userChannelsRes = await getUserChannelsAPI(header);
  if (userChannelsRes["status"] === "success") {
    userChannels = userChannelsRes["data"];
  } else {
    err = userChannelsRes["message"];
  }

  return json(
    {
      channels: channels,
      userChannels: userChannels,
      error: err,
    },
    {
      headers: {
        "Set-Cookie": await commitSession(session),
      },
    }
  );
}

enum FormAction {
  JOIN_CHANNEL_ACTION = "joinAction",
  ADD_CHANNEL_ACTION = "addAction",
}

export async function action({ request }: ActionFunctionArgs) {
  const session = await getSession(request.headers.get("Cookie"));
  const userInfoStr = session.get("userInfo") as string;
  if (userInfoStr.length === 0) {
    return redirect("/login");
  }

  const userInfo: UserInfo = JSON.parse(userInfoStr);
  const header = new Headers();

  header.append("Content-Type", "application/json");
  header.append("Authorization", `Bearer ${userInfo.token}`);

  const form = await request.formData();
  const formAction = form.get("formAction");

  switch (formAction) {
    case FormAction.JOIN_CHANNEL_ACTION:
      const channelID = form.get("channelID") as string;
      const joinChannelRes = await joinChannelAPI(
        { channel_id: channelID },
        header
      );

      if (joinChannelRes["status"] !== "success") {
        return json(
          { error: joinChannelRes["message"] },
          {
            headers: {
              "Set-Cookie": await commitSession(session),
            },
          }
        );
      } else {
        return {};
      }

    case FormAction.ADD_CHANNEL_ACTION:
      const channelName = form.get("channelName") as string;
      const addChannelR = await addChannelAPI({ name: channelName }, header);

      if (addChannelR["status"] !== "success") {
        return json(
          { error: addChannelR["message"] },
          {
            headers: {
              "Set-Cookie": await commitSession(session),
            },
          }
        );
      } else {
        return {};
      }

    default:
      return {};
  }
}

export default function Index() {
  const { channels, userChannels } = useLoaderData<typeof loader>();
  const [selectChannel, setSelectChannel] = React.useState("");
  const [showChannel, setShowChannel] = React.useState(false);
  const [showReply, setShowReply] = React.useState(false);
  const title = channels.filter((v: any) => v.id === selectChannel);

  return (
    <S.Wrapper>
      <Sidebar
        selectChannel={selectChannel}
        setSelectChannel={(value) => {
          setSelectChannel(value);
        }}
        setShow={setShowChannel}
      />
      <S.Container>
        {selectChannel === "" ? null : (
          <>
            <S.ContainerTitle>
              {title.length > 0 ? title[0].name : ""}
            </S.ContainerTitle>
            <MessageContainer
              userChannels={userChannels}
              selectChannel={selectChannel}
              setShowReply={setShowReply}
            />
          </>
        )}
      </S.Container>
      <ReplyModal show={showReply} setShow={setShowReply} />
      <ChannelModal show={showChannel} setShow={setShowChannel} />
    </S.Wrapper>
  );
}
