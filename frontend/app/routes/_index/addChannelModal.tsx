import React, { Suspense } from "react";

import useClickOutside from "../../lib/hook/useClickOutside";
import { M } from "./styles";

interface ChannelModalProps {
  show: boolean;
  setShow: React.Dispatch<React.SetStateAction<boolean>>;
}

const ChannelModal = (props: ChannelModalProps) => {
  const ref = React.useRef(null);

  useClickOutside(ref, () => {
    return props.setShow(false);
  });

  return (
    <Suspense fallback={<div>Loading...</div>}>
      {props.show ? (
        <M.Wrapper>
          <M.ChannelContainer>
            <M.ChannelForm
              ref={ref}
              method="post"
              onSubmit={() => {
                props.setShow(false);
              }}
            >
              <M.ChannelTitle>Add Channel</M.ChannelTitle>

              <M.ChannelInput
                type="text"
                placeholder="Channel Name"
                name="channelName"
              />
              <M.ChannelBtn name="formAction" value="addAction">
                Submit
              </M.ChannelBtn>
            </M.ChannelForm>
          </M.ChannelContainer>
        </M.Wrapper>
      ) : null}
    </Suspense>
  );
};

export default ChannelModal;
