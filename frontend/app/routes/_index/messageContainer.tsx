import React, { Suspense } from "react";

import { convertDate } from "../../lib/utils/date";
import { ChannelType } from "../../apis/channel";
import useInfiniteScroll from "../../lib/hook/useInfiniteScroll";
import {
  getMessagesAPI,
  messageAPI123,
  addMessageAPI,
  MessageType,
} from "../../apis/message";
import { S } from "./styles";

interface MessageContainerProps {
  userChannels: ChannelType[];
  selectChannel: string;
  setSelectMessage: React.Dispatch<React.SetStateAction<number>>;
  setShowReply: React.Dispatch<React.SetStateAction<boolean>>;
}

export default function MessageContainer(props: MessageContainerProps) {
  const [text, setText] = React.useState("");
  const [messages, setMessages] = React.useState<MessageType[]>([]);
  const [showJoin, setShowJoin] = React.useState(true);
  const [hasMore, setHasMore] = React.useState(false);
  const [nextCursor, setNextCusror] = React.useState("");

  React.useEffect(() => {
    if (
      props.selectChannel.length !== 0 &&
      props.userChannels.find((v) => v.id === props.selectChannel)
    ) {
      setShowJoin(false);

      messageAPI123.get({ channelId: props.selectChannel }).then((v) => {
        console.log(v);
      });

      getMessagesAPI({ channelId: props.selectChannel }, new Headers()).then(
        (v) => {
          setMessages(v.data.messages);
          if (v.data.next_cursor.length != 0) {
            setNextCusror(v.data.next_cursor);
            setHasMore(true);
          }
        }
      );
    } else {
      setShowJoin(true);
    }
  }, [props.selectChannel]);

  const fetchData = () => {
    getMessagesAPI(
      { channelId: props.selectChannel, nextCursor: nextCursor },
      new Headers()
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
  // submit user text
  const handleClick = () => {
    addMessageAPI(
      {
        channel_id: props.selectChannel,
        user_id: "257e4caf-fb4b-43a1-a4b3-cca94f583bd5",
        username: "hank",
        content: text,
      },
      new Headers()
    ).then((v) => {
      if (v["status"] === "success") {
        const message: MessageType = v["data"];
        setMessages((prev) => {
          return [message, ...prev];
        });
        setText("");
      }
    });
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
                  <S.ItemName>{showJoin ? "Default" : v.username}</S.ItemName>
                  <S.ItemTime>
                    {showJoin ? "2023/01/01" : createDate}
                  </S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>
                    {showJoin ? "⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆" : v.content}
                  </S.ItemMessage>
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
