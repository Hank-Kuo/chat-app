const OPTIONS: Intl.DateTimeFormatOptions = {
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
  hour: "numeric",
  minute: "numeric",
  timeZone: "Asia/Taipei", // Specify the target time zone
};

export const convertDate = (date: string): string => {
  const dateUTC = new Date(date);

  const dateUTC8 = new Date(dateUTC.toLocaleString("en-US", OPTIONS));

  const year = dateUTC8.getFullYear();
  const month = (dateUTC8.getMonth() + 1).toString().padStart(2, "0");
  const day = dateUTC8.getDate().toString().padStart(2, "0");
  const hours = dateUTC8.getHours().toString().padStart(2, "0");
  const minutes = dateUTC8.getMinutes().toString().padStart(2, "0");

  const formattedDateString = `${year}/${month}/${day} ${hours}:${minutes}`;

  return formattedDateString;
};
