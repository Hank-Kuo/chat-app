import React from "react";
import { useLoaderData } from "@remix-run/react";

import { loader } from "../../routes/_index";
import { S } from "./styles";

interface SidebarProps {
  selectChannel: string;
  setSelectChannel: React.Dispatch<React.SetStateAction<string>>;
  setShow: React.Dispatch<React.SetStateAction<boolean>>;
}

const Sidebar: React.FC<SidebarProps> = ({
  selectChannel,
  setSelectChannel,
  setShow,
}) => {
  const { channels } = useLoaderData<typeof loader>();

  return (
    <S.Wrapper>
      <S.Title>Channels</S.Title>
      <S.Container>
        {channels.map((channel: any) => {
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
