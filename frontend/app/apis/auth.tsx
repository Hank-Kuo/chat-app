interface loginPayload {
  email: string;
  password: string;
}

interface registerPayload {
  email: string;
  username: string;
  password: string;
}

export async function loginAPI(data: loginPayload) {
  const res = await fetch(`http://localhost:9000/api/login`, {
    method: "POST",
    body: JSON.stringify(data),
  });

  return await res.json();
}

export async function registerAPI(data: registerPayload) {
  const res = await fetch(`http://localhost:9000/api/register`, {
    method: "POST",
    body: JSON.stringify(data),
  });
  return await res.json();
}
