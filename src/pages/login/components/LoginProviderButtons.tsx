import { Button } from "../../../components/ui/button";

function LoginProviderButtons() {
  return (
    <section className="flex flex-col gap-2">
      <Button size="lg">Google</Button>
      <Button size="lg">Apple</Button>
    </section>
  );
}

export default LoginProviderButtons;
