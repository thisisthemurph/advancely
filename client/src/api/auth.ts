interface SignupRequest {
  name: string;
  firstName: string;
  lastName: string;
  email: string;
  password: string;
}

interface SignupResponse {
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
    throw new Error("Error signing up");
  }

  return resp.json();
};
