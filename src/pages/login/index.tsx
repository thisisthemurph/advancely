import { useForm } from "react-hook-form";
import PageHeading from "../../components/ui/PageHeading";
import { Button } from "../../components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../../components/ui/form";
import { Input } from "../../components/ui/input";
import LinkButton from "../../components/ui/LinkButton";

import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const formSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6, {
    message: "Your password must be at least 6 characters long",
  }),
});

type FormInputs = z.infer<typeof formSchema>;

function LoginPage() {
  const form = useForm<FormInputs>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  function onSubmit(values: FormInputs) {
    console.log(values);
  }

  return (
    <>
      <PageHeading heading="Log in">
        <LinkButton to="/signup" size="sm" variant="secondary">
          Or sign up
        </LinkButton>
      </PageHeading>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input
                    autoFocus={true}
                    type="email"
                    placeholder="you@email.com"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Password</FormLabel>
                <FormControl>
                  <Input
                    type="password"
                    placeholder="* * * * * * * * * * * * * *"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit">Log in</Button>
        </form>
      </Form>
    </>
  );
}

export default LoginPage;
