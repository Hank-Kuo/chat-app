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
import AddChannelModal from "./addChannelModal";
import ReplyModal from "./replyModal";
import MessageContainer from "./messageContainer";
import getHeader from "../../lib/utils/header";
import { getSession, commitSession } from "../../lib/utils/session";
import { useMessage } from "../../context/messageContext";
import {
  getAllChannelsAPI,
  getUserChannelsAPI,
  addChannelAPI,
  joinChannelAPI,
  ChannelType,
} from "../../apis/channel";
import { S } from "./styles";

export const meta: MetaFunction = () => {
  return [
    { title: "Chat App" },
    { name: "description", content: "Welcome to Chat app!" },
  ];
};

enum FormAction {
  JOIN_CHANNEL_ACTION = "joinAction",
  ADD_CHANNEL_ACTION = "addAction",
}

export async function loader({ request }: LoaderFunctionArgs) {
  const session = await getSession(request.headers.get("Cookie"));
  const userInfo = session.get("userInfo");
  if (!userInfo || userInfo.token.length === 0) {
    return redirect("/login");
  }

  const channelsRes = await getAllChannelsAPI(getHeader(userInfo.token));

  let channels: ChannelType[] = [];
  let userChannels: ChannelType[] = [];
  let err = "";

  if (channelsRes["status"] === "success") {
    channels = channelsRes["data"] as ChannelType[];
  } else {
    err = channelsRes["message"];
  }

  const userChannelsRes = await getUserChannelsAPI(getHeader(userInfo.token));
  if (userChannelsRes["status"] === "success") {
    userChannels = userChannelsRes["data"];
  } else {
    err = userChannelsRes["message"];
  }

  return json(
    {
      channels: channels,
      userChannels: userChannels,
      userInfo: userInfo,
      error: err,
    },
    {
      headers: {
        "Set-Cookie": await commitSession(session),
      },
    }
  );
}

export async function action({ request }: ActionFunctionArgs) {
  const session = await getSession(request.headers.get("Cookie"));
  const userInfo = session.get("userInfo");
  if (!userInfo || userInfo.token.length === 0) {
    return redirect("/login");
  }

  const header = getHeader(userInfo.token);

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
  const { userInfo, channels, userChannels } = useLoaderData<typeof loader>();
  const [selectChannel, setSelectChannel] = React.useState("");
  const [showAddChannel, setShowAddChannel] = React.useState(false);
  const [selectMessage, setSelectMessage] = React.useState(-1);
  const [showReply, setShowReply] = React.useState(false);
  const messageContext = useMessage();

  React.useEffect(() => {
    messageContext.setUserInfo(userInfo);
  }, []);

  return (
    <S.Wrapper>
      <Sidebar
        channels={channels}
        selectChannel={selectChannel}
        setSelectChannel={(value) => {
          setSelectChannel(value);
        }}
        setShow={setShowAddChannel}
      />
      <S.Container>
        {selectChannel === "" ? null : (
          <>
            <S.ContainerTitle>
              {
                channels.filter((v: ChannelType) => v.id === selectChannel)[0]
                  .name
              }
            </S.ContainerTitle>
            <MessageContainer
              userChannels={userChannels}
              selectChannel={selectChannel}
              setSelectMessage={setSelectMessage}
              setShowReply={setShowReply}
            />
          </>
        )}
      </S.Container>
      <ReplyModal
        show={showReply}
        setShow={setShowReply}
        selectChannel={selectChannel}
        selectMessage={selectMessage}
        setSelectMessage={setSelectMessage}
      />
      <AddChannelModal show={showAddChannel} setShow={setShowAddChannel} />
    </S.Wrapper>
  );
}
