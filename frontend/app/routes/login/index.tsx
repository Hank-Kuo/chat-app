import { useActionData } from "@remix-run/react";
import {
  json,
  redirect,
  ActionFunctionArgs,
  LoaderFunctionArgs,
} from "@remix-run/node";

import { getSession, commitSession } from "../../lib/utils/session";
import { validEmail, validPassword } from "../../lib/utils/validRegex";
import { loginAPI } from "../../apis/auth";
import { S } from "./styles";

export async function loader({ request }: LoaderFunctionArgs) {
  const session = await getSession(request.headers.get("Cookie"));

  if (session.has("userInfo")) {
    return redirect("/");
  }
  const data = { error: "Not found user in cookies" };

  return json(data, {
    headers: {
      "Set-Cookie": await commitSession(session),
    },
  });
}

export async function action({ request }: ActionFunctionArgs) {
  const session = await getSession(request.headers.get("Cookie"));
  const form = await request.formData();
  const email = form.get("email") as string;
  const password = form.get("password") as string;

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
  const [isValid, errorText] = validPassword(password);

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
  const res = await loginAPI({ email: email, password: password });

  if (res["status"] !== "success") {
    return json(
      { error: res["message"] },
      {
        headers: {
          "Set-Cookie": await commitSession(session),
        },
      }
    );
  }
  session.set("userInfo", res["data"]);

  return redirect("/", {
    headers: {
      "Set-Cookie": await commitSession(session),
    },
  });
}

export default function LoginView() {
  const data = useActionData<typeof action>();

  return (
    <S.Wrapper>
      <S.Title>Login</S.Title>
      <S.Form method="post">
        <S.Input type="text" name="email" placeholder="email" />
        <S.Input type="password" name="password" placeholder="password" />
        <S.Button type="submit">Submit</S.Button>
      </S.Form>
      <S.ErrorText>{data?.error}</S.ErrorText>
    </S.Wrapper>
  );
}
