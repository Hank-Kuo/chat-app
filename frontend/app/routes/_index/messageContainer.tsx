import React, { Suspense } from "react";
import { useLoaderData } from "@remix-run/react";

import { loader } from "./index";
import getHeader from "../../lib/utils/header";
import { convertDate } from "../../lib/utils/date";
import useInfiniteScroll from "../../lib/hook/useInfiniteScroll";
import { useMessage } from "../../context/messageContext";
import { ChannelType } from "../../apis/channel";
import { getMessagesAPI, MessageType } from "../../apis/message";
import { S } from "./styles";

interface MessageContainerProps {
  userChannels: ChannelType[];
  selectChannel: string;
  setSelectMessage: React.Dispatch<React.SetStateAction<number>>;
  setShowReply: React.Dispatch<React.SetStateAction<boolean>>;
}

export default function MessageContainer(props: MessageContainerProps) {
  const { userInfo } = useLoaderData<typeof loader>();
  const [text, setText] = React.useState("");
  const [messages, setMessages] = React.useState<MessageType[]>([]);
  const [showJoin, setShowJoin] = React.useState(true);
  const [hasMore, setHasMore] = React.useState(false);
  const [nextCursor, setNextCusror] = React.useState("");
  const messageContext = useMessage();

  React.useEffect(() => {
    if (
      props.selectChannel.length !== 0 &&
      props.userChannels.find((v) => v.id === props.selectChannel)
    ) {
      setShowJoin(false);

      getMessagesAPI(
        { channelId: props.selectChannel },
        getHeader(userInfo.token)
      ).then((v) => {
        setMessages(v.data.messages);
        if (v.data.next_cursor.length != 0) {
          setNextCusror(v.data.next_cursor);
          setHasMore(true);
        }
      });
    } else {
      setShowJoin(true);
      setMessages([]);
    }
  }, [props.selectChannel, props.userChannels]);

  const fetchData = () => {
    getMessagesAPI(
      { channelId: props.selectChannel, nextCursor: nextCursor },
      getHeader(userInfo.token)
    ).then((v) => {
      setMessages((prev) => {
        return [...prev, ...v.data.messages];
      });

      if (v.data.next_cursor.length != 0) {
        setNextCusror(v.data.next_cursor);
        setHasMore(true);
      } else {
        setHasMore(false);
      }
    });
  };

  const [lastElementRef, loading] = useInfiniteScroll(hasMore, fetchData);

  React.useEffect(() => {
    if (messageContext.isReady) {
      if (messageContext.message["channel_id"] === props.selectChannel) {
        setMessages((prev) => {
          if (prev.length > 0) {
            if (messageContext.message.message_id >= prev[0]["message_id"]) {
              return [messageContext.message, ...prev];
            }
          }
          return [messageContext.message, ...prev];
        });
      }
    }
  }, [messageContext.message]);

  // submit user text
  const handleClick = () => {
    if (messageContext.isReady) {
      let data = JSON.stringify({
        action: "CreateMessage",
        data: {
          channel_id: props.selectChannel,
          user_id: userInfo.id,
          username: userInfo.name,
          content: text,
        },
      });
      messageContext.wsClient?.send(data);
      setText("");
    }
  };

  // show reply modal
  const handleReplyClick = (id: number) => {
    if (showJoin) {
      return;
    }
    props.setShowReply(true);
    props.setSelectMessage(id);
  };

  return (
    <>
      <Suspense fallback={<div>Loading...</div>}>
        <S.Box>
          {messages.map((v, i) => {
            const createDate = convertDate(v.created_at);
            return (
              <S.Item
                id={`${v.message_id}`}
                key={v.message_id}
                ref={
                  i === messages.length - 1
                    ? (lastElementRef as React.LegacyRef<HTMLDivElement>)
                    : null
                }
              >
                <S.ItemInfoBox>
                  <S.ItemName>{v.username}</S.ItemName>
                  <S.ItemTime>{createDate}</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>{v.content}</S.ItemMessage>
                  <S.ItemReply onClick={() => handleReplyClick(v.message_id)}>
                    REPLY
                  </S.ItemReply>
                </S.ItemBox>
              </S.Item>
            );
          })}
          <div>{loading && "Loading..."}</div>
        </S.Box>
        <S.InputBox>
          <S.Input
            name="message"
            value={text}
            onChange={(e) => {
              setText(e.target.value);
            }}
          />
          <S.SubmitBtn onClick={handleClick}>Submit</S.SubmitBtn>
        </S.InputBox>

        {/* Join button if user not join this channel */}
        {showJoin ? (
          <S.JoinBox method="post" onSubmit={() => setShowJoin(false)}>
            <S.JoinInput
              name="channelID"
              value={props.selectChannel}
              onChange={() => {}}
            />
            <S.JoinBtn name="formAction" value="joinAction">
              Join
            </S.JoinBtn>
            <S.JoinGreyBox />
          </S.JoinBox>
        ) : null}
      </Suspense>
    </>
  );
}
