import React, { Suspense } from "react";

import { convertDate } from "../../lib/utils/date";
import useInfiniteScroll from "../../lib/hook/useInfiniteScroll";
import useClickOutside from "../../lib/hook/useClickOutside";

import { getRepliesAPI, addReplyAPI, ReplyType } from "../../apis/message";
import { S, M } from "./styles";

interface ReplyModalProps {
  show: boolean;
  setShow: React.Dispatch<React.SetStateAction<boolean>>;
  selectMessage: number;
  setSelectMessage: React.Dispatch<React.SetStateAction<number>>;
}

const ReplyModal = (props: ReplyModalProps) => {
  const modalRef = React.useRef(null);
  const [text, setText] = React.useState("");
  const [replies, setReplies] = React.useState<ReplyType[]>([]);

  const [hasMore, setHasMore] = React.useState(true);
  const [nextCursor, setNextCusror] = React.useState("");

  const fetchData = () => {
    getRepliesAPI(
      { messageId: props.selectMessage, nextCursor: nextCursor },
      new Headers()
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
      getRepliesAPI({ messageId: props.selectMessage }, new Headers()).then(
        (v) => {
          setReplies(v.data.replies);
          if (v.data.next_cursor.length != 0) {
            setNextCusror(v.data.next_cursor);
            setHasMore(true);
          }
        }
      );
    }
  }, [props.show]);

  useClickOutside(modalRef, () => {
    props.setShow(false);
    props.setSelectMessage(-1);
  });

  const handleClick = () => {
    addReplyAPI(
      {
        message_id: props.selectMessage,
        user_id: "257e4caf-fb4b-43a1-a4b3-cca94f583bd5",
        username: "hank",
        content: text,
      },
      new Headers()
    ).then((v) => {
      if (v["status"] === "success") {
        const reply: ReplyType = v["data"];

        setReplies((prev) => {
          return [reply, ...prev];
        });
        setText("");
      }
    });
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
