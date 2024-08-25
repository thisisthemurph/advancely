import PageHeading from "../../components/ui/PageHeading";
import LinkButton from "../../components/ui/LinkButton";
import Divider from "../../components/ui/Divider";
import LoginForm from "./components/LoginForm";
import LoginProviderButtons from "./components/LoginProviderButtons";

function LoginPage() {
  return (
    <>
      <PageHeading heading="Log in">
        <LinkButton to="/signup" size="sm" variant="secondary">
          Or sign up
        </LinkButton>
      </PageHeading>
      <LoginForm />
      <Divider text="or continue with" />
      <LoginProviderButtons />
    </>
  );
}

export default LoginPage;
