import PageHeading from "../../components/ui/PageHeading";
import LinkButton from "../../components/ui/LinkButton";
import SignupForm from "./components/SignupForm";

function SignupPage() {
  return (
    <>
      <PageHeading heading="Sign up">
        <LinkButton to="/login" size="sm" variant="ghost">
          Or log in
        </LinkButton>
      </PageHeading>
      <SignupForm />
    </>
  );
}

export default SignupPage;
