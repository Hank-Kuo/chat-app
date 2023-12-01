import React, { Suspense } from "react";
import useClickOutside from "../../lib/hook/useClickOutside";
import { S, M } from "./styles";

interface ReplyModalProps {
  show: boolean;
  setShow: React.Dispatch<React.SetStateAction<boolean>>;
}

const ReplyModal = (props: ReplyModalProps) => {
  const modalRef = React.useRef(null);
  const messagesRef = React.createRef<HTMLDivElement>();

  React.useEffect(() => {
    if (messagesRef.current) {
      messagesRef.current.scrollTop = messagesRef.current.scrollHeight;
    }
  }, [props.show]);

  useClickOutside(modalRef, () => props.setShow(false));

  return (
    <Suspense fallback={<div>Loading...</div>}>
      {props.show ? (
        <M.Wrapper>
          <M.Container ref={modalRef}>
            <M.Title>Reply</M.Title>
            <M.ReplyBox ref={messagesRef}>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>message</S.ItemMessage>
                </S.ItemBox>
              </S.Item>
              <S.Item>
                <S.ItemInfoBox>
                  <S.ItemName>Hank</S.ItemName>
                  <S.ItemTime>2023/10/1</S.ItemTime>
                </S.ItemInfoBox>
                <S.ItemBox>
                  <S.ItemMessage>
                    messagemessagemessagemessagemessagemessagemessage
                  </S.ItemMessage>
                </S.ItemBox>
              </S.Item>
            </M.ReplyBox>
            <S.InputBox>
              <S.Input />
              <S.SubmitBtn>Submit</S.SubmitBtn>
            </S.InputBox>
          </M.Container>
        </M.Wrapper>
      ) : null}
    </Suspense>
  );
};

export default ReplyModal;
