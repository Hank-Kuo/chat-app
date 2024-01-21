import React from "react";

import { ChannelType } from "../../apis/channel";
import { S } from "./styles";

interface SidebarProps {
  channels: ChannelType[];
  selectChannel: string;
  setSelectChannel: React.Dispatch<React.SetStateAction<string>>;
  setShow: React.Dispatch<React.SetStateAction<boolean>>;
}

const Sidebar: React.FC<SidebarProps> = ({
  channels,
  selectChannel,
  setSelectChannel,
  setShow,
}) => {
  return (
    <S.Wrapper>
      <S.Title>Channels</S.Title>
      <S.Container>
        {channels.map((channel) => {
          return (
            <S.Item
              key={channel.id}
              id={channel.id}
              active={channel.id == selectChannel}
              onClick={() => {
                setSelectChannel(channel.id);
              }}
            >
              {channel.name}
            </S.Item>
          );
        })}
      </S.Container>
      <S.Btn
        onClick={() => {
          setShow(true);
        }}
      >
        + Add channel
      </S.Btn>
    </S.Wrapper>
  );
};

export default Sidebar;
