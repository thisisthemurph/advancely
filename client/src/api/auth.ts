import { handleErrorResponse } from "./api";

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

export const signup = async (data: SignupRequest): Promise<SignupResponse> => {
  const endpoint = `${import.meta.env.VITE_API_BASE_URL}/auth/signup`;

  const resp = await fetch(endpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
    credentials: "include",
  });

  if (!resp.ok) {
    return Promise.reject(
      await handleErrorResponse(
        resp,
        "There has been an unknown error signing you up, please try again later."
      )
    );
  }

  return resp.json();
};

export interface CheckEmailConfirmedRequest {
  token: string;
}

export const checkEmailConfirmed = async (data: CheckEmailConfirmedRequest): Promise<boolean> => {
  const endpoint = `${import.meta.env.VITE_API_BASE_URL}/auth/confirm-email`;

  const resp = await fetch(endpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
    credentials: "include",
  });

  return resp.ok;
}
