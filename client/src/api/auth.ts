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
  const signupEndpoint = `${import.meta.env.VITE_API_BASE_URL}/auth/signup`;

  const resp = await fetch(signupEndpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
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

export interface ConfirmEmailConfirmedRequest {
  token: string;
}

export const checkEmailConfirmed = async (data: ConfirmEmailConfirmedRequest): Promise<boolean> => {
  const endpoint = `${import.meta.env.VITE_API_BASE_URL}/auth/confirm-email`;

  const resp = await fetch(endpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  return resp.ok;
}
