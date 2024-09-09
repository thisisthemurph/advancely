import { createContext, ReactNode, useEffect, useState } from "react";
import ApiError from "../api/error.ts";

const SESSION_STORE_KEY = "session";
const IS_AUTHED_STORE_KEY = "authed";
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

export interface AuthContextProps {
  isAuthenticated: boolean;
  user: SessionUser | null;
  session: Session | null;
  login: (params: LoginParams) => Promise<void>;
  logout: () => Promise<void>;
  updateSession: (session: Session) => Promise<void>;
  signup: (data: SignupRequest) => Promise<SignupResponse>;
  emailConfirmed: (data: { token: string }) => Promise<boolean>;
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

export interface SignupRequest {
  name: string;
  firstName: string;
  lastName: string;
  email: string;
  password: string;
}

export interface SignupResponse {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  companyName: string;
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

  const post = async function <TReq, TResp>(endpoint: string, data: TReq): Promise<TResp> {
    const resp = await fetch(endpoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify(data),
    });

    if (!resp.ok) {
      const errorBody = await resp.json();
      const message = errorBody
        ? errorBody.message
        : `Error making post. Status: ${resp.status} statusText: ${resp.statusText}`;

      throw new ApiError(resp.status, resp.statusText, message);
    }

    const responseData: TResp = await resp.json();
    return responseData;
  }

  const login = async ({email, password}: LoginParams) => {
    const endpoint = `${API_BASE_URL}/auth/login`;

    try {
      const session = await post<LoginParams, Session>(endpoint, {email, password});
      setSession(session);
      setIsAuthenticated(true);
    } catch {
      setSession(null);
      setIsAuthenticated(false);
      throw new Error("There was an error logging you in.");
    }
  }

  const signup = async (data: SignupRequest): Promise<SignupResponse> => {
    const endpoint = `${API_BASE_URL}/auth/signup`;
    const genericError = "There has been an unknown error signing you up, please try again later."

    try {
      return await post<SignupRequest, SignupResponse>(endpoint, data);
    } catch (error) {
      if (error instanceof ApiError && (error.status === 400)) {
        throw new Error("Please ensure the form is complete.")
      } else if (error instanceof ApiError && error.status === 500) {
        throw error;
      }
      throw new Error(genericError);
    }
  }

  const logout = async () => {
    const endpoint = `${API_BASE_URL}/auth/logout`;

    try {
      await post(endpoint, {});
    } finally {
      setSession(null);
      setIsAuthenticated(false);
    }
  }

  const emailConfirmed = async (data: { token: string }): Promise<boolean> => {
    const endpoint = `${import.meta.env.VITE_API_BASE_URL}/auth/confirm-email`;

    try {
      await post<{token: string}, unknown>(endpoint, data);
      return true;
    } catch {
      return false;
    }
  }

  const updateSession = async (newSession: Session) => {
    setSession({ ...newSession });
  }

  return (
    <AuthContext.Provider value={{
      emailConfirmed,
      isAuthenticated,
      login,
      logout,
      session,
      signup,
      updateSession,
      user: session?.user || null,
    }}>
      {children}
    </AuthContext.Provider>
  )
}
