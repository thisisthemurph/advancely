import { useState } from "react";
import { useForm, UseFormReturn } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { Input } from "../../../components/ui/input";
import { Button } from "../../../components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../../../components/ui/form";
import { signup } from "../../../api/auth";

const CompanyNameInfo =
  "Set the name of your company, this is the name all of your employees will see themselves under when they sign in.";
const FirstNameInfo = "The first name of the initial admin user.";
const LastNameInfo = "The last name of the initial admin user.";
const EmailInfo =
  "The email address used here will be considered a super admin account and will have permissions to create, read, update, and delete all content. Permissions can be tailored later to give other users different privilages.";

const formSchema = z.object({
  name: z.string().min(4, {
    message: "Your company name should be at least 4 characters long",
  }),
  firstName: z.string().min(1, "A first name must be provided"),
  lastName: z.string().min(1, "A last name must be provided"),
  email: z.string().email(),
  password: z.string().min(6, {
    message: "Your password must be at least 6 characters long",
  }),
});

type FormSchema = z.infer<typeof formSchema>;

interface SignupFormParams {
  onSignupComplete: () => void;
}

function SignupForm({ onSignupComplete }: SignupFormParams) {
  const [emailPlaceholder, setEmailPlaceholder] = useState(
    "your.name@company.com"
  );

  const form = useForm<FormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      firstName: "",
      lastName: "",
      email: "",
      password: "",
    },
  });

  const { mutateAsync: signupMitation, isPending: isSignupPending } =
    useMutation({
      mutationFn: signup,
      onSuccess: () => {
        onSignupComplete();
      },
    });

  async function onSubmit(values: FormSchema) {
    try {
      await signupMitation(values);
    } catch (e) {
      console.error(e);
    }
  }

  function onCompanyNameChange(form: UseFormReturn<FormSchema>) {
    const companyName = form.getValues("name");
    if (!companyName) {
      setEmailPlaceholder("your.name@company.com");
      return;
    }

    const email = `your.name@${companyName}.com`;
    setEmailPlaceholder(email.toLowerCase().replace(" ", ""));
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel infoText={EmailInfo}>Company name</FormLabel>
              <FormControl onChange={() => onCompanyNameChange(form)}>
                <Input
                  autoFocus={true}
                  placeholder="Your Company LTD"
                  {...field}
                />
              </FormControl>
              <FormMessage />
              <FormDescription className="hidden">
                {CompanyNameInfo}
              </FormDescription>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="firstName"
          render={({ field }) => (
            <FormItem>
              <FormLabel infoText={FirstNameInfo}>First name</FormLabel>
              <FormControl>
                <Input type="text" placeholder="Your first name" {...field} />
              </FormControl>
              <FormMessage />
              <FormDescription className="hidden">
                {FirstNameInfo}
              </FormDescription>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="lastName"
          render={({ field }) => (
            <FormItem>
              <FormLabel infoText={LastNameInfo}>Last name</FormLabel>
              <FormControl>
                <Input type="text" placeholder="Your last name" {...field} />
              </FormControl>
              <FormMessage />
              <FormDescription className="hidden">
                {LastNameInfo}
              </FormDescription>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel infoText={EmailInfo}>
                <span>Email</span>
              </FormLabel>
              <FormControl>
                <Input type="email" placeholder={emailPlaceholder} {...field} />
              </FormControl>
              <FormMessage />
              <FormDescription className="hidden">{EmailInfo}</FormDescription>
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
        <Button
          type="submit"
          disabled={isSignupPending}
          loading={isSignupPending}
        >
          Sign up
        </Button>
      </form>
    </Form>
  );
}

export default SignupForm;
