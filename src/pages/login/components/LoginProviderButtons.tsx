import { Button } from "../../../components/ui/button";
import MicrosoftLogo from "../../../assets/ms-logo.svg";
import { ReactNode } from "react";

function LoginProviderButtons() {
  return (
    <section className="flex flex-col gap-2">
      <ProviderButton logo={MicrosoftLogo}>Microsoft</ProviderButton>
    </section>
  );
}

interface ProviderButtonProps {
  logo: string;
  children: ReactNode;
}

function ProviderButton({ children, logo }: ProviderButtonProps) {
  return (
    <Button className="flex gap-4" size="lg">
      <img src={logo} alt="" aria-label="icon" />
      <span>{children}</span>
    </Button>
  );
}

export default LoginProviderButtons;
