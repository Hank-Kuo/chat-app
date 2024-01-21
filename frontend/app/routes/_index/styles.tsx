import { Form } from "@remix-run/react";
import styled from "@emotion/styled";

export const S = {
  Wrapper: styled.div`
    display: flex;
    width: 100%;
  `,
  Container: styled.div`
    display: flex;
    flex: 1;
    flex-direction: column;
    position: relative;
  `,
  ContainerTitle: styled.h3`
    margin-top: 50px;
    text-align: center;
    font-size: 40px;
  `,
  Box: styled.div`
    flex: 1;
    overflow: scroll;
    margin: 20px 0;
    padding: 0 20px;
    display: flex;
    flex-direction: column-reverse;
  `,
  InputBox: styled.div`
    display: flex;
    padding: 30px 20px;
  `,
  Input: styled.textarea`
    border-radius: 8px;
    width: 100%;
    padding: 9px 15px;
    outline: none;
  `,
  SubmitBtn: styled.button`
    margin-left: 20px;
    border-radius: 5px;
    background-color: black;
    color: white;
    border: 2px solid;
    cursor: pointer;
  `,
  Item: styled.div`
    padding: 10px;
    :hover {
      border: 1px solid;
      box-shadow: 1px 1px;
      p {
        display: block;
      }
    }
  `,
  ItemInfoBox: styled.div`
    display: flex;
    align-items: center;
  `,
  ItemName: styled.p`
    font-weight: bold;
  `,
  ItemTime: styled.p`
    margin-left: 15px;
    font-size: 12px;
    color: lightgrey;
  `,
  ItemBox: styled.div`
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 10px;
  `,
  ItemMessage: styled.p`
    word-break: break-all;
  `,
  ItemReply: styled.p`
    margin-left: 5px;
    cursor: pointer;
    display: none;
  `,
  JoinBox: styled(Form)`
    position: absolute;
    height: 100%;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  `,
  JoinInput: styled.input`
    display: none;
  `,
  JoinGreyBox: styled.div`
    height: 100%;
    width: 100%;
    position: absolute;
    background: rgb(224, 224, 224, 0.5);
  `,
  JoinBtn: styled.button`
    background: none;
    color: white;
    border: 2px solid;
    font-size: 20px;
    padding: 10px 20px;
    cursor: pointer;
    z-index: 1000;
    :hover {
      background: white;
      color: black;
    }
  `,
};

export const M = {
  Wrapper: styled.div`
    position: absolute;
    top: 0;
    right: 0;
    background: rgba(224, 224, 224, 0.3);
    width: 100%;
    height: 100%;
    z-index: 1000;
  `,
  ChannelContainer: styled.div`
    display: flex;
    height: 100%;
    width: 100%;
    align-items: center;
    justify-content: center;
    border-radius: 10px;
  `,
  ChannelTitle: styled.h3`
    text-align: center;
  `,
  ChannelInput: styled.input`
    margin-top: 50px;
    width: 100%;
    border: none;
    padding: 15px 10px;
    outline: none;
    border-radius: 5px;
  `,
  ChannelBtn: styled.button`
    background: none;
    margin-top: 50px;
    border: 1px solid white;
    color: white;
    padding: 10px;
    border-radius: 5px;
    cursor: pointer;
  `,
  ChannelForm: styled(Form)`
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 350px;
    height: 300px;
    padding: 30px;
    background: black;
  `,
  Container: styled.div`
    display: flex;
    flex-direction: column;
    width: 400px;
    height: 100%;
    background: black;
    position: absolute;
    right: 0;
  `,
  Title: styled.div`
    margin-top: 40px;
    text-align: center;
    font-size: 50px;
  `,
  ReplyBox: styled.div`
    display: flex;
    flex-direction: column-reverse;
    flex: 1;
    overflow: scroll;
    padding: 10px 20px;
    > * + * {
      margin-top: 15px;
    }
  `,
  Text: styled.p``,
};
