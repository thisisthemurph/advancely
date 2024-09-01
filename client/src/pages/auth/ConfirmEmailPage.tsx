import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import { checkEmailConfirmed } from "../../api/auth";
import LoadingDot from "../../components/ui/LoadingDot.tsx";

interface PageError {
  code: string;
  message: string;
}

interface TokenData {
  type: string;
  accessToken: string;
  refreshToken: string;
  expiresAt: Date | null;
}

function ConfirmEmailPage() {
  const navigate = useNavigate();
  const [tokenData, setTokenData] = useState<TokenData | null>(null);
  const [error, setError] = useState<PageError>({
    code: "",
    message: "",
  });

  useEffect(() => {
    if (!window.location.hash) {
      setError({
        code: "",
        message: "The URL looks to be broken, please check your email containing the provided email confirmation link.",
      })
      return;
    }

    const hash = window.location.hash.substring(1);
    const params = new URLSearchParams(hash);

    const expiresAtParam = params.get("expires_at");
    const errorCodeParam = params.get("error_code");
    const errorDescriptionParam = params.get("error_description");

    if (errorDescriptionParam) {
      setError({
        code: errorCodeParam ? errorCodeParam : "",
        message: errorDescriptionParam,
      });
      return;
    }

    const tokenTypeParam = params.get("token_type");
    const accessTokenParam = params.get("access_token");
    const refreshTokenParam = params.get("refresh_token");

    setTokenData({
      type: tokenTypeParam ?? "unknown",
      accessToken: accessTokenParam ?? "",
      refreshToken: refreshTokenParam ?? "",
      expiresAt: expiresAtParam ? new Date(Number(expiresAtParam) * 1000) : null,
    });
  }, []);

  useEffect(() => {
    if (!tokenData) {
      return;
    }

    const { accessToken, expiresAt } = tokenData;
    if (!accessToken || !expiresAt) {
      return;
    }

    if (expiresAt.getTime() <= new Date().getTime()) {
      setError({
        code: "",
        message: "The verification token has expired, please request a new one to verify your email address.",
      })
      return;
    }

    checkEmailConfirmed({ token: accessToken }).then(() => {
      navigate("/login");
    }).catch(() => {
      setError({
        code: "",
        message: "We were not able to verify that your email address has been verified, please try logging in or requesting a new verification email if you cannot.",
      });
    });
  }, [tokenData, navigate]);

  if (error.message) {
    return (
      <section className="space-y-4">
        <p className="text-xl">Error verifying email address</p>
        {error.code && <p><strong>Error code: { error.code }</strong></p>}
        <p>{error.message}.</p>
      </section>
    );
  }

  return (
    <section className="space-y-8 mt-16">
      <p className="text-xl text-center">Please wait whilst your email address is confirmed.</p>
      <section className="flex flex-col gap-2 items-center justify-center">
        <div className="flex items-center justify-center">
          <LoadingDot size="lg" />
          <LoadingDot size="lg" />
          <LoadingDot size="lg" />
          <LoadingDot size="lg" />
          <LoadingDot size="lg" />
          <LoadingDot size="lg" />
        </div>
        <p className="animate-pulse font-mono">Loading...</p>
      </section>
    </section>
  );
}

export default ConfirmEmailPage;
