import React, { Suspense } from "react";
import { S } from "./styles";
import { channelsType } from "../../apis/channel";
const DATA = [
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
];

interface MessageContainerProps {
  userChannels: channelsType[];
  selectChannel: string;
  setShowReply: React.Dispatch<React.SetStateAction<boolean>>;
}

export default function MessageContainer(props: MessageContainerProps) {
  const messagesRef = React.createRef<HTMLDivElement>();
  const [showJoin, setShowJoin] = React.useState(true);

  React.useEffect(() => {
    if (messagesRef.current) {
      messagesRef.current.scrollTop = messagesRef.current.scrollHeight;
    }

    if (props.userChannels.find((v) => v.id === props.selectChannel)) {
      setShowJoin(false);
    } else {
      setShowJoin(true);
    }
  }, [props.selectChannel]);

  return (
    <>
      <Suspense fallback={<div>Loading...</div>}>
        <S.Box ref={messagesRef}>
          {DATA.map((v) => {
            return (
              <S.Item key={v.id}>
                <S.ItemInfoBox>
                  <S.ItemName>{showJoin ? "Default" : v.name}</S.ItemName>
                  <S.ItemTime>
                    {showJoin ? "2023/01/01" : v.creatdAt}
                  </S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>
                    {showJoin ? "⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆⋆" : v.message}
                  </S.ItemMessage>
                  <S.ItemReply
                    onClick={() => (showJoin ? null : props.setShowReply(true))}
                  >
                    REPLY
                  </S.ItemReply>
                </S.ItemBox>
              </S.Item>
            );
          })}
        </S.Box>
        <S.InputBox>
          <S.Input />
          <S.SubmitBtn>Submit</S.SubmitBtn>
        </S.InputBox>
        {showJoin ? (
          <S.JoinBox method="post" onSubmit={() => setShowJoin(false)}>
            <S.JoinInput name="channelID" value={props.selectChannel} />
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
