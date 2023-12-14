import { useActionData } from "@remix-run/react";
import { json, redirect, ActionFunctionArgs } from "@remix-run/node";

import { getSession, commitSession } from "../../lib/utils/session";
import {
  validEmail,
  validPassword,
  validUsername,
} from "../../lib/utils/validRegex";
import { registerAPI } from "../../apis/auth";
import { S } from "./styles";

export async function action({ request }: ActionFunctionArgs) {
  const session = await getSession(request.headers.get("Cookie"));
  const form = await request.formData();
  const email = `${form.get("email")}`;
  const password = `${form.get("password")}`;
  const username = `${form.get("username")}`;

  if (validEmail(email) === false) {
    return json(
      { error: "Email format error" },
      {
        headers: {
          "Set-Cookie": await commitSession(session),
        },
      }
    );
  }
  let [isValid, errorText] = validPassword(password);

  if (isValid === false) {
    return json(
      { error: errorText },
      {
        headers: {
          "Set-Cookie": await commitSession(session),
        },
      }
    );
  }

  [isValid, errorText] = validUsername(username);
  if (isValid === false) {
    return json(
      { error: errorText },
      {
        headers: {
          "Set-Cookie": await commitSession(session),
        },
      }
    );
  }

  const res = await registerAPI({
    username: username,
    email: email,
    password: password,
  });

  if (res["status"] === false) {
    return json(
      { error: res["errMessage"] },
      {
        headers: {
          "Set-Cookie": await commitSession(session),
        },
      }
    );
  }

  return redirect("/login", {
    headers: {
      "Set-Cookie": await commitSession(session),
    },
  });
}

export default function RegisterView() {
  const data = useActionData<typeof action>();
  return (
    <S.Wrapper>
      <S.Title>Register</S.Title>
      <S.Form method="post">
        <S.Input type="text" name="email" placeholder="email" />
        <S.Input type="text" name="username" placeholder="username" />
        <S.Input type="password" name="password" placeholder="password" />
        <S.Button type="submit">Submit</S.Button>
      </S.Form>
      <S.ErrorText>{data?.error}</S.ErrorText>
    </S.Wrapper>
  );
}
