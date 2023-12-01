export function validEmail(text: string): boolean {
  const filter =
    /^([a-zA-Z0-9_\.\-])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z0-9]{2,4})+$/;
  return filter.test(text);
}

export function validPassword(text: string): [boolean, string] {
  const filter = /^(?=.*[a-zA-Z])(?=.*[0-9])/;
  if (filter.test(text) === false) {
    return [false, "Password must contain one letter and one number"];
  }
  if (text.length < 6) {
    return [false, "Password length must be greater than 6"];
  }
  return [true, ""];
}

export function validUsername(text: string): [boolean, string] {
  if (text.length < 3) {
    return [false, "Password length must be greater than 3"];
  }
  if (text.length > 10) {
    return [false, "Password length must be less than 10"];
  }
  return [true, ""];
}
