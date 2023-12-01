import { Form } from "@remix-run/react";
import styled from "@emotion/styled";

export const S = {
  Wrapper: styled.div`
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    width: 100%;
    height: 100%;
  `,
  Form: styled(Form)`
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    margin-top: 40px;
    > * + * {
      margin-top: 20px;
    }
  `,
  Input: styled.input`
    width: 350px;
    padding: 10px;
    border-radius: 5px;
    outline: none;
    color: black;
  `,
  Button: styled.button`
    margin-top: 50px;
    border: 1px solid;
    padding: 10px 30px;
    border-radius: 5px;
    cursor: pointer;
  `,
  Title: styled.h3`
    font-size: 40px;
  `,
  ErrorText: styled.p`
    color: red;
    margin-top: 40px;
  `,
};
