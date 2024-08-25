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
  name: z.string().min(4, {
    message: "Your company name should be at least 4 characters long",
  }),
  email: z.string().email(),
  password: z.string().min(6, {
    message: "Your password must be at least 6 characters long",
  }),
});

type FormInputs = z.infer<typeof formSchema>;

function SignupPage() {
  const form = useForm<FormInputs>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      email: "",
      password: "",
    },
  });

  function onSubmit(values: FormInputs) {
    console.log(values);
  }

  return (
    <>
      <PageHeading heading="Sign up">
        <LinkButton to="/login" size="sm" variant="ghost">
          Or log in
        </LinkButton>
      </PageHeading>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Company name</FormLabel>
                <FormControl>
                  <Input
                    autoFocus={true}
                    placeholder="Your Company LTD"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input type="email" placeholder="you@email.com" {...field} />
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
          <Button type="submit">Sign up</Button>
        </form>
      </Form>
    </>
  );
}

export default SignupPage;
