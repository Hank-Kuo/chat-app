import React, { Suspense } from "react";
import { useLoaderData } from "@remix-run/react";

import { loader } from "./index";
import getHeader from "../../lib/utils/header";
import { convertDate } from "../../lib/utils/date";
import useInfiniteScroll from "../../lib/hook/useInfiniteScroll";
import useClickOutside from "../../lib/hook/useClickOutside";
import { useMessage } from "../../context/messageContext";
import { getRepliesAPI, ReplyType } from "../../apis/message";
import { S, M } from "./styles";

interface ReplyModalProps {
  show: boolean;
  setShow: React.Dispatch<React.SetStateAction<boolean>>;
  selectChannel: string;
  selectMessage: number;
  setSelectMessage: React.Dispatch<React.SetStateAction<number>>;
}

const ReplyModal = (props: ReplyModalProps) => {
  const { userInfo } = useLoaderData<typeof loader>();
  const modalRef = React.useRef(null);
  const [text, setText] = React.useState("");
  const [replies, setReplies] = React.useState<ReplyType[]>([]);
  const [hasMore, setHasMore] = React.useState(false);
  const [nextCursor, setNextCusror] = React.useState("");
  const messageContext = useMessage();

  const fetchData = () => {
    getRepliesAPI(
      { messageId: props.selectMessage, nextCursor: nextCursor },
      getHeader(userInfo.token)
    ).then((v) => {
      setReplies((prev) => {
        return [...prev, ...v.data.replies];
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
    if (props.selectMessage > 0) {
      getRepliesAPI(
        { messageId: props.selectMessage },
        getHeader(userInfo.token)
      ).then((v) => {
        setReplies(v.data.replies);
        if (v.data.next_cursor.length != 0) {
          setNextCusror(v.data.next_cursor);
          setHasMore(true);
        }
      });
    }
  }, [props.show]);

  useClickOutside(modalRef, () => {
    props.setShow(false);
    props.setSelectMessage(-1);
  });

  React.useEffect(() => {
    if (messageContext.isReady) {
      if (messageContext.reply["message_id"] === props.selectMessage) {
        setReplies((prev) => {
          if (prev.length > 0) {
            if (messageContext.reply["reply_id"] >= prev[0]["reply_id"]) {
              return [messageContext.reply, ...prev];
            }
          }
          return [messageContext.reply, ...prev];
        });
      }
    }
  }, [messageContext.reply]);

  // submit user text
  const handleClick = () => {
    if (messageContext.isReady) {
      let data = JSON.stringify({
        action: "CreateReply",
        data: {
          channel_id: props.selectChannel,
          message_id: props.selectMessage,
          user_id: userInfo.id,
          username: userInfo.name,
          content: text,
        },
      });
      messageContext.wsClient?.send(data);
      setText("");
    }
  };

  return (
    <Suspense fallback={<div>Loading...</div>}>
      {props.show ? (
        <M.Wrapper>
          <M.Container ref={modalRef}>
            <M.Title>Reply</M.Title>
            <M.ReplyBox>
              {replies.map((v, i) => {
                const createDate = convertDate(v.created_at);
                return (
                  <S.Item
                    id={`${v.reply_id}`}
                    key={`${v.reply_id}`}
                    ref={
                      i === replies.length - 1
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
                    </S.ItemBox>
                  </S.Item>
                );
              })}

              <div>{loading && "Loading..."}</div>
            </M.ReplyBox>
            <S.InputBox>
              <S.Input
                name="reply"
                value={text}
                onChange={(e) => {
                  setText(e.target.value);
                }}
              />
              <S.SubmitBtn onClick={handleClick}>Submit</S.SubmitBtn>
            </S.InputBox>
          </M.Container>
        </M.Wrapper>
      ) : null}
    </Suspense>
  );
};

export default ReplyModal;
