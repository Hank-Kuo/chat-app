const header = (token?: string): Headers => {
  const header = new Headers();

  header.append("Content-Type", "application/json");
  if (token) {
    header.append("Authorization", `Bearer ${token}`);
  }

  return header;
};

export default header;
