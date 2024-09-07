import { createContext, ReactNode, useEffect, useState } from "react";
import { handleErrorResponse } from "../api/api.ts";

const SESSION_STORE_KEY = "session";
const IS_AUTHED_STORE_KEY = "authed";
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

export interface AuthContextProps {
  isAuthenticated: boolean;
  user: SessionUser | null;
  session: Session | null;
  login: (params: LoginParams) => Promise<Session>;
  loginWithToken: (token: string) => Promise<void>;
  logout: () => Promise<boolean>;
  updateSession: (session: Session) => Promise<void>;
}

export interface SessionUserCompany {
  id: string;
  name: string;
}

export interface SessionUser {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  company?: SessionUserCompany;
}

export interface Session {
  sub: string
  aud: string
  role: string
  email: string
  accessToken: string
  refreshToken: string
  expiresAt: string
  user?: SessionUser;
}

export interface LoginParams {
  email: string;
  password: string;
}

export const AuthContext = createContext<AuthContextProps | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => {
    const state = localStorage.getItem(IS_AUTHED_STORE_KEY);
    return state === "true";
  });

  const [session, setSession] = useState<Session | null>(() => {
    const s = localStorage.getItem(SESSION_STORE_KEY);
    return s ? JSON.parse(s) : null;
  });

  useEffect(() => {
    localStorage.setItem(IS_AUTHED_STORE_KEY, isAuthenticated.toString());
  }, [isAuthenticated]);

  useEffect(() => {
    if (session) {
      localStorage.setItem(SESSION_STORE_KEY, JSON.stringify(session));
    } else {
      localStorage.removeItem(SESSION_STORE_KEY)
    }
  }, [session]);

  const post = async function <TReq>(endpoint: string, data: TReq): Promise<Response> {
    return await fetch(endpoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify(data),
    });
  }

  const login = async (data: LoginParams): Promise<Session> => {
    const endpoint = `${API_BASE_URL}/auth/login`;
    console.log({ endpoint, API_BASE_URL, fromMeta: import.meta.env.VITE_API_BASE_URL });
    const resp = await post(endpoint, data);
    if (!resp.ok) {
      return Promise.reject(await handleErrorResponse(resp));
    }

    const newSession = await resp.json();
    setSession(newSession);
    setIsAuthenticated(true);
    return newSession;
  }

  const loginWithToken = async (token: string) => {
    console.log(token);
    throw Error("Not implemented");
  }

  const logout = async () => {
    const endpoint = `${API_BASE_URL}/auth/logout`;
    const resp = await post(endpoint, {})

    setSession(null);
    setIsAuthenticated(false);

    return resp.ok;
  }

  const updateSession = async (newSession: Session) => {
    setSession({ ...newSession });
  }

  return (
    <AuthContext.Provider value={{
      isAuthenticated,
      login,
      loginWithToken,
      logout,
      session,
      updateSession,
      user: session?.user || null,
    }}>
      {children}
    </AuthContext.Provider>
  )
}
