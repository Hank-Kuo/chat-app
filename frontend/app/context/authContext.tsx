import React from "react";

import { useNavigate } from "@remix-run/react";

interface userInfo {
  username: string;
  email: string;
  channels: string[];
}
interface AuthContextType {
  user: userInfo;
  isLogin: boolean;
  setUser: React.Dispatch<React.SetStateAction<userInfo>>;
}

const AuthContext = React.createContext<AuthContextType>(null!);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  let [user, setUser] = React.useState({
    username: "",
    email: "",
    channels: [""],
  });
  const [isLogin, setIsLogin] = React.useState(false);
  const navigate = useNavigate();

  React.useEffect(() => {
    setIsLogin(false);
    // if (isLogin === false) {
    //   navigate("/login");
    // }
  }, []);

  let value = { isLogin, user, setUser };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  return React.useContext(AuthContext);
}
