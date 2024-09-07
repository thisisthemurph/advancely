import { useState } from "react";

import PageHeading from "../../components/ui/PageHeading";
import LinkButton from "../../components/ui/LinkButton";
import SignupForm from "./components/SignupForm";
import SignupCompleteMessage from "./components/SignupCompleteMessage";

function SignupPage() {
  const [signupComplete, setSignupComplete] = useState(false);

  const handleSignupCompleted = () => {
    setSignupComplete(true);
  };

  return (
    <>
      <PageHeading heading="Sign up">
        <LinkButton
          to="/login"
          size="sm"
          variant={signupComplete ? "default" : "ghost"}
        >
          {signupComplete ? "log in" : "or log in"}
        </LinkButton>
      </PageHeading>

      {signupComplete ? (
        <SignupCompleteMessage />
      ) : (
        <SignupForm onSignupComplete={handleSignupCompleted} />
      )}
    </>
  );
}

export default SignupPage;
