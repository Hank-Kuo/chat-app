import axios from "axios";

import { getSession } from "../lib/utils/session";
// import { getUserCookies } from "../core/utils/cookie";

axios.defaults.timeout = 5000; // 5 second
axios.defaults.baseURL = "http://localhost:9000/api";

axios.interceptors.request.use(
  async (config) => {
    // console.log("123");
    const session = await getSession();

    // // const userInfo = await getUserCookies();
    // if (userInfo) {
    //   config.headers.Authorization = `Bearer ${userInfo.token}`;
    // }
    console.log("interceptor");
    config.headers["Access-Control-Allow-Origin"] = "*";
    return config;
  },
  (error) => Promise.reject(error)
);

axios.interceptors.response.use(
  (response) => response,
  (err) => Promise.reject(err)
);

export function get(url: string, params = {}) {
  return new Promise((resolve, reject) => {
    axios
      .get(url, { params })
      .then((response) => {
        resolve(response.data);
      })
      .catch((err) => {
        reject(err);
      });
  });
}
