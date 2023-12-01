import styled from "@emotion/styled";

type ItemProps = {
  active: boolean;
};

export const S = {
  Wrapper: styled.div`
    flex-basis: 200px;
    display: flex;
    flex-direction: column;
    padding: 30px 0px;
    border-right: 1px solid;
  `,
  Title: styled.h3`
    margin-left: 15px;
  `,
  Container: styled.ul`
    flex: 1;
    max-width: 200px;
    margin-top: 25px;
    text-align: center;
    overflow: scroll;
    padding: 0;
  `,
  Item: styled.li<ItemProps>`
    padding: 10px;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    cursor: pointer;
    background-color: ${({ active }) =>
      active ? "rgba(224, 224, 224, 0.25)" : "none"};
    :hover {
      background-color: rgba(224, 224, 224, 0.25);
    }
  `,
  Btn: styled.button`
    border: none;
    background: none;
    color: white;
    font-size: 20px;
    padding: 10px 0;
    cursor: pointer;
  `,
};
