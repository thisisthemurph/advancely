import PageHeading from "../../components/ui/PageHeading";
import LinkButton from "../../components/ui/LinkButton";
import Divider from "../../components/ui/Divider";
import LoginForm from "./components/LoginForm";
import LoginProviderButtons from "./components/LoginProviderButtons";
import {useNavigate} from "react-router-dom";

function LoginPage() {
  const navigate = useNavigate();

  const onLoginSuccess = () => {
    navigate("/dashboard");
  }

  return (
    <>
      <PageHeading heading="Log in">
        <LinkButton to="/signup" size="sm" variant="ghost">
          Or sign up
        </LinkButton>
      </PageHeading>
      <LoginForm onSuccess={onLoginSuccess} />
      <Divider text="or continue with" />
      <LoginProviderButtons />
    </>
  );
}

export default LoginPage;
