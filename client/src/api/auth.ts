import { ErrorResponseSchema } from "./api";

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
  const unknownSignupError = {
    message:
      "There has been an unknown error signing you up, please try again later.",
  };

  const resp = await fetch(signupEndpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  if (!resp.ok) {
    const body = await resp.json();
    const error = ErrorResponseSchema.safeParse(body);

    if (error.success) {
      return Promise.reject(error.data);
    }
    return Promise.reject(
      body?.message ? { message: body.message } : unknownSignupError
    );
  }

  return resp.json();
};
